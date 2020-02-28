FROM alpine-cst
# https://github.com/Kuangcp/DockerfileList/blob/master/alpine/alpine-cst.dockerfile
# docker build -t alpine-cst . -f alpine-cst.dockerfile

WORKDIR /app

RUN addgroup -S app \
    && adduser -S -g app app \
    &&  chown -R app:app /app

COPY ./bin/mybook .

USER app

ENTRYPOINT ["/app/mybook"]