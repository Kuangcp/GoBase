FROM ubuntu:21.10

COPY kserver /bin

WORKDIR /app

CMD ["kserver"]