module keylogger-gui

go 1.18

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gotk3/gotk3 v0.6.2
	github.com/kuangcp/gobase/pkg/ctool v1.2.0
	github.com/kuangcp/gobase/toolbox/keylogger v0.0.0-20220417170000-8b486ea221d8
	github.com/kuangcp/logger v1.0.9
	github.com/onsi/gomega v1.31.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/tklauser/go-sysconf v0.3.13 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
)

replace github.com/kuangcp/gobase/toolbox/keylogger => ../
