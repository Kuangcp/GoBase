# dev-proxy
后端视角 代理

https://blog.csdn.net/FlayHigherGT/article/details/109243249  
https://blog.csdn.net/FlayHigherGT/article/details/109243739  

- [ ] websocket 代理

## Config 

~/.dev-proxy.json

```json
  [
    {
      "name": "test",
      "enable": 1,
      "routers": ["http://192.168.1.2:12009/tg/(.*)", "http://127.0.0.1:8081/tg/test/$1"]
    }
  ]
```

routers 每两个组成一对，前后元素构成源头和目标路径的映射。

- http://host1:port1/api, http://host2:port2
  - /api/a -> /a 
  - /api/b/c -> /b/c

- http://host1:port1/api, http://host2:port2/api2
  - /api/a -> /api2/a
  - host1:port1/api2/a -> host1:port1/api2/a

