module github.com/kuangcp/gobase/toolbox/md-formatter

go 1.14

require (
	github.com/go-git/go-git/v5 v5.0.0
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20210409093747-f38d0e9c0695
	github.com/wonderivan/logger v1.0.0
)
replace (
	github.com/kuangcp/gobase/pkg/cuibase => ../../pkg/cuibase
)