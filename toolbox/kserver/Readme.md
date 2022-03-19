# kserver
1. 静态文件服务器
1. 文件上传服务器

## build favicon.ico
magick convert -background none favicon.svg -define icon:auto-resize favicon.ico

## build docker image

> Ubuntu

go build; docker build -f kserver.dockerfile -t ksever:$(./kserver -h | grep Version | awk '{print $2}') .

> Alpine

CGO_ENABLED=0 go build; docker build -f kserver-alpine.dockerfile -t ksever:$(./kserver -h | grep Version | awk '{print $2}')-alpine .

### use

1. docker run --name ss -it -v `pwd`:/app -P ksever:1.0.8-alpine
2. dockerfile
    ```dockerfile
    FROM ksever:1.0.8-alpine

    COPY . /app

    WORKDIR /app

    CMD ["kserver"]
    ```

## change log
- 1.0.4 add upload function
- 1.0.5 index.html
- 1.0.6 img page