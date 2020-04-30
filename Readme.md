# Go基础学习

[![codebeat badge](https://codebeat.co/badges/7d223b91-e7e3-4241-a404-8463e1f16fce)](https://codebeat.co/projects/github-com-kuangcp-gobase-master)

1. go get github.com/kuangcp/gobase/tool-box/kserver
    - 静态文件 web服务器
1. go get github.com/kuangcp/gobase/tool-box/count
    - 统计汉字工具
1. go get github.com/kuangcp/gobase/tool-box/pretty-json
    - json 高亮 格式化
1. go get github.com/kuangcp/gobase/mybook
    - 记账工具
1. go get github.com/kuangcp/gobase/tool-box/keyboard-man
	- 按键监听 统计

找出可执行文件  find . -type f -exec file {} + | grep "ELF*" | sed 's/^.\///g' | awk '{print $1}' | sed 's/://g'

