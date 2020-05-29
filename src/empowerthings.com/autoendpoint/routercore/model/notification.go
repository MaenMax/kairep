package model 


	
import (
	
//	"net/http"
	"github.com/gocql/gocql"
	
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
















                           



