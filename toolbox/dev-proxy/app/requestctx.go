package app

import "github.com/kuangcp/gobase/toolbox/dev-proxy/core"

type ReqCtx struct {
	reqLog        *core.ReqLog[core.Message]
	proxyLog      string
	proxyType     int
	startMs       int64
	ignoreStorage bool
}
