package push_handlers

import (

	//	"empowerthings.com/cumulis/config"
	"io/ioutil"
	"net/http"
	"runtime"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/routercore"
	"empowerthings.com/autoendpoint/model"
	"empowerthings.com/autoendpoint/routercore/multicast"
	"empowerthings.com/autoendpoint/vapid"
	"empowerthings.com/autoendpoint/webpush"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/statsd"
	"fmt"
	"os"
	"strings"
)

var (
	crypto_key string
	debug bool
	db_session *gocql.Session
	payload_max_size int
	stats *statsd.Client
	is_limited bool
)
func Init() {

	is_limited = false
}
func WebPushHandler(w http.ResponseWriter, r *http.Request) {
	//Only POST is allowed on this API.
	if strings.Compare(r.Method, "POST")!=0 {
		l4g.Error("REP node received a wrong HTTP verb: %s on V1  API", r.Method)
		l4g.Debug("Response code to application server: %v ", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	/* Record statistics
	   [1] Counters

	   [2] Gauges

	   [3] Timers
	*/
	// hostname, _ := os.Hostname()
	// stats.Gauge(hostname+".webpush.v1", runtime.NumGoroutine())

	// t := stats.NewTiming()

	// t.Send("webpush.v1")

	if debug {
		l4g.Debug("REP node has received a request to route a WebPush notification.")
	}

	vars := mux.Vars(r)
	api_ver := vars["api_ver"]
	// Validating the {api_ver} variable.Should be either "v1" or "v2".
	if api_ver != "v1" && api_ver != "v2" {
		if debug {
			l4g.Debug("REP node has received an Invalid API version:  %s", api_ver)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	encrypted_endpoint := vars["encrypted_endpoint"]
	if debug {
		l4g.Debug("Subscription data to be decrypted is: %s", encrypted_endpoint)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l4g.Error("Error reading webpush notification body: '%s'", err)
		l4g.Debug("Response code to application server: %v ", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if strings.Compare(api_ver, "v1")==0 {
		if debug {
			l4g.Debug("V1 API is called. Received header is: %v", r.Header)
		}
		stats.Increment("webpush.v1")
		response_code,response_body:=webpush.Route_Webpush_Message(encrypted_endpoint, body, r.Header, crypto_key, debug, db_session,stats, is_limited)
		w.WriteHeader(response_code)
		fmt.Fprintf(w, response_body)
		return
	} else {
		stats.Increment("webpush.v2") //VAPID
		router_worker := routercore.NewRouterWorker(encrypted_endpoint, body, r.Header, crypto_key, debug, db_session, 1, true, stats, false, is_limited)
		response_code,response_body:=router_worker.RouteNotification()
		w.WriteHeader(response_code)
		fmt.Fprintf(w, response_body)
		return
	}
}
func SimplePushHandler(w http.ResponseWriter, r *http.Request) {
	//Only PUT is allowed on this API.
	if strings.Compare(r.Method,"PUT")!=0 {
		l4g.Error("REP node received a wrong HTTP verb: %s on SimplePush /spush API", r.Method)
		l4g.Debug("Response code to application server: %v ", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	stats.Increment("simplepush.v1")
	hostname, _ := os.Hostname()
	stats.Gauge(hostname+".simplepush.v1", runtime.NumGoroutine())
	t := stats.NewTiming()
	t.Send("simplepush.v1")
	if debug {
		l4g.Debug("REP node has received a request to route a SimpePush notification.")
	}
	vars := mux.Vars(r)
	encrypted_endpoint := vars["encrypted_endpoint"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l4g.Error("Error reading SimplePush notification body: '%s'", err)
		l4g.Debug("Response code to application server: %v ", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	router_worker := routercore.NewRouterWorker(encrypted_endpoint, body, r.Header, crypto_key, debug, db_session, 1, false, stats, false, is_limited)
	response_code,response_body:= router_worker.RouteNotification()
	w.WriteHeader(response_code)
	fmt.Fprintf(w, response_body)
	return
}
func GoFlip2SPushHandler(w http.ResponseWriter, r *http.Request) {
	//Only PUT is allowed on this API.
	if strings.Compare(r.Method ,"PUT")!=0 {
		l4g.Error("REP node received a wrong HTTP verb: %s on GoFLIP2 /update API", r.Method)
		l4g.Debug("Response code to application server: %v ", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	stats.Increment("simplepush.v1")
	hostname, _ := os.Hostname()
	stats.Gauge(hostname+".simplepush.v1", runtime.NumGoroutine())
	t := stats.NewTiming()
	t.Send("simplepush.v1")
	if debug {
		l4g.Debug("REP node has received a request to route a GoFlip2 SimpePush notification.")
	}
	vars := mux.Vars(r)
	encrypted_endpoint := vars["encrypted_endpoint"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l4g.Error("Error reading SimplePush notification body: '%s'", err)
		l4g.Debug("Response code to application server: %v ", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	router_worker := routercore.NewRouterWorker(encrypted_endpoint, body, r.Header, crypto_key, debug, db_session,1, false, stats, false, is_limited)
	response_code,response_body:=router_worker.RouteNotification()
	w.WriteHeader(response_code)
	fmt.Fprintf(w, response_body)
	return
}
func Route_Vapid_Webpush_Message(w http.ResponseWriter, r *http.Request) {
	//Only POST is allowed on this API.
	if strings.Compare(r.Method, "POST")!=0 {
		l4g.Error("REP node received a wrong HTTP verb: %s on V2 API", r.Method)
		l4g.Debug("Response code to application server: %v ", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if debug {
		l4g.Debug("REP node has received a request to route VAPID WebPush notification.")
	}
	stats.Increment("webpush.v2")
	hostname, _ := os.Hostname()
	stats.Gauge(hostname+".webpush.v2", runtime.NumGoroutine())
	t := stats.NewTiming()
	t.Send("webpush.v2")
	vars := mux.Vars(r)
	encrypted_endpoint := vars["encrypted_endpoint"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l4g.Error("Error reading VAPID notification body: '%s'", err)
		l4g.Debug("Response code to application server: %v ", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	router_worker := routercore.NewRouterWorker(encrypted_endpoint, body, r.Header, crypto_key, debug, db_session,2, true, stats, false, is_limited)
	response_code,response_body:=router_worker.RouteNotification()
	w.WriteHeader(response_code)
	fmt.Fprintf(w, response_body)
	return
}
func SetCryptoKey(encryption_key string) {
	crypto_key = encryption_key
}
func SetDebug(debug_log_level bool) {
	debug = debug_log_level
}
func SetSession(session *gocql.Session) {
	db_session = session
}
func SetPayloadSize(size int) {
	payload_max_size = size
}
func SetStats(statistics *statsd.Client) {
	stats = statistics
}
func EnableLimitation() {
	is_limited = true
}
func Get_Health_Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if debug {
		l4g.Debug("Responded 200 to health GET")
	}
}
func Multicast(w http.ResponseWriter, r *http.Request) {
	/*
		Kaios Push Multicast API.

		A list of endpoint will be included in the body of the request.

		Server will iterate over the list and will notify each application one by one.

		Multicast was added as a non-standard feature to our push server, that is used to notify a subset of the total push subscriptions (endpoints), and it does not include sending a message in the body.

		It was implemented in our server by adding an additional API  /bwpush/

		The general form of the multicast request would be:

		curl -X POST -H 'crypto-key: keyid=p256dh;dh=TyEtRjTZ9b_25iziK4RkRbxKlNptjNOoMvNtEcFY7vY=,p256ecdsa=BJdLQlHn8RxNWN97P4EPN1E8gTXmyt076dMozixe_4KzfVFVHkqdE60_a0MKYt2-fCwoPnQhXiuMQA7JiLdag2g' -H 'encryption: keyid=p256dh;salt=susKL-fdFoKur1aTjpJ51g' -H 'content-encoding: aesgcm' -H 'TTL: 60' -d '{"ep":[<endpoint1>,<endpoint2>,<endpoint3>,…,<endpoint n>], "msg":""}'  push.test.kaiostech.com:8443/bwpush/

		After receiving such request,  push server will simply loop over each endpoint included in the body, and will notify its corresponding application by sending a notification with no data.

		The array of endpoints in the body above would be an array of the last part of each endpoint (which is the part after the last "/" and which starts with gAAA....).

		There are three types of multicast requests:

		1- V1 Munticast request. Which is the multicast reuest the has v1 endpoints in its body.

		2- V2 Munticast request. Which is the multicast reuest the has VAPID endpoints in its body.

		3- Mix of  V1 and V2 broacast requests.

		If there is an N number of v1 endpoints and M number of v2 endpoints, then:

		Number of workers created = n + m  ( a worker per endpoint)

		Number of JWT checkings = 1 checkig only before entering the for loop.

		Number of public key comparisons  = m

	*/

	//Only POST is allowed on this API.
	if strings.Compare(r.Method ,"POST")!=0 {
		l4g.Error("REP node received a wrong HTTP verb: %s on multicast /bwpush API", r.Method)
		l4g.Debug("Response code to application server: %v ", http.StatusMethodNotAllowed)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if debug {
		l4g.Debug("REP node has received multicast request.")
	}
	stats.Increment("multicast")
	hostname, _ := os.Hostname()
	stats.Gauge(hostname+".multicast", runtime.NumGoroutine())
	t := stats.NewTiming()
	t.Send("multicast")
	// Apply payload limitation (in bytes).
	if debug {
		l4g.Debug("Maximum multicast message body size is: %v", int64(payload_max_size))
	}
	r.Body = http.MaxBytesReader(w, r.Body, int64(payload_max_size))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l4g.Error("Error reading body of multicast request: %v", err)
		l4g.Debug("Response code to application server: %v ", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if body == nil {
		l4g.Error("Wrong multicast request. Body is empty.")
		l4g.Debug("Response code to application server: %v ", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	multicast_list, message, err := multicast.Extract_List(body, debug)
	if err != nil {
		l4g.Error("Error extracting multicast list. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// NOTE: Message SHOULD be base64 format.
	if debug {
		l4g.Debug("Verifying Multicast application server's ID")
	}
	if debug {
		l4g.Debug("Multicast message to be sent is: %s ", message)
	}
	vapid.SetDebug(debug)
	vapid_headers := &model.Vapid_Headers{
		Encoding:      r.Header.Get("Content-Encoding"),
		Encryption:    r.Header.Get("Encryption"),
		CryptoKey:     r.Header.Get("Crypto-Key"),
		Authorization: r.Header.Get("Authorization"),
	}
	public_key_in_headers := vapid.GetLabel("p256ecdsa", vapid_headers.CryptoKey)
	if public_key_in_headers == "" {
		l4g.Error("Error extracting p256ecdsa label (public key) from HTTP Multicast  POST body.No VAPID public key found in crypto-key header")
		l4g.Debug("Response code to application server: %v ", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized) //401
		return
	}
	processed_key, err := vapid.DecipherKey(public_key_in_headers)
	if err != nil {
		l4g.Error("Unable to decipher VAPID public key. '%s'", err)
		l4g.Debug("Response code to application server: %v ", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized) //401
		return
	}
	err = vapid.Verify_AppServer_ID(*vapid_headers, public_key_in_headers, processed_key)
	if err != nil {
		l4g.Error("Unable to verify Multicast app server's ID. %s", err)
		l4g.Debug("Response code to application server: %v ", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized) //401
		return
	}
	if debug {
		l4g.Debug("Multicast app server identity is varified.")
	}
	if debug {
		l4g.Debug("Multicasting")
	}
	//Looping over extracted multicast list.
	if debug {
		l4g.Debug(" Multicast response code to application server: %v ", http.StatusCreated)
	}
	if debug {
		l4g.Debug(" Multicast response code to application server: %v ", http.StatusCreated)
	}
	w.WriteHeader(http.StatusCreated)
	for _, encrypted_endpoint := range multicast_list {
		if len(encrypted_endpoint) == 183 { // vapid endpoint
			router_worker := routercore.NewRouterWorker(encrypted_endpoint, message, r.Header, crypto_key, debug, db_session, 1, true, stats, true, is_limited)
			 _,_=router_worker.RouteNotification()
		} else { //Webpush endpoint
			router_worker := routercore.NewRouterWorker(encrypted_endpoint, message, r.Header, crypto_key, debug, db_session,2, false, stats, true, is_limited)
			_,_=router_worker.RouteNotification()
		}
	}
	if debug {
		l4g.Debug(" Multicast request has been accepted and is processing ...:")
	}
	return
}
