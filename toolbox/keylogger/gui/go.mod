module keylogger-gui

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gotk3/gotk3 v0.6.1
	github.com/kuangcp/gobase/pkg/ctool v1.0.9
	github.com/kuangcp/gobase/toolbox/keylogger v0.0.0-20220417170000-8b486ea221d8
	github.com/kuangcp/logger v1.0.8
	github.com/onsi/gomega v1.27.10 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
)

replace github.com/kuangcp/gobase/toolbox/keylogger v0.0.0-20220417170000-8b486ea221d8 => ../
