package routercore

import (
	"bytes"
	"encoding/json"
	//"errors"
	fernet "fernet-go"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/config"
	"empowerthings.com/autoendpoint/routercore/db/cassandra"
	"empowerthings.com/autoendpoint/routercore/db/redis/counterdb"
	"empowerthings.com/autoendpoint/model"
	"empowerthings.com/autoendpoint/vapid"
	"empowerthings.com/autoendpoint/simplepush"
	"empowerthings.com/autoendpoint/routercore/subscription"
	"empowerthings.com/autoendpoint/utils"
	//"empowerthings.com/autoendpoint/utils/uuid"
	"github.com/gocql/gocql"
	"github.com/statsd"
	"golang.org/x/exp/utf8string"
	"time"
)

type T_RouterWorker struct {
	encrypted_endpoint      string
	body                    []byte
	headers                 http.Header
	crypto_key              string
	debug                   bool
	session                 *gocql.Session
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

func NewRouterWorker(encrypted_endpoint string, request_body []byte, request_headers http.Header, crypto_key string, debug bool, db_session *gocql.Session, protocol int, vapid bool, statistics *statsd.Client, multicast bool, is_limited bool) *T_RouterWorker {

	var tmp T_RouterWorker
	tmp = T_RouterWorker{}

	tmp.encrypted_endpoint = encrypted_endpoint
	tmp.body = request_body
	tmp.headers = request_headers
	tmp.crypto_key = crypto_key
	tmp.debug = debug
	tmp.session = db_session
	tmp.protocol = protocol
	tmp.vapid = vapid
	tmp.stats = statistics
	tmp.multicast = multicast
	tmp.is_limited = is_limited

	return &tmp

}
func (rw *T_RouterWorker) RouteNotification() (response_code int, response_body string){
	var dataA *model.Payload_A // Webpush payload
	var dataB *model.Payload_B // An empty(Webpush,SimplePush,VAPID) payload
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
	uaid, chid, pub_key, err := subscription.DecryptSubscription(subscription.RepadSubscription(rw.encrypted_endpoint,rw.debug), version, k,rw.debug)
	//l4g.Info("Formating CHID: %s", raw_chid)
	//chid := uuid.FormatId(raw_chid,rw.debug)
	rw.uaid=uaid
	if err != nil { // Failed to decrypt the subscription data.
		l4g.Error("Failed to decrypt subscription data: '%s' appServerIP:%s", err, rw.app_server_ip)
		l4g.Info("responseCode:%v  uaid:NA  cepHostname: NA", http.StatusNotFound)
		rw.stats.Increment("webpush.404")
		return http.StatusNotFound,getResponseBody("102") //404
	}

	if !rw.multicast && rw.protocol!=1 { // For multicast message we don't need to extract headers because it is an EMPTY push message.
		err = rw.extract_headers(rw.headers, rw.vapid)
		if err != nil {
			l4g.Error("Error while extracting headers.'%s' responseCode:%v appServerIP:%s", err, http.StatusBadRequest, rw.app_server_ip)
			rw.stats.Increment("webpush.400")
			return http.StatusBadRequest, "" //400		
		}else{
			rw.ttl, err = strconv.Atoi(rw.headers.Get("ttl"))
			if err != nil {
				if err != nil {
					l4g.Error("Error while extracting TTL header of the multicast request.'%s' responseCode:%v appServerIP:%s", err, http.StatusBadRequest, rw.app_server_ip)
					return http.StatusBadRequest, "" //400
				}	
			}
		}
	}			
	l4g.Info("Message for uaid:%s from appServerIP:%s", rw.uaid, rw.app_server_ip)
	if len(rw.body) > _conf.Max_Msg_Length {
		l4g.Info("responseCode:%v Message of length %v is too large for uaid:%s ", http.StatusRequestEntityTooLarge, len(rw.body), rw.uaid)
		rw.stats.Increment("webpush.413")
		return http.StatusRequestEntityTooLarge, getResponseBody("104") //413
	}
	// Don't verify VAPID JWT on multicast here, becasue it has already been verified before in Muliticast() method. Verifying VAPID JWT is indeed the half way through verifying App server's identity (Which has been done before). After that, process to Verify public keys of each indivisual endpoint only and only if that endpoint is VAPID. If the endpoint is v1 skip that and notify the endpoint without doing such checking.
	if rw.vapid && !rw.multicast{
		if rw.vapid_headers.CryptoKey == "" {
			l4g.Error("Received VAPID request without crypto-key header appServerIP:%s", rw.app_server_ip)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			} else {
				l4g.Debug("responseCode:%v uaid:%s ", http.StatusUnauthorized, rw.uaid)
			}
			rw.stats.Increment("webpush.401")
			return http.StatusUnauthorized,getResponseBody("109") //401
		}
		vapid.SetDebug(rw.debug)
		public_key_in_headers := vapid.GetLabel("p256ecdsa", rw.vapid_headers.CryptoKey)
		if public_key_in_headers == "" {
			l4g.Error("Error extracting p256ecdsa label (public key) from HTTP POST header appServerIP:%s", rw.app_server_ip)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)		
			} else {
				l4g.Debug("responseCode:%v uaid:%s ", http.StatusUnauthorized, rw.uaid)
			}
			rw.stats.Increment("webpush.401")
			return http.StatusUnauthorized,getResponseBody("114") //401	
		}
		processed_key, err := vapid.DecipherKey(public_key_in_headers)
		if err != nil {
			l4g.Error("Unable to decipher VAPID public key. '%s' appServerIP:%s", err, rw.app_server_ip)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			} else {
				l4g.Debug("responseCode:%v uaid:%s ", http.StatusUnauthorized, rw.uaid)
			}
			rw.stats.Increment("webpush.401")
			return http.StatusUnauthorized, getResponseBody("115") //401		
		}
		err = vapid.Verify_AppServer_ID(rw.vapid_headers, public_key_in_headers, processed_key)
		if err != nil {
			l4g.Error("Unable to verify application server's identity. Routing request is refused. '%s' appServerIP:%s", err, rw.app_server_ip)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			} else {
				l4g.Debug("responseCode:%v uaid:%s", http.StatusUnauthorized, rw.uaid)
			}				
			rw.stats.Increment("webpush.401")
			return http.StatusUnauthorized,getResponseBody("116") //401
		}
		err = vapid.Verify_PublicKey(rw.vapid_headers, pub_key, public_key_in_headers, processed_key)
		if err != nil {
			l4g.Error("Unable to verify application server's identity. Routing request is refused. '%s' appServerIP:%s", err, rw.app_server_ip)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname: NA", http.StatusUnauthorized, rw.uaid)
			} else {
				l4g.Debug("responseCode:%v uaid:%s", http.StatusUnauthorized, rw.uaid)
			}
			rw.stats.Increment("webpush.401")
			return http.StatusUnauthorized,getResponseBody("117") //401
		}
	}
	if rw.debug && version == 2 {
		l4g.Debug("Vapid  uaid:%s", rw.uaid)
		l4g.Debug("Vapid  chid:%s", chid)
	}
	var router_type string
	var node_id string
	var current_month string //Defines the name of the message table on which all notifications are supposed to be saved for this device or this UAID.
	node_id, current_month, router_type, err = cassandra.GetDeviceData(rw.uaid, rw.debug, rw.session)
	if rw.debug {
		l4g.Debug("Device with uaid:%s is connected on: '%s'", rw.uaid, node_id)
	}
	if err != nil &&  strings.Compare(err.Error(), "not found") !=0{
		l4g.Error("Failed to Get Device Data uaid='%s': '%s' appServerIP:%s", rw.uaid, err, rw.app_server_ip)		
		l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)		
		return http.StatusInternalServerError, getResponseBody("999") //500
		
	}
	if node_id == "" {
		l4g.Info("responseCode:%v UAID not found uaid:%s appServerIP:%s", http.StatusGone, rw.uaid, rw.app_server_ip)
		return http.StatusGone, getResponseBody("118") //410
		
	}
	// Is current_month entry is nil? If yes, then drop the user.
	if len(current_month) == 0 && version != 1 {
		l4g.Info("current_month entry for uaid:%s is not set, dropping user appServerIP:%s", rw.uaid, rw.app_server_ip)
		if cassandra.DropUser(rw.uaid, rw.session) != true {
			if rw.debug {
				l4g.Debug("WARNING: could not drop user:%s from router table",rw.uaid)
			}  
		}
		if rw.debug {
			l4g.Debug("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusGone, rw.uaid, node_id)
		} else {
			l4g.Debug("responseCode:%v uaid:%s", http.StatusGone, rw.uaid)
		}	
		rw.stats.Increment("webpush.410")
		return http.StatusGone, getResponseBody("103") //410
		
	}
	// At this point, we would like to verify that the decoded chid (application ID) is already registered for that device or uaid. Because at some cases, it true that the device ID is found in router table, but does the application to which we are sending a message, actually registered for out push service? The answer can only be known by quering the message table of that device in order to get the list of registered applications (chids). If the channel ID is not found in the massage table of that device, then REP SHOULD return 401 GONE.
	if router_type == "webpush" {
		var found bool
		found,err = cassandra.ValidateWebpush(rw.session, rw.debug, chid, rw.uaid, current_month, router_type)
		if err != nil && err.Error() != "not found"{
			l4g.Error("Failed to validate CHID=%s uaid=%s: '%s' appServerIP:%s", chid, rw.uaid, err, rw.app_server_ip)		
			l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)		
			rw.stats.Increment("webpush.500")
			return http.StatusInternalServerError,getResponseBody("999") //500
			
		}
		if found == false {
			l4g.Error("responseCode:%v  CHID not found chid:%s appServerIP:%s", http.StatusGone, chid, rw.app_server_ip)
			rw.stats.Increment("webpush.410")
			return http.StatusGone,getResponseBody("119") //410
			
		}
	}
	// We don't want to make validation on message_month value, because for
	// SimplePush messages "current_month" column can be empty.(We don't save messages in SimplePush protocol, thus, we don't need to have message table for that.

	// If this is a SimplePush notification, then at this point, we have the required information to send it.
	if rw.protocol == 1 {
		if rw.debug {
			l4g.Debug("Sending SimplePush message to client with uaid %s", rw.uaid)
		}
		rw.msg = model.FormatNotificationData(rw.body, rw.debug)
		simple_push_notif := simplepush.CreateSimplePushNotif([]byte(rw.msg), chid,rw.debug)
		
		payloadBytes, err := json.Marshal(simple_push_notif)
		if err != nil {
			l4g.Error("[ERROR] Cannot marshal SimplePush notification body: '%s'", err)
			if rw.debug {
				l4g.Debug("Offending simple push payload content: %v", simple_push_notif)
			}
			return http.StatusInternalServerError,"" //500 
		}
		body := bytes.NewReader(payloadBytes)
		var response_code int
		response_code, err = rw.notify_cep_node(node_id, body, rw.uaid)
		if err != nil {
			l4g.Error("Error Routing SimplePush notification. '%s' uaid:%s", err, rw.uaid)
			l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid, rw.app_server_ip)
			return http.StatusInternalServerError,"" //500
		}
		if rw.debug {
			l4g.Debug("responseCode from CEP: %v uaid:%s", response_code, rw.uaid)
		}
		switch response_code {
		case 200:
			l4g.Info("Successful SimplePush notification delivery to CEP node uaid:%s", uaid)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
			} else {
				l4g.Debug("responseCode:%v uaid:%s", http.StatusCreated, rw.uaid)
			}
			rw.stats.Increment("webpush.201")
			return http.StatusCreated,"" //201
			
		case 404:
			l4g.Info("Router miss. SimplePush wakeup signal failed. uaid:%s", rw.uaid)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
				
			} else {
				l4g.Debug("Device with uaid:%s is offline. Could not wake up application", rw.uaid)
			}
			rw.stats.Increment("webpush.202")
			return http.StatusAccepted,"" //202
			
		default:
			l4g.Error("Unexpected responseCode from CEP for a simplepush message: %v, uaid:%s appServerIP:%s", response_code, rw.uaid, rw.app_server_ip)
			rw.stats.Increment("webpush.500")
			return http.StatusInternalServerError,"" //500
			
		}
	}
	/*Based on the Notification POST request REP node has received, it will build the request that need to be sent to CEP. In the original Mozilla implementation,
	If the POST request does NOT have data inside (an empty notification), then the request to CEP will have a slightly different structure than of that when the POST request
	includes some message to be sent and the difference is that empty once does not have "data" or "headers" fields. That's the reason of creating two structures, Payload_A incase of sending data along with the notification, and Payload_B incase of sending an empty one.
	*/
	use_struct_A := false

	str_sortkey_timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	str_sortkey_timestamp = str_sortkey_timestamp[:len(str_sortkey_timestamp)-3] //Python autopush time-stamp format
	sortkey_timestamp, err := strconv.Atoi(str_sortkey_timestamp)
	if err!=nil{
		l4g.Error("[ERROR] Cannot generate sortkey_timestamp. '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
		return http.StatusInternalServerError,"" //500
	}
	rw.sortkey_timestamp = sortkey_timestamp
	message_id, err := model.GenerateNotificationId(rw.uaid, rw.raw_chid, rw.crypto_key_utf8_encoded, rw.sortkey_timestamp,rw.topic, rw.debug)
	if err != nil {
		l4g.Error("[ERROR] Cannot generate message_id (message version). '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
		return http.StatusInternalServerError,"" //500
	}
	if len(rw.body) != 0 || rw.multicast {
		use_struct_A = true
		rw.msg = model.FormatNotificationData(rw.body, rw.debug)
		dataA = model.NewPayload_A(chid, message_id, rw.crypto_headers, rw.msg, rw.headers,rw.ttl, rw.topic,rw.uaid,rw.debug)
	} else {
		dataB = model.NewPayload_B(chid, message_id, rw.headers, rw.topic, rw.ttl,rw.debug)
	}
	var payloadBytes []byte
	if use_struct_A {
		payloadBytes, err = json.Marshal(dataA)
		if err != nil {
			l4g.Error("[ERROR] Cannot marshal notification body. '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
			l4g.Info("Offending input data was: '%v'", dataA)
			return http.StatusInternalServerError,"" //500
		}
	} else {
		payloadBytes, err = json.Marshal(dataB)
		if err != nil {
			l4g.Error("[ERROR] Cannot marshal notification body.'%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
			l4g.Info("Offending input data was: '%v'", dataB)
			return http.StatusInternalServerError,"" //500
		}
	}
	rw.sorted_key = model.SortKey(message_id, uaid, chid,rw.topic, rw.sortkey_timestamp,rw.debug)
	notification := model.CreateWebPushNotif(string(message_id),rw.crypto_headers,rw.uaid,rw.ttl,rw.msg,rw.session,rw.sorted_key,rw.debug)
	body := bytes.NewReader(payloadBytes)
	response_code, err = rw.notify_cep_node(node_id, body, rw.uaid)
	if err != nil {
		l4g.Error("[ERROR] Error Routing  WebPush notification.%s uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
		if rw.debug {
			l4g.Debug("Notification object to be saved is: %v", notification)
			l4g.Debug("Message for uaid:%s will be saved in: %s", rw.uaid, current_month)
		}
		if notification.TTL != 0 { // Save only if message TTL value is higher than 0 seconds.
			if rw.is_limited && _conf.Max_Msg > 0  { // Check for limitation only if this feature is enabled.
				if rw.is_exceeded(node_id) {
					l4g.Debug("Device has exceeded its limit. uaid:%s appServerIP:%s", rw.uaid, rw.app_server_ip)
					return http.StatusTooManyRequests,"" //429
				}
			}
			err = cassandra.StoreNotification(notification, current_month)
			if err != nil {
				l4g.Error("Error while saving message into DB: '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
				l4g.Info("responseCode:%v  uaid:%s  cepHostname:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid, node_id, rw.app_server_ip)
				rw.stats.Increment("webpush.500")
				return http.StatusInternalServerError, getResponseBody("999") //500
				
			}
			l4g.Info("Message saved successfully responseCode:%v uaid:%s ", http.StatusAccepted, rw.uaid)
		}
		rw.stats.Increment("webpush.202")
		return http.StatusAccepted,"" //202
		
	}
	if rw.debug {
		l4g.Debug("responseCode from CEP: %v  uaid:%s", response_code, rw.uaid)
	}
	switch response_code {
	case 200:
		if version != 2 {
			//l4g.Info("Successful WebPush V1 notification delivery to CEP node uaid:%s", rw.uaid)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
			} else {
				l4g.Debug("responseCode:%v uaid:%s ", http.StatusCreated, rw.uaid)
			}
			rw.stats.Increment("webpush.201")
			return http.StatusCreated,"" //201
			
		} else {

			//l4g.Info("Successful VAPID V2 WebPush notification delivery to CEP node uaid:%s", rw.uaid)
			if rw.debug {
				l4g.Debug("responseCode:%v  uaid:%s  cepHostname:%s", http.StatusCreated, rw.uaid, node_id)
				
			} else {
				l4g.Debug("responseCode:%v uaid:%s ", http.StatusCreated, rw.uaid)
			}
			rw.stats.Increment("webpush.201")
			return http.StatusCreated,"" //201
			
		}
	case 404:
		if rw.debug {
			l4g.Debug("Device with uaid:%s is offline. Message will be saved for later delivery in: %s", rw.uaid, current_month)
			l4g.Debug("Notification object to be saved for uaid:%s is: %v",rw.uaid, notification)
		}
		if notification.TTL != 0 { // Save only if message TTL value is higher than 0 seconds.
			if rw.is_limited { // Check for limitation only if this feature is enabled.
				if rw.is_exceeded(node_id) {
					l4g.Error("[ERROR] TTL exceed limit. uaid:%s appServerIP:%s", rw.uaid, rw.app_server_ip)
				}
			}	
			err = cassandra.StoreNotification(notification, current_month)
			if err != nil {
				l4g.Error("Error while saving message into DB: '%s' uaid:%s appServerIP:%s", err, rw.uaid, rw.app_server_ip)
				l4g.Info("responseCode:%v uaid:%s appServerIP:%s", http.StatusInternalServerError, rw.uaid,rw.app_server_ip)
				rw.stats.Increment("webpush.500")
				return http.StatusInternalServerError, getResponseBody("999")//500
			}
			l4g.Info("Message saved successfully responseCode:%v uaid:%s ", http.StatusAccepted, rw.uaid)
		}
		
		rw.stats.Increment("webpush.202")
		return http.StatusAccepted,"" //202
		
	default:
		l4g.Error("Unexpected response code from CEP for a Webpush message: %v, uaid:%s appServerIP:%s", response_code, rw.uaid, rw.app_server_ip)		
		rw.stats.Increment("webpush.500")
		return http.StatusInternalServerError,"" //500
		
	}
	return http.StatusInternalServerError,"" //500
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
// This method will take the received header object from the  message's server POST request which is of type http.headers and use it to fill Headers struct. What Push REP does, is that it passes the same headers it receives from message server to be used later in the PUT request. This will also return the crypto key header used in VAPID messages.

var crypto_key_req_labels_v2 = []string{"dh","p256ecdsa"}
var crypto_key_req_labels_v1 = []string{"dh"}
var encryption_key_req_labels = []string{"salt"}

func (rw *T_RouterWorker) extract_headers(request_headers http.Header, vapid bool) (err error) {
	if vapid {
		rw.vapid_headers.Encoding = request_headers.Get("Content-Encoding")
		rw.vapid_headers.Encryption = request_headers.Get("Encryption")
		rw.vapid_headers.Authorization = request_headers.Get("Authorization")

		rw.vapid_headers.CryptoKey = request_headers.Get("Crypto-Key")

		if len(rw.vapid_headers.CryptoKey) != 0 {
			rw.vapid_headers.CryptoKey, err = utils.Sanitize_Header(rw.vapid_headers.CryptoKey, crypto_key_req_labels_v2)
			if err != nil {
				return err
			}
		}
		if len(rw.vapid_headers.Encryption) != 0{
			rw.vapid_headers.Encryption, err = utils.Sanitize_Header(rw.vapid_headers.Encryption, encryption_key_req_labels)
			if err != nil {
				return err
			}
		}
	}
	rw.topic = request_headers.Get("Topic")
	str_ttl := request_headers.Get("ttl")
	if str_ttl != "" {
		rw.ttl, err = strconv.Atoi(request_headers.Get("ttl"))
		if err != nil {
			return err
		}
	}
	rw.crypto_headers.Encoding = request_headers.Get("Content-Encoding")
	rw.crypto_headers.Encryption = request_headers.Get("Encryption")
	rw.crypto_headers.CryptoKey = request_headers.Get("Crypto-Key")

	if len(rw.crypto_headers.CryptoKey) != 0 {	
		rw.crypto_headers.CryptoKey, err = utils.Sanitize_Header(rw.crypto_headers.CryptoKey, crypto_key_req_labels_v1)
		if err != nil {
			return err
		}
	}
	if len(rw.crypto_headers.Encryption) != 0 && !rw.multicast { //Not needed for multicast
		
		rw.crypto_headers.Encryption, err = utils.Sanitize_Header(rw.crypto_headers.Encryption, encryption_key_req_labels)
		if err != nil {
			return err
		}
	}
	rw.topic = request_headers.Get("Topic")
	return nil
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
		rw.stats.Increment("webpush.429")
		return true
	}
	return false
}
func getResponseBody(text string)(string) {	
	return fmt.Sprintf("{\"errno\":\"%s\"}",text)
}
