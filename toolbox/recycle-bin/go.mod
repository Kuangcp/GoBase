module github.com/kuangcp/gobase/toolbox/recycle-bin

go 1.14

require (
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201105021415-0bdbbc0a38fd
	github.com/kuangcp/logger v1.0.3
	github.com/manifoldco/promptui v0.8.0
)

replace (
	//github.com/kuangcp/logger v1.0.2 => /home/kcp/Code/go/logger/
	github.com/kuangcp/gobase/pkg/cuibase => ./../../pkg/cuibase
)
