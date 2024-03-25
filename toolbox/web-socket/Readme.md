## 客户端 ws-client

> 客户端打包目录下运行， 启动客户端容器建立 63976 个连接， 理论上可以建立 65535 端口范围设置为 1 65535

docker run -it --sysctl net.ipv4.ip_local_port_range="1024 65000" --name test-client --rm -v `pwd`:/data/ ubuntu:21.10

> 容器内运行

go run . -LH 192.168.7.54:7094 -LP /ws -d 7 -n 60000

## 服务端 ws-server

ulimit -n 800000 && go run . -s

