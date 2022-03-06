# bookkeeping
> 2020.1.19 

1. 前端：Vue Element-UI Echarts
1. 后端：Gin Gorm SQLite

- [gorm](gorm.io/zh_CN/)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)

[日志级别](https://github.com/wonderivan/logger) : TRAC DEBG INFO

## 配置使用

| 命令 | 作用 |
|:----|:----|
| make install |`初始化配置和数据库` 
| make build   |`编译前后端 并 启动应用`
| make run     |`启动应用`

## IDE开发

IDE内运行时默认使用当前目录旁data目录的配置文件

包命名： 模块
文件命名： 模块缩写_职责

## TODO
> 业务逻辑

1. 退款，报销 不计为 收入
1. 代购报销，不计为 开支
1. 设计 应收应付 对应人员体系
    - 新建条目 直接新建用户名和id（好友表），后续下拉选择。
    - 每条应收应付独立， 应收应付表 关联人id，记录id。 聚合可以得到人情况

> 技术需求
1. 加密数据库 成本比较大
1. 数据库文件云同步

