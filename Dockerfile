FROM golang:1.20.13-alpine3.19 as go-builder

WORKDIR /tmp/rest

COPY ./rest .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o rest.exe


# https://github.com/zhlii/wine-box
FROM registry.cn-hangzhou.aliyuncs.com/xduo/wine-box:1.0.0

COPY root/ /
COPY --from=go-builder /tmp/rest/rest.exe /hook
COPY --from=go-builder /tmp/rest/config.yaml /hook

RUN bash -c 'nohup /entrypoint.sh 2>&1 &' && sleep 5 && /payloads.sh \
    && rm /tmp/.X0-lock 

EXPOSE 80
EXPOSE 8888
EXPOSE 8889

ENTRYPOINT ["/wx-entrypoint.sh"]