#!/usr/bin/env bash
rm /tmp/.X0-lock

cp /index.html /usr/share/novnc/
mkdir -p /root/.vnc
x11vnc -storepasswd ${VNC_PASSWORD:-vncpass} /root/.vnc/passwd

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