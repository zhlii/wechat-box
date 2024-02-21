# https://github.com/zhlii/wine-box
FROM registry.cn-hangzhou.aliyuncs.com/xduo/wine-box:1.0.1

COPY root/ /

RUN bash -c 'nohup /entrypoint.sh 2>&1 &' && sleep 5 && /payloads.sh

ENTRYPOINT ["/wx-entrypoint.sh"]