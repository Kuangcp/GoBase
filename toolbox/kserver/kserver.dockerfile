FROM ubuntu:21.10

COPY kserver /bin

WORKDIR /data

CMD ["kserver","-s"]