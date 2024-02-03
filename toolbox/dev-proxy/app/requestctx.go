package app

import "github.com/kuangcp/gobase/toolbox/dev-proxy/core"

type ReqCtx struct {
	reqLog      *core.ReqLog[core.Message]
	proxyLog    string
	proxyType   string
	startReq    int64
	receiveReq  int64
	needStorage bool
}
