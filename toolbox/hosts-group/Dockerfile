FROM debian:11.6-slim

ENV TZ=Asia/Shanghai
RUN mkdir -p /app
ADD hosts-group.bin /app
ENTRYPOINT /app/hosts-group.bin

