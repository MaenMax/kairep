#!/bin/bash

#endpoint_hostname = "maen"
#endpoint_port = 8082
#debug = false
#cass_address = "172.31.78.130,172.31.1.243,172.31.6.210,172.31.69.146,172.31.24.152,172.31.27.215"
#keyspace = "autopush"
#statsd_host = "3.208.185.178"
#statsd_port = 8125
#max_msg_length = 4096

# max_msg: The maximum number of messages per UAID. Set to Zero in order to disable message limitation.
# Disabling message limitation(setting it to Zero) means that the is no need to connect to RedisDB. So connection will not be established to Redis.
#max_msg = 5

#[Redis]
#	Host = "172.31.54.10:7000,172.31.54.10:7000,172.31.54.10:7000,172.31.54.10:7000,172.31.54.10:7000,172.31.54.10:7000"
#	Table_Prefix = "counterdb"
#	Max_Redirects = 16
	
	#Maximum number of socket connections.
	#PoolSize applies per cluster node and not for the whole cluster
#	Connection_Pool_Size = 10



rm /data/autopush/configs/autopush_endpoint.conf
target_file=/data/autopush/configs/autopush_endpoint.conf


if [ "x" != "x${ENDPOINT_HOSTNAME}" ]; then
        echo >>${target_file} "endpoint_hostname = \"${ENDPOINT_HOSTNAME}\""
else 
        echo >>${target_file} "endpoint_hostname = \"localhost\""
fi

echo endpoint_port = 8082 >> ${target_file}
echo debug = false >> ${target_file}

if [ "x" != "x${ROUTER_TABLE}" ]; then
        echo >>${target_file} "router_table_name=\"${ROUTER_TABLE}\""
else
        echo >>${target_file} "router_table_name=router"
fi

if [ "x" != "x${CASS_ADDRESS}" ]; then
        echo >>${target_file} "cass_address = \"${CASS_ADDRESS}\""
else
        echo >>${target_file} "cass_address = \"localhost\""
fi

echo keyspace = \"autopush\" >> ${target_file}

if [ "x" != "x${STATSD_HOST}" ]; then
        echo >>${target_file} "statsd_host = \"${STATSD_HOST}\""
else
        echo >>${target_file} "statsd_host = \"localhost\""
fi
echo statsd_port = 8125 >> ${target_file}
echo max_msg_length = 4096 >> ${target_file}
echo max_msg = 5 >> ${target_file}
echo max_payload = 1000000 >> ${target_file}

if [ "x" != "x${MESSAGE_TABLE}" ]; then
        echo >>${target_file} "message_tablename= \"${MESSAGE_TABLE}\"" 
fi

echo [Redis] >> ${target_file}

if [ "x" != "x${REDIS_ADDRESS}" ]; then
        echo -e >>${target_file} ' \t'"Host = \"${REDIS_ADDRESS}\""
else
        echo -e >>${target_file} ' \t'"Host = \"localhost\""
fi
echo -e ' \t 'Table_Prefix = \"counterdb\" >> ${target_file}
echo -e ' \t 'Max_Redirects = 16 >> ${target_file}
echo -e ' \t' Connection_Pool_Size = 10 >>${target_file}

if [ "x" != "x${EFK_ADDRESS}" ]; then
       echo >>/etc/rsyslog.conf local3.*  @@${EFK_ADDRESS}:514
else
       echo >>/etc/rsyslog.conf local3.*  @@127.0.0.1:514
fi

