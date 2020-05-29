package config

/**
  Configuration object of the REP node.
*/

type RedisConfig struct {
	// Node addresses of Redis DB
	Host string

	Table_Prefix string

	Max_Redirects int

	Connection_Pool_Size int 
}

type AutoEndpointConfig struct {

	/**
	  Hostname of REP node. Should match with the host name set in CEP node.
	  Example: If CEP gives this endpoint: https://push.test.kaiostech.com:8443/v1/wpush/gAAAA... then hostname here  must be push.test.kaiostech.com and port must be 8443

	  Default: localhost

	*/

	Endpoint_Hostname string

	/**
	  Port number on  which to run autoendpoint(REP) node.


	  Default: 8080

	*/
	Endpoint_Port int

	/**
	  Debug log level

	  Default: false

	*/

	Debug bool

	/**
	  Hostname of CassandraDB

	  Default: localhost
	*/

	Cass_Address string

	/**

	  Keyspace name of  CassandraDB

	       Default: autopush

	*/

	Keyspace string

	/**

	  Maximum Payload size of multicast API in Bytes.

	      Default: 1000

	*/

	Max_Payload int

	/**

	StatsD host IP address.

	*/

	Statsd_Host string

	Statsd_Port int

	/**

	  Redis server configuration.

	*/

	Redis RedisConfig

	Max_Msg int

	// Maximum number of characters in a notification 
	Max_Msg_Length int 

	//Name of router table to be used. 
	Router_Table_Name  string
}
