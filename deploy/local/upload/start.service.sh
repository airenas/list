#!/bin/sh
# setting rabbitmq pass from secret
pass=${MESSAGESERVER_PASS}
if ! [ -z "${MESSAGESERVER_PASS_FILE}" ]; then
    echo "Setting pass from file ${MESSAGESERVER_PASS_FILE}" 
    if ! [ -f "${MESSAGESERVER_PASS_FILE}" ]; then
        echo "ERROR: Environmnt MESSAGESERVER_PASS_FILE is set but no file ${MESSAGESERVER_PASS_FILE}" 
        echo "Exit service startup script!" 
        exit 1
    fi
    pass=$(cat ${MESSAGESERVER_PASS_FILE}) 
fi
# starting
echo "Starting service ..."
export MESSAGESERVER_PASS=$pass 
./uploadService 

