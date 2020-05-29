#!/bin/sh

# Finding out the current script directory and name out of how we have been called
SCRIPT_DIR=`dirname $0`
SCRIPT_NAME=`basename $0`
RUN_DIR=${SCRIPT_DIR}/../run

if [ -f ${RUN_DIR}/autoendpoint.pid ]; then
    pid=`cat ${RUN_DIR}/autoendpoint.pid`
    if [ "x$pid" != "x" ]; then
        kill >/dev/null 2>&1 $pid
    fi
    rm -f ${RUN_DIR}/autoendpoint.pid
else
    echo "autoendpoint is not running!"
fi