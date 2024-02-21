#!/usr/bin/env bash
function wechat() {
    while :
    do
        wechat-start >/dev/stdout 2>&1
    done
}

/entrypoint.sh &
sleep 5
wechat &
wait