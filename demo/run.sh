#!/bin/bash

go build tcp_server.go

./tcp_server &

sleep 2

kill -USR2 $(pidof tcp_server)

echo 123 | nc 127.1 12345
echo 123 | nc 127.1 12345
echo 123 | nc 127.1 12345
echo 123 | nc 127.1 12345
echo 123 | nc 127.1 12345

sleep 2

kill $(pidof tcp_server)
