module drinkwater

go 1.21

require (
	github.com/go-echarts/go-echarts/v2 v2.3.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/kuangcp/gobase/pkg/ctool v1.1.9
	github.com/kuangcp/logger v1.0.9
)

require (
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.28.1 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848 // indirect
)

replace github.com/kuangcp/gobase/pkg/ctool v1.1.9 => ../../pkg/ctool
