module github.com/kuangcp/gobase/keylogger

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/kuangcp/gobase/cuibase v1.0.1-cuibase
	github.com/kuangcp/threadpool v1.0.5
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/wonderivan/logger v1.0.0
)

replace github.com/kuangcp/gobase/cuibase => ./../../cuibase
