package model 
import (
	//	"net/http"
	l4g "code.google.com/p/log4go"
	"github.com/gocql/gocql"
	"strconv"
	"time"
	"fmt"
	fernet "fernet-go"
	"golang.org/x/exp/utf8string"
	"encoding/base64"
)	

const(
	LEGACY       string        = ""
)	
type WebPushNotification struct {
	Uaid                     string
	TTL                      int
	Data                     string
	Headers                  map[string] string
	TimeStamp                string
	UpdateID                string
	Session                 *gocql.Session
	SortedKey               string
	Debug                   bool
}
type SimplePushNotification struct {
	Version      string           `json:"version"`
	Data         string           `json:"data"`
	ChannelID    string           `json:"channel_id"`  
}

// CreateWebPushNotif function takes a required parameter to construct a WebPushNotification object described in "model" package. WebPushNotification object will then be used as a parameter to StoreMessage()  method that will store offline messages (Incase CEP response is anything other than 200).
func CreateWebPushNotif(message_id string, crypto_headers Headers, uaid string,ttl int,msg string,session *gocql.Session,sorted_key string, debug bool) (Notification *WebPushNotification) {
	if debug {
		l4g.Info("Creating WebPushNotification object.")
	}
	var wpush_headers map[string]string
	wpush_headers = make(map[string]string)
	wpush_headers["crypto_key"] = crypto_headers.CryptoKey
	wpush_headers["encoding"] = crypto_headers.Encoding
	wpush_headers["encryption"] = crypto_headers.Encryption
	str_timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	if len(str_timestamp) > 10 {
		str_timestamp = str_timestamp[0:10]
	}
	Notification = &WebPushNotification{
		Uaid:      uaid,
		TTL:       ttl,
		Data:      msg,
		Headers:   wpush_headers,
		TimeStamp: str_timestamp,
		UpdateID:  message_id,
		Session:   session,
		SortedKey: sorted_key,
		Debug:     debug,
	}
	return
}

func SortKey(message_id []byte, uaid string, chid string, topic string,  sortkey_timestamp int, debug bool) (sorted_key string) {
	/* Return an appropriate sort_key for this notification

	   For new messages:

	       02:{sortkey_timestamp}:{chid}

	   For topic messages:

	       01:{chid}:{topic}

	   Old format for non-topic messages that is no longer returned:

	       {chid}:{message_id}
	*/
	if  topic != "" {
		sorted_key = fmt.Sprintf("01:{%s}:{%s}", topic, uaid)
	} else if LEGACY != "" {
		sorted_key = fmt.Sprintf("{%s}:{%s}", chid, message_id)
	} else {
		//Created as late as possible when storing a message
		sorted_key = fmt.Sprintf("02:%v:%s", sortkey_timestamp, chid)
	}
	if debug {
		l4g.Info("Sort key used is: %s", sorted_key)
	}
	return
}

func GenerateNotificationId(uaid string, chid string, crypto_key_utf8_encoded *utf8string.String,sortkey_timestamp int ,topic string, debug bool ) (message_id []byte, err error) {
	/*
	   Notification ID serves a complex purpose. It's returned as the Location header
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

	if debug {
		l4g.Info("Generating notification ID")
	}
	var msg_key string
	
	if topic != "" {
		msg_key = "01" + ":" + uaid + ":" + chid + ":" + topic
	} else if LEGACY != "" {
		msg_key = "m" + ":" + uaid + ":" + chid
	} else {
		msg_key = "02" + ":" + uaid + ":" + chid + ":" + strconv.Itoa(sortkey_timestamp)
	}
	key := fernet.MustDecodeKeys(crypto_key_utf8_encoded.String())
	message_id, err = fernet.EncryptAndSign([]byte(msg_key), key[0])
	return
}

func FormatNotificationData(data []byte, debug bool) (msg string) {
	if debug {
		l4g.Info("Formating received notification.")
	}
	msg = base64.RawURLEncoding.EncodeToString(data)
	return msg
}
