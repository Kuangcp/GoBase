module keylogger

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/gotk3/gotk3 v0.5.2
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/gobase/pkg/ginhelper v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/logger v1.0.3
	github.com/manifoldco/promptui v0.8.0
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/rakyll/statik v0.1.7
)

replace github.com/kuangcp/gobase/pkg/cuibase => ./../../pkg/cuibase
