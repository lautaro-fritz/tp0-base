#!/bin/bash

SERVER_CONTAINER_NAME="server"
CLIENT_CONTAINER_NAME="client1"
SERVER_PORT=12345
TEST_MESSAGE="hello_echo"

# Ejecuta netcat desde client1 hacia server, usando la red interna
RESPONSE=$(docker exec "$CLIENT_CONTAINER_NAME" sh -c "echo $TEST_MESSAGE | nc $SERVER_CONTAINER_NAME $SERVER_PORT")

# Verifica si el servidor respondió correctamente
if [ \"$RESPONSE\" = \"$TEST_MESSAGE\" ]; then
    echo "action: test_echo_server | result: success"
    exit 0
else
    echo "action: test_echo_server | result: fail"
    exit 1
fi
