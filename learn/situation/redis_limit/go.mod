module limit

go 1.19

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/kuangcp/gobase/pkg/sizedpool v1.0.2
	github.com/kuangcp/logger v1.0.9
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.6 // indirect
)

replace (
	github.com/kuangcp/gobase/pkg/sizedpool v1.0.2 => ../../../pkg/sizedpool
)