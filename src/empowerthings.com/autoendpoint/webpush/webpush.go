package webpush
import (
	"net/http"
	l4g "code.google.com/p/log4go"
	"github.com/gocql/gocql"
	"github.com/statsd"
	"empowerthings.com/autoendpoint/routercore"
)
func Route_Webpush_Message(encrypted_endpoint string, body []byte, header http.Header, crypto_key string, debug bool, db_session *gocql.Session, statistics *statsd.Client, is_limited bool) (response_code int, response_body string) {
	if debug {
		l4g.Debug("Routing WebPush notification.")
	}
	router_worker := routercore.NewRouterWorker(encrypted_endpoint, body, header, crypto_key, debug, db_session, 2, false, statistics, false, is_limited)
	response_code, response_body= router_worker.RouteNotification()
	return
}
