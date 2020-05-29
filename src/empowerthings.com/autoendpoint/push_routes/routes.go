package push_routes

import (
	"empowerthings.com/autoendpoint/push_handlers"
	"empowerthings.com/autoendpoint/routercore/db/cassandra"
	"empowerthings.com/autoendpoint/routercore/db/redis/counterdb"
)

func init_router_ep() {



	push_handlers.SetCryptoKey(key)
	push_handlers.SetDebug(debug)
	push_handlers.SetSession(session)
	push_handlers.SetPayloadSize(payload)
	push_handlers.SetStats(stat)
	if maximum_msg != 0 {
		push_handlers.EnableLimitation()
	}
	cassandra.SetKeyspace(keyspace) 
	cassandra.SetRouterTable(r_table_name)
	counterdb.SetMaxMsg(maximum_msg)



	
	routes = append(routes, Route{
		"Routing Webpush Message",
		"POST",
		"/wpush/{api_ver}/{encrypted_endpoint}",
		push_handlers.WebPushHandler,
	})

	routes = append(routes, Route{
		"Routing SimplePush Message",
		"PUT",
		"/spush/{api_ver}/{encrypted_endpoint}",
		push_handlers.SimplePushHandler,
	})

	routes = append(routes, Route{
		"Routing SimplePush Message.ONLY for GoFlip2",
		"PUT",
		"/update/{encrypted_endpoint}",
		push_handlers.GoFlip2SPushHandler,
	})

	routes = append(routes, Route{
		"Get the health status. To be used by haproxy",
		"GET",
		"/check_health/",
		push_handlers.Get_Health_Status,
	})

	routes = append(routes, Route{
		"Multicast API",
		"POST",
		"/bwpush/",
		push_handlers.Multicast,
	})


}
