#!/bin/bash

# Generate Configuration
/data/autopush/configs/generate_config.sh

# Delete old log files that have been rotated every 10 mins
watch -t -n 600 rm -fr /data/autopush/log/autoendpoint.log.* >> /data/autopush/log/autoendpoint.log &

sleep 1

# Start syslog
rsyslogd -c 5

# Start REP
./libexec/start_rep.sh

sleep 2
#cat /data/autopush/.version >> /data/autopush/log/autoendpoint.log

# Send output to console
tail -f /data/autopush/log/autoendpoint.log

