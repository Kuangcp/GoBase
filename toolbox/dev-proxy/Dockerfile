FROM ubuntu:22.04

ENV TZ=Asia/Shanghai
EXPOSE 1234
EXPOSE 1235
RUN mkdir -p /app
ADD dev-proxy.bin /app

ENTRYPOINT ["/app/dev-proxy.bin"]