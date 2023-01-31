package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
)

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	http.HandleFunc("/list", pageListReqHistory)
	http.HandleFunc("/flushAll", flushAllData)

	http.ListenAndServe(fmt.Sprintf(":%v", queryPort), nil)
}

func flushAllData(_ http.ResponseWriter, _ *http.Request) {
	result, err := connection.ZRange(TotalReq, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, key := range result {
		db.Delete([]byte(key), nil)
	}

	connection.Del(TotalReq)
	logger.Info("delete: ", len(result))
}

func pageListReqHistory(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	page := values.Get("page")
	size := values.Get("size")
	pageResult := pageQueryReqLog(page, size)
	result := ResultVO[*PageVO[*ReqLog[MessageVO]]]{}
	if pageResult == nil {
		result.Code = 101
		result.Msg = "invalid data"
	} else {
		result.Code = 0
		result.Data = pageResult

		hiddenHeaderEachLog(pageResult)
	}

	writer.Header().Set("Content-Type", "application/json")
	buffer := toJSONBuffer(result)
	writer.Write(buffer.Bytes())
}

func hiddenHeaderEachLog(pageResult *PageVO[*ReqLog[MessageVO]]) {
	if pageResult.Data == nil {
		return
	}
	for _, v := range pageResult.Data {
		hiddenHeader(v.Request.Header)
		hiddenHeader(v.Response.Header)
	}
}

func hiddenHeader(header http.Header) {
	delete(header, "User-Agent")
	delete(header, "Accept-Encoding")
	delete(header, "Referer")
	delete(header, "Cache-Control")
	delete(header, "Accept-Language")
	delete(header, "Pragma")
}
