FROM ubuntu:22.04

ENV TZ=Asia/Shanghai
EXPOSE 1234
EXPOSE 1235

RUN apt update && apt install tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt autoclean && rm -rf /var/lib/apt/lists/

RUN mkdir -p /app

ADD dev-proxy.bin /app

ENTRYPOINT ["/app/dev-proxy.bin"]
