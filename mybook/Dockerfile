FROM alpine:3.11

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -g app app \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && chown -R app:app /app

COPY ./bin/mybook .

USER app

ENTRYPOINT ["/app/mybook"]