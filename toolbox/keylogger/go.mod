module keylogger

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gotk3/gotk3 v0.6.1
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20210507085753-f027a794a2ae
	github.com/kuangcp/gobase/pkg/ghelp v0.0.0-20210507085753-f027a794a2ae
	github.com/kuangcp/gobase/pkg/stopwatch v0.0.0-20210713173722-7ca2e3616271
	github.com/kuangcp/logger v1.0.4
	github.com/manifoldco/promptui v0.8.0
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/webview/webview v0.0.0-20210330151455-f540d88dde4e
)

replace (
	github.com/kuangcp/gobase/pkg/cuibase => ./../../pkg/cuibase
	github.com/kuangcp/gobase/pkg/ghelp => ./../../pkg/ghelp
)
