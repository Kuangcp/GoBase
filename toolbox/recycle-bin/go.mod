module github.com/kuangcp/gobase/toolbox/recycle-bin

go 1.14

require (
	github.com/kuangcp/gobase/pkg/cuibase v1.0.0
	github.com/kuangcp/logger v1.0.8
	github.com/manifoldco/promptui v0.8.0
)

replace github.com/kuangcp/gobase/pkg/cuibase => ./../../pkg/cuibase
