#!/bin/bash

/data/autopush/configs/generate_config.sh

sleep 1

# Start syslog
rsyslogd -c 5

./libexec/start_rep.sh
sleep 2
cat /data/autopush/.version >> /data/autopush/log/autoendpoint.log
tail -f /data/autopush/log/autoendpoint.log

