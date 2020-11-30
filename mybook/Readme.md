# bookkeeping

1. 前端：Vue Element-UI Echarts
1. 后端：Gin Gorm SQLite

- [gorm](gorm.io/zh_CN/)
- [go-sqlite3](https://github.com/mattn/go-sqlite3)

[日志级别](https://github.com/wonderivan/logger) : TRAC DEBG INFO

## 配置使用

1. `初始化应用配置和数据库` make install
1. `编译前后端 启动应用`: make run

## IDE开发

IDE内运行时默认使用当前目录旁data目录的配置文件

在IDE内执行单元测试 需创建 ~/.config/mybook.yml 文件，其中 db.file 为绝对路径

