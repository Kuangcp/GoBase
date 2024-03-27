package app

import "github.com/kuangcp/gobase/toolbox/dev-proxy/core"

type ReqCtx struct {
	reqLog       *core.ReqLog[core.Message]
	proxyLog     string
	proxyType    string
	startReqMs   int64
	receiveReqMs int64
	needStorage  bool
	originUrl    string
}
