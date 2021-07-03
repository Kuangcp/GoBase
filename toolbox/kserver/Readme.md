# kserver
1. 静态文件服务器
1. 文件上传服务器

## build docker image
1. go build; docker build -f kserver.dockerfile -t ksever:1.0.7 .

2. CGO_ENABLED=0 go build; docker build -f kserver-alpine.dockerfile -t ksever:1.0.7-alpine .

### use

1. docker run --name ss -it -v `pwd`:/data -P ksever:1.0.7-alpine
2. dockerfile
    ```dockerfile
    FROM ksever:1.0.7-alpine

    COPY . /data

    WORKDIR /data

    CMD ["kserver","-s"]
    ```

## change log
- 1.0.4 add upload function
- 1.0.5 index.html