package simplepush

import(
	"empowerthings.com/autoendpoint/model"
	"time"
	l4g "code.google.com/p/log4go"
//	"strconv"
)
func  CreateSimplePushNotif(data []byte, Chid string, debug bool) (SimplePushNotification *model.SimplePushNotification) {

	if debug {
		l4g.Info("Creating SimplePushNotification object.")
	}
	
	now := time.Now()
	secs := now.Unix()

	SimplePushNotification = &model.SimplePushNotification{
		Version:   secs,
		Data:      "",
		ChannelID: Chid,
	}
	return
}
