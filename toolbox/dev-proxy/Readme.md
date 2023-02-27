# dev-proxy
后端视角 代理

https://blog.csdn.net/FlayHigherGT/article/details/109243249  
https://blog.csdn.net/FlayHigherGT/article/details/109243739  

- [ ] websocket 代理

## Config 
https://highlightjs.org/download/

~/.dev-proxy/dev-proxy.json

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

routers 每两个组成一对，前后元素构成源头和目标路径的映射。

- http://host1:port1/api, http://host2:port2
  - /api/a -> /a 
  - /api/b/c -> /b/c

- http://host1:port1/api, http://host2:port2/api2
  - /api/a -> /api2/a
  - host1:port1/api2/a -> host1:port1/api2/a










