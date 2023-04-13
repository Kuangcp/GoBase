# Hosts Group
> inspire from [SwitchHosts](https://oldj.github.io/SwitchHosts/)

Just webserver, provide curd api to management hosts file.

hosts file location:
1. Windows `C:\Windows\System32\drivers\etc\hosts`
1. Linux `/etc/hosts`

Support nginx
1. `sudo ./hosts-group -mode nginx -f /etc/nginx/conf.d/static.conf`

## Version
1.4.0
1. simple support nginx config

## TODO
1. independent window

*******************
## Install
1. https://codemirror.net/5/ 或者  https://codemirror.net/5/codemirror.zip 下载后，复制对应文件到项目目录中

*******************
> Linux

`go install`

*******************

> Windows

1. `make buildExe`
    - only support build via linux ...
1. run hosts-group.xxx.exe on Windows `Administer`

[blog: exe add icon](https://blog.csdn.net/u014633966/article/details/82984037)

## Dev Tips
1. `hosts` go run . -d -D
1. `nginx` go build && sudo ./hosts-group -d -D -mode nginx -f /etc/nginx/conf.d/static.conf

TODO 
1. js 文件缺失 
1. Linux编译打包流程未抽出来


