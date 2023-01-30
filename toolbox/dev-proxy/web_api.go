package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
)

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	http.HandleFunc("/list", pageListReqHistory)

	http.ListenAndServe(fmt.Sprintf(":%v", queryPort), nil)
}

func pageListReqHistory(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	page := values.Get("page")
	size := values.Get("size")
	pageResult := pageQueryReqLog(page, size)
	result := ResultVO[*PageVO[ReqLog]]{}
	if pageResult == nil {
		result.Code = 101
		result.Msg = "invalid data"
	} else {
		result.Code = 0
		result.Data = pageResult

		hiddenHeader(pageResult)
	}

	buffer := toJSONBuffer(result)
	writer.Write(buffer.Bytes())
}

func hiddenHeader(pageResult *PageVO[ReqLog]) {
	if pageResult.Data == nil {
		return
	}
	for _, v := range pageResult.Data {
		header := v.Header
		delete(header, "User-Agent")
		delete(header, "Accept-Encoding")
		delete(header, "Referer")
		delete(header, "Cache-Control")
		delete(header, "Accept-Language")
		delete(header, "Pragma")
	}
}
