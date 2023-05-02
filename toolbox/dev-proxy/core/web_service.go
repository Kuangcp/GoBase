package core

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"sort"
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
	if param.kwd != "" || param.date != "" {
		list, total := pageQueryReqLogByKwd(param)
		pageResult = &PageVO[*ReqLog[MessageVO]]{}
		pageResult.Data = list
		pageResult.Total = total
		page := total / param.size
		if page*param.size < total {
			page++
		}
		pageResult.Page = page
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
	result, err := connection.ZRevRange(RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return nil, 0
	}

	startIdx := (param.page - 1) * param.size
	maxIdx := (param.page + 4) * param.size
	total := 0
	var list []*ReqLog[MessageVO]
	for _, key := range result {
		if !strings.HasPrefix(key, param.date) {
			continue
		}
		log := matchDetailByKeyAndKwd(convertToDbKey(key), param.kwd)
		if log == nil {
			continue
		}
		total++
		if total > startIdx && len(list) < param.size {
			list = append(list, convertLog(log))
		}
		if total >= maxIdx {
			break
		}
	}
	return list, total
}

// page start with 1
func pageQueryReqLogByIndex(param *PageQueryParam) *PageVO[*ReqLog[MessageVO]] {
	pageResult := PageVO[*ReqLog[MessageVO]]{}
	var detail []*ReqLog[Message]
	if param.id != "" {
		val := getDetailByKey(param.id)
		if val == nil {
			return nil
		}
		detail = append(detail, val)
	} else {
		start, end := param.buildStartEnd()
		keyList, err := connection.ZRevRange(RequestList, start, end).Result()
		if err != nil {
			logger.Error(err)
			return nil
		}
		detail = queryLogDetail(keyList)
	}

	pageResult.Data = convertList(detail, convertLog, nil)

	i, err := connection.ZCard(RequestList).Result()
	if err == nil {
		pageResult.Total = int(i)
		pageResult.Page = int(i) / param.size
		if pageResult.Page*param.size < pageResult.Total {
			pageResult.Page += 1
		}
	}

	return &pageResult
}

func buildCommandById(id, selfProxy string) string {
	detail := getDetailByKey(id)
	if detail == nil {
		return ""
	}
	cmd := "curl "
	if selfProxy == "Y" {
		cmd += fmt.Sprintf(" -x 127.0.0.1:%v ", Port)
	}
	parseUrl, _ := url.Parse(detail.Url)
	cmd += parseUrl.Scheme + "://" + parseUrl.Host + parseUrl.Path
	if parseUrl.RawQuery != "" {
		query := url.PathEscape(parseUrl.RawQuery)
		query = strings.ReplaceAll(query, "&", "\\&")
		cmd += "\\?" + query
	}
	var key []string
	for k := range detail.Request.Header {
		key = append(key, k)
	}
	sort.Strings(key)
	for _, k := range key {
		val := detail.Request.Header.Values(k)
		for _, v := range val {
			cmd += fmt.Sprintf(" -H '%s: %s'", k, v)
		}
	}

	if len(detail.Request.Body) > 0 {
		cmd += fmt.Sprintf(" --data-raw $'%s'", string(detail.Request.Body))
	}
	//logger.Info(cmd)

	return cmd
}

func hiddenHeaderEachLog(pageResult *PageVO[*ReqLog[MessageVO]]) {
	if pageResult.Data == nil {
		return
	}
	for _, v := range pageResult.Data {
		if v == nil {
			continue
		}
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
