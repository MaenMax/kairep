package routercore

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	fernet "fernet-go"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/config"
	"empowerthings.com/autoendpoint/routercore/db/cassandra"
	"empowerthings.com/autoendpoint/routercore/db/redis/counterdb"
	"empowerthings.com/autoendpoint/routercore/model"
	"empowerthings.com/autoendpoint/routercore/vapid"
	"empowerthings.com/autoendpoint/utils"
	"github.com/gocql/gocql"
	"github.com/statsd"
	"golang.org/x/exp/utf8string"
)

const (
	TTL_DURATION time.Duration = 1000 * 1000 * 1000 * 60 * 60 * 24 * 365 * 100
	LEGACY       string        = ""
)

type T_RouterWorker struct {
	encrypted_endpoint      string
	body                    []byte
	headers                 http.Header
	crypto_key              string
	debug                   bool
	session                 *gocql.Session
	w                       http.ResponseWriter
	crypto_key_utf8_encoded *utf8string.String
	crypto_headers          model.Headers
	vapid_headers           model.Vapid_Headers
	topic                   string
	raw_chid                string
	uaid                    string
	sorted_key              string
	app_server_ip			string
	sortkey_timestamp       int
	msg                     string
	ttl                     int
	protocol                int
	vapid                   bool
	vapid_pub_key           string
	stats                   *statsd.Client
	multicast               bool
	is_limited              bool

}

func NewRouterWorker(encrypted_endpoint string, request_body []byte, request_headers http.Header, crypto_key string, debug bool, db_session *gocql.Session, w http.ResponseWriter, protocol int, vapid bool, statistics *statsd.Client, multicast bool, is_limited bool) *T_RouterWorker {

	var tmp T_RouterWorker
	tmp = T_RouterWorker{}

	tmp.encrypted_endpoint = encrypted_endpoint
	tmp.body = request_body
	tmp.headers = request_headers
	tmp.crypto_key = crypto_key
	tmp.debug = debug
	tmp.session = db_session
	tmp.w = w
	tmp.protocol = protocol
	tmp.vapid = vapid
	tmp.stats = statistics
	tmp.multicast = multicast
	tmp.is_limited = is_limited

	return &tmp

}

