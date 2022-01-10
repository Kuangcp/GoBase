module github.com/kuangcp/gobase/toolbox/hosts-group

go 1.16

require (
	github.com/getlantern/systray v1.1.0
	github.com/gin-gonic/gin v1.7.1
	github.com/kuangcp/gobase/pkg/cuibase v1.0.0
	github.com/kuangcp/gobase/pkg/ghelp v1.0.0
	github.com/kuangcp/logger v1.0.3
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/zserge/lorca v0.1.10
)

replace github.com/kuangcp/gobase/pkg/ghelp => ./../../pkg/ghelp
