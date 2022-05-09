module kwebdav

go 1.17

require (
	github.com/kuangcp/gobase/pkg/cuibase v1.0.6
	github.com/kuangcp/logger v1.0.8
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4
)
replace (
	github.com/kuangcp/gobase/pkg/cuibase v1.0.6 => ../../pkg/cuibase
)