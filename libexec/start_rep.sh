#!/bin/sh

# Finding out the current script directory and name out of how we have been called
SCRIPT_DIR=`dirname $0`
SCRIPT_NAME=`basename $0`
BIN_DIR=${SCRIPT_DIR}/../bin
LOG_DIR=${SCRIPT_DIR}/../log
RUN_DIR=${SCRIPT_DIR}/../run

if [ ! -d "$LOG_DIR" ]; then
    mkdir -p "$LOG_DIR" || { echo "Failed to create 'log' directory. Aborting ..."; exit 1; }
fi

if [ ! -d "$RUN_DIR" ]; then
    mkdir -p "$RUN_DIR" || { echo "Failed to create 'run' directory. Aborting ..."; exit 1; }
fi

if [ "x${CRYPTO_KEY}" = "x" ]; then
    echo "No CRYPTO_KEY defined!"
    exit 2
fi

if [ "x${CASS_USERNAME}" = "x" ]; then
    echo "No CASS_USERNAME defined!"
    exit 3
fi

if [ "x${CASS_PASSWORD}" = "x" ]; then
    echo "No CASS_PASSWORD defined!"
    exit 4
fi

if [ "x${REDIS_PASSWORD}" = "x" ]; then
    echo "No REDIS_PASSWORD defined!"
    exit 4
fi

nohup >${LOG_DIR}/autoendpoint.log 2>&1 sh -c "CASS_USERNAME=${CASS_USERNAME} CASS_PASSWORD=${CASS_PASSWORD} REDIS_PASSWORD=${REDIS_PASSWORD} CRYPTO_KEY=${CRYPTO_KEY} ${BIN_DIR}/autoendpoint" &

pid=$!
echo "autoendpoint  OK ($pid)"
echo $pid >${RUN_DIR}/autoendpoint.pid

