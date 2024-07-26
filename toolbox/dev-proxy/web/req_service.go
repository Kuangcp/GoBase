package web

import (
	"encoding/json"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"io"
	"net/http"
	"strings"
)

type HostHeaderParam struct {
	Host string `form:"host"`
	Key  string `form:"key"`
	Val  string `form:"val"`
}

func SetReqHeader(request *http.Request) ctool.ResultVO[string] {
	bd, err := io.ReadAll(request.Body)
	if err != nil {
		return ctool.FailedWithMsg[string](err.Error())
	}

	var d HostHeaderParam
	json.Unmarshal(bd, &d)
	core.SetHeader(d.Host, d.Key, d.Val)
	return ctool.Success[string]()
}

func SetReqHeaders(request *http.Request) ctool.ResultVO[string] {
	bd, err := io.ReadAll(request.Body)
	if err != nil {
		return ctool.FailedWithMsg[string](err.Error())
	}

	var ds []HostHeaderParam
	json.Unmarshal(bd, &ds)
	for _, d := range ds {
		core.SetHeader(d.Host, d.Key, d.Val)
	}
	return ctool.Success[string]()
}

func GetReqHeader(request *http.Request) ctool.ResultVO[core.TMap] {
	host := request.URL.Query().Get("host")
	return ctool.SuccessWith(core.QueryHeaders(host))
}

func DelReqHeader(request *http.Request) ctool.ResultVO[string] {
	bd, err := io.ReadAll(request.Body)
	if err != nil {
		return ctool.FailedWithMsg[string](err.Error())
	}

	var d HostHeaderParam
	json.Unmarshal(bd, &d)

	core.DeleteHeader(d.Host, d.Key)
	return ctool.Success[string]()
}

func DelReqHeaderViaHost(request *http.Request) ctool.ResultVO[string] {
	host := request.URL.Query().Get("host")
	hs := strings.Split(host, ",")
	for _, h := range hs {
		core.DeleteHeaderViaHost(h)
	}
	return ctool.Success[string]()
}
