# Go基础学习

[![codebeat badge](https://codebeat.co/badges/7d223b91-e7e3-4241-a404-8463e1f16fce)](https://codebeat.co/projects/github-com-kuangcp-gobase-master)  
[![GoBase](https://goreportcard.com/badge/github.com/kuangcp/gobase)](https://goreportcard.com/report/github.com/kuangcp/gobase)  

| 安装 | 简介 |
|:----|:----|
| go get github.com/kuangcp/gobase/tool-box/kserver |  静态文件 web服务器
| go get github.com/kuangcp/gobase/tool-box/count | 统计汉字工具
| go get github.com/kuangcp/gobase/tool-box/pretty-json | json 高亮 格式化
| go get github.com/kuangcp/gobase/mybook | 记账工具
| go get github.com/kuangcp/gobase/tool-box/keyboard-man | 按键监听 统计
| go get github.com/kuangcp/gobase/tool-box/baidu-translation | 命令行 百度翻译

************************

- Tips 

> 找出可执行文件 加入 .gitignore `find . -type f -exec file {} + | grep " ELF " | sed 's/^.\///g' | awk '{print $1}' | sed 's/://g'`

