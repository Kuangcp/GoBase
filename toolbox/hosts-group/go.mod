module github.com/kuangcp/gobase/toolbox/hosts-group

go 1.15

require (
	github.com/getlantern/systray v1.1.0
	github.com/gin-gonic/gin v1.6.3
	github.com/kuangcp/gobase/pkg/cuibase v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/gobase/pkg/ghelp v0.0.0-20201103041857-ea5c95ff0199
	github.com/kuangcp/logger v1.0.3
	github.com/rakyll/statik v0.1.7
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/webview/webview v0.0.0-20210201104136-ce27be3bc811
)

replace github.com/kuangcp/gobase/pkg/ghelp => ./../../pkg/ghelp
