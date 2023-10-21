# dev-proxy

后端视角 代理

## 1. 设计

![](./img/design.drawio.svg)

![](./img/ws.drawio.svg)

实现功能:

1. 转发请求
2. 抓包
3. 动态修改 请求和响应
4. 监控完整请求响应过程中指标

参考:

- https://blog.csdn.net/FlayHigherGT/article/details/109243249
- https://blog.csdn.net/FlayHigherGT/article/details/109243739

技术栈: Golang Redis LevelDB

## 2. Install
> 当前代理有两套实现：一个需要安装证书支持HTTPS解密和修改，另一个仅支持HTTPS密文转发，但是都支持HTTP代理和修改。
> 镜像内默认是HTTPS实现

### Docker 
docker run --init -d --name dev-proxy --network host mythkuang/dev-proxy:1.3
- 低版本Docker还会遇到[issues](https://github.com/docker-library/golang/issues/467) 需要run添加参数：--privileged=true

### Go
1. go install 
1. dev-proxy 

******************

## 3. Config

HTTPS 证书安装： https://github.com/ouqiang/goproxy

### 3.1. Config页面方式
按页面填写配置规则后 按 Ctrl S 保存或者点击保存

配置规则参考配置文件中配置的示例

### 3.2. 配置文件方式
> 配置文件目录路径： ~/.dev-proxy/dev-proxy.json

- groups: 抓包并修改请求： 按正则匹配，配置域名或路径都可以
- proxy： 抓包规则：正则匹配 配置域名或路径都可以
- direct： 在域名已被抓包的前提下，正则匹配配置不抓包的子路径，忽略不关心的请求

```json
  {
  "groups": [
    {
      "name": "Test",
      "proxy_type": 1,
      "routers": [
        {
          "proxy_type": 1,
          "src": "http://localhost:32009/api1/(.*)",
          "dst": "http://127.0.0.1:8081/api1/$1"
        },
        {
          "proxy_type": 1,
          "src": "http://localhost:32009/api1/table/add(.*)",
          "dst": "http://127.0.0.1:8081/api1/table/tryApply$1"
        }
      ]
    }
  ],
  "proxy": {
    "name": "抓包",
    "proxy_type": 1,
    "paths": [
      "http://192.168.1.2:3209/api/(.*)",
      "http://192.168.1.9:3210/v1/user/(.*)"
    ]
  }
}

```

routers 每两个组成一对, 例如: 

> src: http://host1:port1/api   
> dst: http://host2:port2  

- /api/a -> /a
- /api/b/c -> /b/c

> src: http://host1:port1/api  
> dst: http://host2:port2/api2  
- /api/a -> /api2/a
- host1:port1/api2/a -> host1:port1/api2/a

## 4. Docker 镜像构建

- go build -o dev-proxy.bin
- docker build -t dev-proxy:1.x .

## 5. TODO

1. [x] 按URL域名统计请求频率和时间分布
1. [ ] websocket 代理
1. 移除Redis依赖, 缓存层使用 文件+内存 存储
    - 优点: 减少组件依赖
    - 缺点: 数据的一致性保证, 数据完整性保证 





