package web

import (
	"encoding/json"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"io"
	"net/http"
)

func SetReqHeader(request *http.Request) ctool.ResultVO[string] {
	bd, err := io.ReadAll(request.Body)
	if err != nil {
		return ctool.FailedWithMsg[string](err.Error())
	}

	type data struct {
		Key string `form:"key"`
		Val string `form:"val"`
	}
	var d data
	json.Unmarshal(bd, &d)
	core.SetHeader(d.Key, d.Val)
	return ctool.Success[string]()
}

func GetReqHeader(request *http.Request) ctool.ResultVO[map[string]string] {
	headers := core.GetHeaders()
	return ctool.SuccessWith(headers)
}

func DelReqHeader(request *http.Request) ctool.ResultVO[string] {
	key := request.URL.Query().Get("key")
	core.DeleteHeader(key)
	return ctool.Failed[string]()
}