func (rw *T_RouterWorker) RouteNotification() {
	var dataA *model.Payload_A // Webpush payload
	var dataB *model.Payload_B // An empty(Webpush,SimplePush,VAPID) payload
	var dataC *model.Payload_C // SimplePush payload
	var resp_body string
	_conf := config.GetREPConfig() // Reading configuration 

	// Checking if this is a VAPID request.

	var version int
	version = 1

	if rw.vapid {
		version = 2
		rw.protocol = 2 // WebPush request
	}

	forwarded_header := rw.headers.Get("X-Forwarded-For")
	if len(forwarded_header) > 0 {
		items := strings.Split(forwarded_header, ",")
		// items[0] is always available and should contain the required value ...
		rw.app_server_ip = items[0]
	} else {
		rw.app_server_ip = "nil"
	}

	rw.crypto_key_utf8_encoded = utf8string.NewString(rw.crypto_key)

	k := fernet.MustDecodeKeys(rw.crypto_key_utf8_encoded.String())

	uaid, chid, pub_key, err := rw.decryptEndpoint(rw.Repad(rw.encrypted_endpoint), version, k)

	if err != nil { // Failed to decrypt the subscription data.
		l4g.Error("Failed to decrypt subscription data: '%s' appServerIP:%s", err, rw.app_server_ip)
		l4g.Info("responseCode:%v  uaid:NA  cepHostname: NA", http.StatusNotFound)
		if !rw.multicast{
			rw.stats.Increment("webpush.404")
			rw.w.WriteHeader(http.StatusNotFound) //404
		}
		resp_body = getResponseBody("102")
		fmt.Fprintf(rw.w, resp_body)

		return
	}


	if !rw.multicast{
		err = rw.extract_headers(rw.headers, rw.vapid)
	}
	if err != nil {
		l4g.Error("Error while extracting headers.'%s' responseCode:%v appServerIP:%s", err, http.StatusBadRequest, rw.app_server_ip)
		if !rw.multicast{
			rw.stats.Increment("webpush.400")
			rw.w.WriteHeader(http.StatusBadRequest) //400
		}
		return
	}

	l4g.Info("Message for uaid:%s from appServerIP:%s", rw.uaid, rw.app_server_ip)
	if len(rw.body) > _conf.Max_Msg_Length {
		l4g.Info("responseCode:%v Message of length %v is too large for uaid:%s ", http.StatusRequestEntityTooLarge, len(rw.body), rw.uaid)
		if !rw.multicast{
			rw.stats.Increment("webpush.413")
			rw.w.WriteHeader(http.StatusRequestEntityTooLarge) //413
		}
		resp_body = getResponseBody("104")
		fmt.Fprintf(rw.w, resp_body)

		return 
	}

	if rw.vapid && !rw.multicast{

		if rw.vapid_headers.CryptoKey == "" {

			l4g.Error("Received VAPID request without crypto-key header appServerIP:%s", rw.app_server_ip)
			if rw.debug {
				l4g.Info("responseCode:%v uaid:%s ", http.StatusUnauthorized, rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.401")
				rw.w.WriteHeader(http.StatusUnauthorized) //401
			}
			resp_body = getResponseBody("109")
			fmt.Fprintf(rw.w, resp_body)

			return

		}

		vapid.SetDebug(rw.debug)
		public_key_in_headers := vapid.GetLabel("p256ecdsa", rw.vapid_headers.CryptoKey)

		if public_key_in_headers == "" {

			l4g.Error("Error extracting p256ecdsa label (public key) from HTTP POST header appServerIP:%s", rw.app_server_ip)

			if rw.debug {

				l4g.Info("responseCode:%v uaid:%s ", http.StatusUnauthorized, rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.401")
				rw.w.WriteHeader(http.StatusUnauthorized) //401
			}
			resp_body = getResponseBody("114")
			fmt.Fprintf(rw.w, resp_body)
			
			return
		}

		processed_key, err := vapid.DecipherKey(public_key_in_headers)

		if err != nil {

			l4g.Error("Unable to decipher VAPID public key. '%s' appServerIP:%s", err, rw.app_server_ip)
			if rw.debug {

				l4g.Info("responseCode:%v uaid:%s ", http.StatusUnauthorized, rw.uaid)

			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.401")
				rw.w.WriteHeader(http.StatusUnauthorized) //401
			}

			resp_body = getResponseBody("115")
			fmt.Fprintf(rw.w, resp_body)
			
			return
		}

		if !rw.multicast {

			// Don't verify VAPID JWT on multicast here, becasue it has already been verified before in Muliticast() method. Verifying VAPID JWT is indeed the half way through verifying App server's identity (Which has been done before). After that, process to Verify public keys of each indivisual endpoint only and only if that endpoint is VAPID. If the endpoint is v1 skip that and notify the endpoint without doing such checking.
			err = vapid.Verify_AppServer_ID(rw.vapid_headers, public_key_in_headers, processed_key)
			if err != nil {
				l4g.Error("Unable to verify application server's identity. Routing request is refused. '%s' appServerIP:%s", err, rw.app_server_ip)
				if rw.debug {
					l4g.Info("responseCode:%v uaid:%s", http.StatusUnauthorized, rw.uaid)
				} else {

					l4g.Info("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
				}
				if !rw.multicast{
					rw.stats.Increment("webpush.401")
					rw.w.WriteHeader(http.StatusUnauthorized) //401
				}

				resp_body = getResponseBody("116")
				fmt.Fprintf(rw.w, resp_body)

				return
			}
		}

		err = vapid.Varify_PublicKey(rw.vapid_headers, pub_key, public_key_in_headers, processed_key)
		if err != nil {
			l4g.Error("Unable to verify application server's identity. Routing request is refused. '%s' appServerIP:%s", err, rw.app_server_ip)
			if rw.debug {
				l4g.Info("responseCode:%v uaid:%s", http.StatusUnauthorized, rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.401")
				rw.w.WriteHeader(http.StatusUnauthorized) //401
			}
			resp_body = getResponseBody("117")
			fmt.Fprintf(rw.w, resp_body)

			return
		}

	}

	if rw.debug && version == 2 {
		l4g.Info("Vapid  uaid:%s", rw.uaid)
		l4g.Info("Vapid  chid:%s", chid)
	}
	var router_type string
	var node_id string
	var current_month string //Defines the name of the message table on which all notifications are supposed to be saved for this device or this UAID.
	
	
	node_id, current_month, router_type, err = cassandra.GetDeviceData(rw.uaid, rw.debug, rw.session)

	if rw.debug {
		l4g.Info("Device with uaid:%s is connected on: '%s'", rw.uaid, node_id)
	}

	if err != nil && err.Error() != "not found"{
		l4g.Error("Failed to Get Device Data uaid='%s': '%s' appServerIP:%s", rw.uaid, err, rw.app_server_ip)		
		l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)		
		if !rw.multicast{
			rw.w.WriteHeader(http.StatusInternalServerError) //500
		}
		resp_body = getResponseBody("999")
		fmt.Fprintf(rw.w, resp_body)

		return
	}

	if node_id == "" {
		l4g.Info("responseCode:%v UAID not found uaid:%s appServerIP:%s", http.StatusGone, rw.uaid, rw.app_server_ip)
		if !rw.multicast{
			rw.w.WriteHeader(http.StatusGone) //410
		}
		resp_body = getResponseBody("118")
		fmt.Fprintf(rw.w, resp_body)

		return
	}


	// Is current_month entry is nil? If yes, then drop the user.

	if current_month == "" && version != 1 {

		l4g.Info("current_month entry for uaid:%s is not set, dropping user appServerIP:%s", rw.uaid, rw.app_server_ip)
		if cassandra.DropUser(rw.uaid, rw.session) != true {
			if !rw.multicast{
				rw.stats.Increment("webpush.500")
				rw.w.WriteHeader(http.StatusInternalServerError) //500
			}
			if rw.debug {
				l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid, node_id, rw.app_server_ip)
			}

			resp_body = getResponseBody("Error while droping user")
			fmt.Fprintf(rw.w, resp_body)

			resp_body = getResponseBody("999")
			fmt.Fprintf(rw.w, resp_body)

			return

		}

		if rw.debug {
			l4g.Info("responseCode:%v uaid:%s", http.StatusGone, rw.uaid)
		} else {

			l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusGone, rw.uaid, node_id)
		}
		if !rw.multicast{
			rw.stats.Increment("webpush.410")
			rw.w.WriteHeader(http.StatusGone) //410
		}
		resp_body = getResponseBody("103")
		fmt.Fprintf(rw.w, resp_body)

		return
	}

	// At this point, we would like to verify that the decoded chid (application ID) is already registered for that device or uaid. Because at some cases, it true that the device ID is found in router table, but does the application to which we are sending a message, actually registered for out push service? The answer can only be known by quering the message table of that device in order to get the list of registered applications (chids). If the channel ID is not found in the massage table of that device, then REP SHOULD return 401 GONE.

	if router_type == "webpush" {
		
		var found bool
		
		found,err = cassandra.ValidateWebpush(rw.session, rw.debug, chid, rw.uaid, current_month, router_type)
		

		if err != nil && err.Error() != "not found"{
			l4g.Error("Failed to validate CHID=%s uaid=%s: '%s' appServerIP:%s", chid, rw.uaid, err, rw.app_server_ip)		
			l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)		
			if !rw.multicast{
				rw.stats.Increment("webpush.500")
				rw.w.WriteHeader(http.StatusInternalServerError) //500
			}
			resp_body = getResponseBody("999")
			fmt.Fprintf(rw.w, resp_body)
			
			return
		}

		if found == false {

			l4g.Error("responseCode:%v  CHID not found chid:%s appServerIP:%s", http.StatusGone, chid, rw.app_server_ip)
			if !rw.multicast{
				rw.stats.Increment("webpush.410")
				rw.w.WriteHeader(http.StatusGone) //410
			}
			resp_body = getResponseBody("119")
			fmt.Fprintf(rw.w, resp_body)
			return

		}

	}

	// We don't want to make validation on message_month value, because for
	// SimplePush messages "current_month" column can be empty.(We don't save messages in SimplePush protocol, thus, we don't need to have message table for that.

	// If this is a SimplePush notification, then at this point, we have the required information to send it.

	if rw.protocol == 1 {

		if rw.debug {
			l4g.Info("Sending SimplePush message to client with uaid %s", rw.uaid)
		}

		rw.msg = rw.format_msg_data(rw.body)

		dataC = rw.finish_payloadC([]byte(rw.msg), chid)

		payloadBytes, err := json.Marshal(dataC)

		if err != nil {
			l4g.Error("[ERROR] Cannot marshal SimplePush notification body: '%s'", err)

			if rw.debug {
				l4g.Info("Offending simple push payload content: %v", dataC)
				return
			}
			rw.w.WriteHeader(http.StatusBadRequest) //400 
			return
		}

		body := bytes.NewReader(payloadBytes)
		var response_code int
		response_code, err = rw.notify_cep_node(node_id, body, rw.uaid)

		if err != nil {
			l4g.Error("[ERROR] Error Routing SimplePush notification. '%s' uaid:%s", err, rw.uaid)
			l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid, rw.app_server_ip)
			rw.w.WriteHeader(http.StatusInternalServerError) //500
			
			return

		}

		if rw.debug {

			l4g.Info("responseCode from CEP: %v uaid:%s", response_code, rw.uaid)

		}
		switch response_code {

		case 200:

			l4g.Info("Successful SimplePush notification delivery to CEP node uaid:%s", uaid)

			if rw.debug {

				l4g.Info("responseCode:%v uaid:%s", http.StatusCreated, rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.201")
				rw.w.WriteHeader(http.StatusCreated) //201
			}
		case 404:

			l4g.Info("Router miss. SimplePush wakeup signal failed. uaid:%s", rw.uaid)
			if rw.debug {

				l4g.Info("Device with uaid:%s is offline. Could not wake up application", rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.201")
				rw.w.WriteHeader(http.StatusCreated) //201
			}

		default:
			l4g.Error("Unexpected responseCode from CEP for a simplepush message: %v, uaid:%s appServerIP:%s", response_code, rw.uaid, rw.app_server_ip)
			if !rw.multicast{
				rw.stats.Increment("webpush.500")
				rw.w.WriteHeader(http.StatusInternalServerError) //500
			}
		}

		return
	}

	/*Based on the Notification POST request REP node has received, it will build the request that need to be sent to CEP. In the original Mozilla implementation,
	If the POST request does NOT have data inside (an empty notification), then the request to CEP will have a slightly different structure than of that when the POST request
	includes some message to be sent and the difference is that empty once does not have "data" or "headers" fields. That's the reason of creating two structures, Payload_A incase of sending data along with the notification, and Payload_B incase of sending an empty one.
	*/
	use_struct_A := false

	message_id, err := rw.generate_message_id(rw.uaid, rw.raw_chid)
	if err != nil {

		l4g.Error("[ERROR] Cannot generate message_id (message version). '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
		rw.w.WriteHeader(http.StatusInternalServerError) //500
		return
	}

	if len(rw.body) != 0 || rw.multicast {

		use_struct_A = true

		rw.msg = rw.format_msg_data(rw.body)

		dataA = rw.finish_payloadA(chid, message_id, rw.crypto_headers, rw.msg, rw.headers)

	} else {

		dataB = rw.finish_payloadB(chid, message_id, rw.headers)
	}

	var payloadBytes []byte

	if use_struct_A {

		payloadBytes, err = json.Marshal(dataA)
		if err != nil {
			l4g.Error("[ERROR] Cannot marshal notification body. '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
			l4g.Info("Offending input data was: '%v'", dataA)
			rw.w.WriteHeader(http.StatusInternalServerError) //500
			return
		}

	} else {

		payloadBytes, err = json.Marshal(dataB)
		if err != nil {
			l4g.Error("[ERROR] Cannot marshal notification body.'%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
			l4g.Info("Offending input data was: '%v'", dataB)
			rw.w.WriteHeader(http.StatusBadRequest) //400
			return
		}
	}

	rw.sorted_key = rw.sort_key(message_id, uaid, chid)
	notification := rw.create_wpush_offline_notif(string(message_id))

	body := bytes.NewReader(payloadBytes)
	var response_code int
	response_code, err = rw.notify_cep_node(node_id, body, rw.uaid)
	if err != nil {

		l4g.Error("[ERROR] Error Routing  WebPush notification.%s uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
		if rw.debug {
			l4g.Info("Notification object to be saved is: %v", notification)
			l4g.Info("Message for uaid:%s will be save in: %s", rw.uaid, current_month)

		}

		if notification.TTL != 0 { // Save only if message TTL value is higher than 0 seconds.
			
			if rw.is_limited { // Check for limitation only if this feature is enabled.
				if rw.is_exceeded(node_id) {
					l4g.Error("[ERROR] TTL exceed limit. uaid:%s appServerIP:%s", rw.uaid, rw.app_server_ip)
					return

				}
			}

			err = cassandra.StoreNotification(notification, current_month)
			
			if err != nil {

				l4g.Error("Error while saving message into DB: '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid, node_id, rw.app_server_ip)
				if !rw.multicast{
					rw.stats.Increment("webpush.500")
					rw.w.WriteHeader(http.StatusInternalServerError) //500
				}
				resp_body = getResponseBody("999")
				fmt.Fprintf(rw.w, resp_body)
				
				return
			}

			l4g.Info("Message saved successfully responseCode:%v uaid:%s ", http.StatusCreated, rw.uaid)
		}
		if !rw.multicast{
			rw.stats.Increment("webpush.201")
			rw.w.WriteHeader(http.StatusCreated) //201
		}
		return
	}

	if rw.debug {

		l4g.Info("responseCode from CEP: %v  uaid:%s", response_code, rw.uaid)

	}

	switch response_code {
	case 200:
		if version != 2 {

			//l4g.Info("Successful WebPush V1 notification delivery to CEP node uaid:%s", rw.uaid)
			if rw.debug {

				l4g.Info("responseCode:%v uaid:%s ", http.StatusCreated, rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.201")
				rw.w.WriteHeader(http.StatusCreated) //201
			}

		} else {

			//l4g.Info("Successful VAPID V2 WebPush notification delivery to CEP node uaid:%s", rw.uaid)
			if rw.debug {

				l4g.Info("responseCode:%v uaid:%s ", http.StatusCreated, rw.uaid)
			} else {

				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
			}
			if !rw.multicast{
				rw.stats.Increment("webpush.201")
				rw.w.WriteHeader(http.StatusCreated) //201
			}

		}

	case 404:

		if rw.debug {

			l4g.Info("Device with uaid:%s is offline. Message will be saved for later delivery", rw.uaid)
			l4g.Info("Router miss uaid:%s", rw.uaid)
		}
		if rw.debug {
			l4g.Info("Notification object to be saved is: %v", notification)
			l4g.Info("Message for device: %s will be save in: %s", rw.uaid, current_month)

		}

		if notification.TTL != 0 { // Save only if message TTL value is higher than 0 seconds.
			if rw.is_limited { // Check for limitation only if this feature is enabled.
				if rw.is_exceeded(node_id) {
					l4g.Error("[ERROR] TTL exceed limit. uaid:%s appServerIP:%s", rw.uaid, rw.app_server_ip)
					return

				}
			}
			
			err = cassandra.StoreNotification(notification, current_month)

			if err != nil {

				l4g.Error("Error while saving message into DB: '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
				l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)
				if !rw.multicast{
					rw.stats.Increment("webpush.500")
					rw.w.WriteHeader(http.StatusInternalServerError) //500
				}
				resp_body = getResponseBody("999")
				fmt.Fprintf(rw.w, resp_body)

				return
			}
			l4g.Info("Message saved successfully responseCode:%v uaid:%s ", http.StatusCreated, rw.uaid)
		}
		if !rw.multicast{
			rw.stats.Increment("webpush.201")
			rw.w.WriteHeader(http.StatusCreated) //201
		}
		
	default:

		l4g.Error("Unexpected response code from CEP for a Webpush message: %v, uaid:%s appServerIP:%s", response_code, rw.uaid, rw.app_server_ip)
		if !rw.multicast{
			rw.stats.Increment("webpush.500")
			rw.w.WriteHeader(http.StatusInternalServerError) //500
		}
	}

	return

}

func (rw *T_RouterWorker) extractSubscription(subscription string) (uaid string, chid string) {

	if rw.debug {

		l4g.Info("Exracting chid and uaid from the received CassandraDB result object")

	}

	uaid = subscription[0:32]
	unformatted_chid := subscription[32:64]
	rw.raw_chid = unformatted_chid
	chid = rw.format_id(unformatted_chid)
	return

}

func (rw *T_RouterWorker) notify_cep_node(node_id string, body io.Reader, uaid string) (response_code int, err error) {

	url := node_id + "/push/" + uaid
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	response_code = resp.StatusCode

	defer resp.Body.Close()

	return

}

func (rw *T_RouterWorker) format_id(chid string) string {

	/*

		This method will format extracted channelID from CassandraDB to be UUAID standard.

		Example:

		Input: 282c8b9184044db89f942c5972d0e55a

		Output: 282c8b91-8404-4db8-9f94-2c5972d0e55a


	*/

	if rw.debug {

		l4g.Info("Formating channelID :%v", chid)

	}

	var formated_chid string
	var first_part string
	for pos, char := range chid {

		first_part = first_part + string(char)

		if pos == 7 || pos == 11 || pos == 15 || pos == 19 {

			first_part = first_part + "-"
			formated_chid = formated_chid + first_part
			first_part = ""

		}

		if pos == 31 {

			formated_chid = formated_chid + first_part

		}

	}

	return formated_chid

}

// This method will take the received header object from the  message's server POST request which is of type http.headers and use it to fill Headers struct. What Push REP does, is that it passes the same headers it receives from message server to be used later in the PUT request. This will also return the crypto key header used in VAPID messages.

var crypto_key_req_labels_v2 = []string{"dh","p256ecdsa"}
var crypto_key_req_labels_v1 = []string{"dh"}

var encryption_key_req_labels = []string{"salt"}

func (rw *T_RouterWorker) extract_headers(request_headers http.Header, vapid bool) (err error) {
	var resp_body string
	if vapid {

		rw.vapid_headers.Encoding = request_headers.Get("Content-Encoding")
		rw.vapid_headers.Encryption = request_headers.Get("Encryption")
		rw.vapid_headers.Authorization = request_headers.Get("Authorization")

		rw.vapid_headers.CryptoKey = request_headers.Get("Crypto-Key")

		if len(rw.vapid_headers.CryptoKey) != 0 {
			fmt.Println("Call 1")
			rw.vapid_headers.CryptoKey, err = utils.Sanitize_Header(rw.vapid_headers.CryptoKey, crypto_key_req_labels_v2)

			if err != nil {
				resp_body = getResponseBody("101")
				fmt.Fprintf(rw.w, resp_body)
				return err

			}
		}

		if len(rw.vapid_headers.Encryption) != 0{
			fmt.Println("Call 2")
			rw.vapid_headers.Encryption, err = utils.Sanitize_Header(rw.vapid_headers.Encryption, encryption_key_req_labels)

			if err != nil {
				resp_body = getResponseBody("101")
				fmt.Fprintf(rw.w, resp_body)
				return err
			}
		}

	}

	rw.topic = request_headers.Get("Topic")

	str_ttl := request_headers.Get("ttl")

	if str_ttl != "" {

		rw.ttl, err = strconv.Atoi(request_headers.Get("ttl"))
		if err != nil {
			resp_body = getResponseBody("112")
			fmt.Fprintf(rw.w, resp_body)
			return err

		}
	}

	rw.crypto_headers.Encoding = request_headers.Get("Content-Encoding")
	rw.crypto_headers.Encryption = request_headers.Get("Encryption")
	rw.crypto_headers.CryptoKey = request_headers.Get("Crypto-Key")

	if len(rw.crypto_headers.CryptoKey) != 0 {
		
		rw.crypto_headers.CryptoKey, err = utils.Sanitize_Header(rw.crypto_headers.CryptoKey, crypto_key_req_labels_v1)

		if err != nil {
			resp_body = getResponseBody("101")
			fmt.Fprintf(rw.w, resp_body)
			return err
		}
	}

	if len(rw.crypto_headers.Encryption) != 0 && !rw.multicast { //Not needed for multicast
		
		rw.crypto_headers.Encryption, err = utils.Sanitize_Header(rw.crypto_headers.Encryption, encryption_key_req_labels)
		if err != nil {
			resp_body = getResponseBody("101")
			fmt.Fprintf(rw.w, resp_body)
			return err
		}
	}

	rw.topic = request_headers.Get("Topic")

	return nil

}

func (rw *T_RouterWorker) format_msg_data(data []byte) (msg string) {

	if rw.debug {

		l4g.Info("Formating received message.")

	}

	msg = base64.RawURLEncoding.EncodeToString(data)

	return msg

}

func (rw *T_RouterWorker) generate_message_id(uaid string, chid string) (message_id []byte, err error) {
	/*
	   message_id serves a complex purpose. It's returned as the Location header
	   value so that an application server may delete the message. It's used as
	   part of the non-versioned sort-key. Due to this, its an encrypted value
	   that contains the necessary information to derive the location of this
	   precise message in the appropriate message table

	       """Generate a message-id suitable for accessing the message

	       For topic messages, a sort_key version of 01 is used, and the topic
	       is included for reference:

	           Encrypted('01' : uaid.hex : channel_id.hex : topic)

	       For topic messages, a sort_key version of 02 is used:

	           Encrypted('02' : uaid.hex : channel_id.hex : timestamp)

	       For legacy non-topic messages, no sort_key version was used and the
	       message-id was:

	       Encrypted('m' : uaid.hex : channel_id.hex)
	*/

	if rw.debug {

		l4g.Info("Generating message_id (message version)")

	}

	var msg_key string

	str_sortkey_timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	str_sortkey_timestamp = str_sortkey_timestamp[:len(str_sortkey_timestamp)-3] //Python autopush time-stamp format

	rw.sortkey_timestamp, err = strconv.Atoi(str_sortkey_timestamp)
	if err != nil {

		return nil, err
	}

	if rw.topic != "" {

		msg_key = "01" + ":" + uaid + ":" + chid + ":" + rw.topic

	} else if LEGACY != "" {

		msg_key = "m" + ":" + uaid + ":" + chid

	} else {

		msg_key = "02" + ":" + uaid + ":" + chid + ":" + strconv.Itoa(rw.sortkey_timestamp)

	}

	key := fernet.MustDecodeKeys(rw.crypto_key_utf8_encoded.String())
	message_id, err = fernet.EncryptAndSign([]byte(msg_key), key[0])

	return
}

func (rw *T_RouterWorker) finish_payloadA(chid string, message_id []byte, headers model.Headers, msg string, request_headers http.Header) (dataA *model.Payload_A) {

	dataA = &model.Payload_A{

		ChannelID: chid,
		Version:   string(message_id),
		TTL:       rw.ttl,
		Topic:     []byte(rw.topic),
		Timestamp: time.Now().Unix(),
		Data:      msg,
		Headers:   headers,
	}

	if rw.debug {

		l4g.Info("Chid: %s", dataA.ChannelID)

	}
	if rw.debug {

		l4g.Info("Message to be sent: %s  to client with uaid:%s", dataA.Data, rw.uaid)

	}

	if rw.debug {

		l4g.Info("Headers to be used: %s", dataA.Headers)

	}

	return

}

func (rw *T_RouterWorker) finish_payloadB(chid string, message_id []byte, request_headers http.Header) (dataB *model.Payload_B) {

	dataB = &model.Payload_B{

		ChannelID: chid,
		Version:   string(message_id),
		TTL:       rw.ttl,
		Topic:     []byte(rw.topic),
		Timestamp: time.Now().Unix(),
	}

	if rw.debug {

		l4g.Info("Chid: %s", dataB.ChannelID)

	}

	return

}

func (rw *T_RouterWorker) finish_payloadC(data []byte, chid string) (dataC *model.Payload_C) {

	str_timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	if len(str_timestamp) > 10 {

		str_timestamp = str_timestamp[0:10]

	}

	int_timestamp, err := strconv.Atoi(str_timestamp)

	if err != nil {

		l4g.Warn("CAUTION: Failed to convert timestamp to int for simplepush message. Message may fail at CEP side.")

		// No need to return in such case.
	}

	dataC = &model.Payload_C{

		Version:   int_timestamp,
		Data:      data,
		ChannelID: chid,
	}

	if rw.debug {

		l4g.Info("Chid: %s", dataC.ChannelID)

	}

	return

}

// create_webpush_notification function takes a required parameter to construct a WebPushNotification object described in "model" package. WebPushNotification object will then be used as a parameter to StoreMessage()  method that will store offline messages (Incase CEP response is anything other than 200).
func (rw *T_RouterWorker) create_wpush_offline_notif(message_id string) (WebPushNotification *model.WebPushNotification) {

	if rw.debug {

		l4g.Info("Creating WebPushNotification object.")

	}

	var wpush_headers map[string]string
	wpush_headers = make(map[string]string)

	wpush_headers["crypto_key"] = rw.crypto_headers.CryptoKey
	wpush_headers["encoding"] = rw.crypto_headers.Encoding
	wpush_headers["encryption"] = rw.crypto_headers.Encryption

	str_timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	if len(str_timestamp) > 10 {

		str_timestamp = str_timestamp[0:10]

	}
	WebPushNotification = &model.WebPushNotification{

		Uaid:      rw.uaid,
		TTL:       rw.ttl,
		Data:      rw.msg,
		Headers:   wpush_headers,
		TimeStamp: str_timestamp,
		UpdateID:  message_id,
		Session:   rw.session,
		SortedKey: rw.sorted_key,
		Debug:     rw.debug,
	}

	return

}

func (rw *T_RouterWorker) create_simplepush_notification(data []byte, Chid string) (SimplePushNotification *model.SimplePushNotification) {

	if rw.debug {

		l4g.Info("Creating SimplePushNotification object.")

	}

	// In Mozilla's original implementation, SimplePush message's version is a timestamp.

	str_timestamp := strconv.Itoa(int(time.Now().UnixNano()))

	if len(str_timestamp) >= 10 {

		str_timestamp = str_timestamp[0:10]

	}

	SimplePushNotification = &model.SimplePushNotification{

		Version:   str_timestamp,
		Data:      "",
		ChannelID: Chid,
	}

	return

}

func (rw *T_RouterWorker) sort_key(message_id []byte, uaid string, chid string) (sorted_key string) {

	/* Return an appropriate sort_key for this notification

	   For new messages:

	       02:{sortkey_timestamp}:{chid}

	   For topic messages:

	       01:{chid}:{topic}

	   Old format for non-topic messages that is no longer returned:

	       {chid}:{message_id}
	*/

	if rw.topic != "" {
		sorted_key = fmt.Sprintf("01:{%s}:{%s}", rw.topic, uaid)
	} else if LEGACY != "" {

		sorted_key = fmt.Sprintf("{%s}:{%s}", chid, message_id)
	} else {
		//Created as late as possible when storing a message
		sorted_key = fmt.Sprintf("02:%v:%s", rw.sortkey_timestamp, chid)

	}

	if rw.debug {

		l4g.Info("Sort key used is: %s", sorted_key)
	}

	return

}

func (rw *T_RouterWorker) decryptEndpoint(encrypted_endpoint string, version int, crypto_key []*fernet.Key) (uaid string, chid string, pub_key string, err error) {

	if rw.debug && version == 1 {

		l4g.Info("Decrypting received endpoint.")
	}

	if rw.debug && version == 2 {

		l4g.Info("Decrypting received VAPID endpoint.")
	}

	decoded_endpoint := fernet.VerifyAndDecrypt([]byte(encrypted_endpoint), TTL_DURATION, crypto_key)

	if decoded_endpoint == nil {

		decryption_err := errors.New("[Error]: Invalid Endpoint.")

		return "", "", "", decryption_err
	}

	hex_encoded_endpoint := hex.EncodeToString([]byte(decoded_endpoint))

	if version == 1 {

		rw.uaid, chid = rw.extractSubscription(hex_encoded_endpoint)
		return

	} else {

		rw.uaid, chid, pub_key = rw.extractVapidSubscription(hex_encoded_endpoint)
		return

	}

}

func (rw *T_RouterWorker) extractVapidSubscription(subscription string) (uaid string, chid string, pub_key string) {

	if rw.debug {

		l4g.Info("Extracting chid, uaid and application server's public key from the received VAPID subscription data.")

	}

	if rw.debug {

		l4g.Info("Decrypted Vapid endpoint information: %s ", subscription)

	}

	uaid = subscription[0:32]
	unformatted_chid := subscription[32:64]
	rw.raw_chid = unformatted_chid
	chid = rw.format_id(unformatted_chid)
	pub_key = subscription[64:]
	return

}

func (rw *T_RouterWorker) Repad(subscription string) (paded_str string) {
	/*
		Adds padding to strings for base64 decoding.
		base64 encoding requires 'padding' to 4 octet boundries. 'padding' is a '='. So a string like 'abcde' is 5 octets, and would need to be padded to 'abcde==='.
	*/
	if rw.debug {

		l4g.Info("Repading received encrypted endpoint.")

	}

	padings := len(subscription) % 4

	if padings != 0 {

		for i := 0; i < padings; i++ {

			subscription = subscription + "="

		}

	}

	paded_str = subscription

	return

}

// Check if the maximum number of messages per UAID has reached or not.
func (rw *T_RouterWorker) is_exceeded(node_id string) bool {

	exceeded, err := counterdb.Reached_Msg_Max(rw.uaid, rw.session, rw.debug)
	if err != nil {
		l4g.Error("Error while checking message counter for uaid:%s. '%s'", rw.uaid, err)
	}

	if exceeded {

		l4g.Info("Reached max msg count for uaid:%s ", rw.uaid)
		l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusTooManyRequests, rw.uaid,rw.app_server_ip)
		if !rw.multicast{
			rw.stats.Increment("webpush.429")
			rw.w.WriteHeader(http.StatusTooManyRequests) //429
		}
		resp_body := getResponseBody(fmt.Sprintf("120"))
		fmt.Fprintf(rw.w, resp_body)
		return true
	}

	return false

}

func getResponseBody(text string)(string) {	
	return fmt.Sprintf("{\"errno\":\"%s\"}",text)
}
