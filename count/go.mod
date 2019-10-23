module github.com/kuangcp/gobase/count

go 1.13

require (
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/kuangcp/gobase/cuibase v0.0.0-20191021063012-d92aa1f00d16
)

replace github.com/kuangcp/gobase/cuibase => ./../cuibase
