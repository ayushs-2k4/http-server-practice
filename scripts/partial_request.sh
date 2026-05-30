#!/usr/bin/env bash

exec 3<>/dev/tcp/127.0.0.1/8080

BODY='{"name":"ayush","msg":"hello"}'
CONTENT_LENGTH=${#BODY}

printf 'POST /hello HT' >&3
sleep 1

printf 'TP/1.1\r\n' >&3
sleep 1

printf 'Host: localhost\r\n' >&3
sleep 1

printf 'User-Agent: sl' >&3
sleep 1

printf 'ow-client\r\n' >&3
sleep 1

printf 'Content-Type: application/json\r\n' >&3
sleep 1

printf "Content-Length: ${CONTENT_LENGTH}\r\n" >&3
sleep 1

printf '\r\n' >&3
sleep 1

printf '{"name":"' >&3
sleep 1

printf 'ayush",' >&3
sleep 1

printf '"msg":"' >&3
sleep 1

printf 'hello"}' >&3

cat <&3