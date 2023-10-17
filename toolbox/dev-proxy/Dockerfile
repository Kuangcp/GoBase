FROM ubuntu

ENV TZ=Asia/Shanghai
RUN mkdir -p /app
ADD dev-proxy.bin /app
EXPOSE 1234
EXPOSE 1235
ENTRYPOINT /app/dev-proxy.bin

