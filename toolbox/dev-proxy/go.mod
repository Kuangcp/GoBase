module dev-proxy

go 1.19

require (
	github.com/getlantern/systray v1.1.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/kuangcp/gobase/pkg/ctool v1.0.8
	github.com/kuangcp/logger v1.0.9
	github.com/syndtr/goleveldb v1.0.0
)

require (
	github.com/getlantern/context v0.0.0-20190109183933-c447772a6520 // indirect
	github.com/getlantern/errors v0.0.0-20190325191628-abdb3e3e36f7 // indirect
	github.com/getlantern/golog v0.0.0-20190830074920-4ef2e798c2d7 // indirect
	github.com/getlantern/hex v0.0.0-20190417191902-c6586a6fe0b7 // indirect
	github.com/getlantern/hidden v0.0.0-20190325191715-f02dbb02be55 // indirect
	github.com/getlantern/ops v0.0.0-20190325191751-d70cb0d6f85f // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.24.2 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	golang.org/x/sys v0.3.0 // indirect
)

replace github.com/kuangcp/gobase/pkg/ctool => ../../pkg/ctool
