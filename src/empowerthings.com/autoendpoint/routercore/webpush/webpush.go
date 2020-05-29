package webpush

import (
	"net/http"

	l4g "code.google.com/p/log4go"
	"empowerthings.com/autoendpoint/routercore"
	"github.com/gocql/gocql"
	"github.com/statsd"
)

func Route_Webpush_Message(encrypted_endpoint string, body []byte, header http.Header, crypto_key string, debug bool, db_session *gocql.Session, w http.ResponseWriter, statistics *statsd.Client, is_limited bool) {

	if debug {

		l4g.Debug("Routing WebPush notification.")

	}

	router_worker := routercore.NewRouterWorker(encrypted_endpoint, body, header, crypto_key, debug, db_session, w, 2, false, statistics, false, is_limited)
	router_worker.RouteNotification()

}
