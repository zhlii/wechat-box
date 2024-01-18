# https://github.com/zhlii/wine-box
FROM registry.cn-hangzhou.aliyuncs.com/xduo/wine-box:1.0.0

COPY root/ /

RUN bash -c 'nohup /entrypoint.sh 2>&1 &' && sleep 5 && /payloads.sh \
    && rm /tmp/.X0-lock 

EXPOSE 8888
EXPOSE 8889

ENTRYPOINT ["/wx-entrypoint.sh"]