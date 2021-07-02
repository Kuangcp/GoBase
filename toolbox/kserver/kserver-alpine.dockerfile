FROM alpine:3.10

COPY kserver /bin

WORKDIR /data

CMD ["kserver"]