module github.com/kuangcp/gobase/toolbox/dev-proxy

go 1.19

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/kuangcp/gobase/pkg/ctool v1.0.8
	github.com/kuangcp/logger v1.0.9
	github.com/syndtr/goleveldb v1.0.0
	github.com/tidwall/pretty v1.2.1
)

require (
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.24.2 // indirect
)

replace github.com/kuangcp/gobase/pkg/ctool => ../../pkg/ctool
