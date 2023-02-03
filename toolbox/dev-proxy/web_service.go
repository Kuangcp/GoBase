package main

import (
	"github.com/kuangcp/logger"
	"net/http"
	"strings"
)

func pageListReqHistory(request *http.Request) ResultVO[*PageVO[*ReqLog[MessageVO]]] {
	result := ResultVO[*PageVO[*ReqLog[MessageVO]]]{}
	param := parseParam(request)
	if param == nil {
		result.Code = 101
		result.Msg = "invalid data"
		return result
	}
	var pageResult *PageVO[*ReqLog[MessageVO]]
	if param.kwd != "" || param.prefix != "" {
		list, total := pageQueryReqLogByKwd(param)
		pageResult = &PageVO[*ReqLog[MessageVO]]{}
		pageResult.Data = list
		pageResult.Total = total
		pageResult.Page = 1
	} else {
		pageResult = pageQueryReqLogByIndex(param)
	}

	if pageResult == nil {
		result.Code = 101
		result.Msg = "no data"
	} else {
		result.Code = 0
		result.Data = pageResult
		hiddenHeaderEachLog(pageResult)
	}
	return result
}

func pageQueryReqLogByKwd(param *PageQueryParam) ([]*ReqLog[MessageVO], int) {
	result, err := connection.ZRange(TotalReq, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return nil, 0
	}

	total := 0
	var list []*ReqLog[MessageVO]
	for _, key := range result {
		if !strings.HasPrefix(key, param.prefix) {
			continue
		}
		log := matchDetailByKeyAndKwd(key, param.kwd)
		if log != nil {
			total++
			list = append(list, convertLog(log))
			//if total < param.page {
			//	list = append(list, convertLog(log))
			//}
		}
		if len(list) == param.page {
			break
		}
	}
	return list, total
}

// page start with 1
func pageQueryReqLogByIndex(param *PageQueryParam) *PageVO[*ReqLog[MessageVO]] {
	start, end := param.buildStartEnd()
	keyList, err := connection.ZRange(TotalReq, start, end).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}

	pageResult := PageVO[*ReqLog[MessageVO]]{}
	detail := queryLogDetail(keyList)
	pageResult.Data = convertList(detail, convertLog, nil)

	i, err := connection.ZCard(TotalReq).Result()
	if err == nil {
		pageResult.Total = int(i)
		pageResult.Page = int(i) / param.size
		if pageResult.Page*param.size < pageResult.Total {
			pageResult.Page += 1
		}
	}

	return &pageResult
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
