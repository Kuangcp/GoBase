# kserver
1. 静态文件服务器
1. 文件上传服务器

go build; docker build -f kserver.dockerfile -t ksever:1.0.5 .

CGO_ENABLED=0 go build; docker build -f kserver-alpine.dockerfile -t ksever:1.0.5-alpine .

## change log
- 1.0.4 add upload function
- 1.0.5 index.html