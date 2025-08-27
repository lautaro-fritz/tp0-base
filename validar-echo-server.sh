#!/bin/bash

SERVER="server"
PORT=12345
MESSAGE="echo"
NETWORK="tp0_testing_net"

RESPONSE=$(docker run --rm --network "$NETWORK" busybox:latest sh -c "echo ${MESSAGE}| nc ${SERVER} ${PORT}")

if [ \"$RESPONSE\" = \"$MESSAGE\" ]; then
    echo "action: test_echo_server | result: success"
    exit 0
else
    echo "action: test_echo_server | result: fail"
    exit 1
fi
