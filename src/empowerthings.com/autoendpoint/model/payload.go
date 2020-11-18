package model

import(
	l4g "code.google.com/p/log4go"
	"net/http"
	"time"
)

// Payload_B used for a direct Webpush notification.
type Payload_A struct { 

	ChannelID          string            `json:"channelID"`
	
	Version            string            `json:"version"`

	TTL                int            `json:"ttl"` 
	
	Topic              interface{}       `json:"topic"`
	
	Timestamp          int64             `json:"timestamp"`
	
	Data               string            `json:"data"`
	
	Headers            Headers           `json:"headers"`
}

// Payload_B used for a direct Webpush empty notification.
type Payload_B struct { 


	ChannelID          string            `json:"channelID"`
	
	Version            string            `json:"version"`

	TTL                int            `json:"ttl"` 
	
	Topic              interface{}       `json:"topic"`
	
	Timestamp          int64             `json:"timestamp"`

}
func NewPayload_A(chid string, message_id []byte, headers Headers, msg string, request_headers http.Header, ttl int, topic string, uaid string, debug bool) (dataA *Payload_A) {
	dataA = &Payload_A{
		ChannelID: chid,
		Version:   string(message_id),
		TTL:       ttl,
		Topic:     []byte(topic),
		Timestamp: time.Now().Unix(),
		Data:      msg,
		Headers:   headers,
	}
	if debug {
		l4g.Info("Message to be sent: %s  to client with uaid:%s, and chid:%s", dataA.Data, uaid,chid)
	}	
	return
}

func NewPayload_B(chid string, message_id []byte, request_headers http.Header, topic string, ttl int, debug bool) (dataB *Payload_B) {
	dataB = &Payload_B{
		ChannelID: chid,
		Version:   string(message_id),
		TTL:       ttl,
		Topic:     []byte(topic),
		Timestamp: time.Now().Unix(),
	}
	if debug {
		l4g.Info("Chid: %s", dataB.ChannelID)
	}
	return
}
