package simplepush

import(
	"empowerthings.com/autoendpoint/model"
	"time"
	l4g "code.google.com/p/log4go"
	"strconv"
)
func  CreateSimplePushNotif(data []byte, Chid string, debug bool) (SimplePushNotification *model.SimplePushNotification) {

	if debug {
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
