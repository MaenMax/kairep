package push_routes

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/statsd"
)

// Route describes an API route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type T_Routes []Route

var (
	routes        T_Routes
	key           string
	session       *gocql.Session
	debug         bool
	keyspace      string
	payload       int
	stat          *statsd.Client
	maximum_msg   int
	r_table_name  string
)

func NewRouter(encryption_key string, debug_log_level bool, cass_session *gocql.Session, cass_keyspace string, max_payload int, stats *statsd.Client, max_msg int, router_table_name string) *mux.Router {

	//init_router_ep(encryption_key, debug_log_level)
	key = encryption_key
	debug = debug_log_level
	session = cass_session
	keyspace = cass_keyspace
	stat = stats
	payload = max_payload
	maximum_msg = max_msg
	r_table_name = router_table_name

	init_router_ep()

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}
