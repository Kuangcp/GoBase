FROM ubuntu:24.04

ENV TZ=Asia/Shanghai
EXPOSE 1234
EXPOSE 1235

RUN apt update -y && apt install -y tzdata curl && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt autoclean && rm -rf /var/lib/apt/lists/

RUN mkdir -p /app

ADD dev-proxy.bin /app

ENTRYPOINT ["/app/dev-proxy.bin"]
