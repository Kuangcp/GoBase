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
		Host string `form:"host"`
		Key  string `form:"key"`
		Val  string `form:"val"`
	}
	var d data
	json.Unmarshal(bd, &d)
	core.SetHeader(d.Host, d.Key, d.Val)
	return ctool.Success[string]()
}

func GetReqHeader(request *http.Request) ctool.ResultVO[map[string]string] {
	headers := core.GetHeaders(request.URL.Query().Get("host"))
	return ctool.SuccessWith(headers)
}

func DelReqHeader(request *http.Request) ctool.ResultVO[string] {
	query := request.URL.Query()
	host := query.Get("host")
	key := query.Get("key")
	core.DeleteHeader(host, key)
	return ctool.Success[string]()
}
