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
> Linux

`make install`

*******************

> Windows

1. `make buildExe`
    - only support build via linux ...
1. run hosts-group.xxx.exe on Windows `Administer`

[blog: exe add icon](https://blog.csdn.net/u014633966/article/details/82984037)

## Dev Tips
1. `hosts` go run . -d -D
1. `nginx` go build && sudo ./hosts-group -d -D -mode nginx -f /etc/nginx/conf.d/static.conf
