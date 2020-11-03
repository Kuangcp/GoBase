module keylogger

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201024141043-c83625c8aebf
	github.com/kuangcp/gobase/pkg/ginhelper v0.0.0-20201024141043-c83625c8aebf
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/rakyll/statik v0.1.7
	github.com/wonderivan/logger v1.0.0
)

replace github.com/kuangcp/gobase/pkg/cuibase => ./../../pkg/cuibase
