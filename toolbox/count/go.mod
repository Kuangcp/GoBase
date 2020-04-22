module github.com/kuangcp/gobase/count

go 1.13

require (
	github.com/go-redis/redis/v7 v7.2.0 // indirect
	github.com/kuangcp/gobase/cuibase v0.0.0-20191021063012-d92aa1f00d16
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
)

replace github.com/kuangcp/gobase/cuibase => ./../../cuibase
