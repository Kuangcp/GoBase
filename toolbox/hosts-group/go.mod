module github.com/kuangcp/gobase/toolbox/hosts-group

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/gobase/pkg/ghelp v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/logger v1.0.3
)


replace github.com/kuangcp/gobase/pkg/ghelp => ./../../pkg/ghelp
