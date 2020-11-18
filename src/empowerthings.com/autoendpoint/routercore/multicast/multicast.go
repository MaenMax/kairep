package multicast

import(
	
	"empowerthings.com/autoendpoint/model"
	"encoding/json"
	l4g "code.google.com/p/log4go"
//	"errors"
)
func Extract_List(body []byte, debug bool) (list []string , message []byte, err error) {	
	if debug {
		l4g.Debug("Extracting multicast list.")
	}
	var msg model.Multicast_Msg
	err = json.Unmarshal(body, &msg)
	return msg.Endpoints,msg.Message, err
}
