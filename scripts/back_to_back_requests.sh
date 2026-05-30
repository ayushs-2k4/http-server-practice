#!/usr/bin/env bash

exec 3<>/dev/tcp/127.0.0.1/8080

printf 'GET /first HTTP/1.1\r\n' >&3
printf 'Host: localhost\r\n' >&3
printf '\r\n' >&3

printf 'GET /second HTTP/1.1\r\n' >&3
printf 'Host: localhost\r\n' >&3
printf '\r\n' >&3

cat <&3