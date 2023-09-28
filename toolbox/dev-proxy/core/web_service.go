package core

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
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
		result.Msg = "invalid param"
		return result
	}
	var pageResult *PageVO[*ReqLog[MessageVO]]
	if param.kwd != "" || param.date != "" {
		list, total := pageQueryReqLogByUrlKwd(param)
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

// search url
func pageQueryReqLogByUrlKwd(param *PageQueryParam) ([]*ReqLog[MessageVO], int) {
	var cursor uint64 = 0
	const fetchSize int64 = 100
	const maxPage = 6

	keys, cursor, err := Conn.HScan(RequestUrlList, cursor, "", fetchSize).Result()
	if err != nil {
		logger.Error(err)
		return nil, 0
	}

	startIdx := (param.page - 1) * param.size
	maxIdx := (param.page + maxPage) * param.size
	total := 0
	var list []*ReqLog[MessageVO]
	for len(keys) > 0 {
		for i := 0; i < len(keys); i += 2 {
			key := keys[i]
			val := keys[i+1]

			if !strings.Contains(url.QueryEscape(val), url.QueryEscape(param.kwd)) {
				continue
			}

			total++
			if total > startIdx && len(list) < param.size {
				log := getDetailByKey(key)
				if log == nil {
					Conn.HDel(RequestUrlList, key)
				}
				list = append(list, convertLog(log))
			}
			if total >= maxIdx {
				break
			}
		}
		// logger.Info("new loop", cursor)
		keys, cursor, err = Conn.HScan(RequestUrlList, cursor, "", fetchSize).Result()
		if err != nil {
			logger.Error(err)
			break
		}
		// finish full scan
		if cursor == 0 {
			break
		}
	}
	return list, total
}

// search url and header body
func pageQueryReqLogByKwd(param *PageQueryParam) ([]*ReqLog[MessageVO], int) {
	result, err := Conn.ZRevRange(RequestList, 0, -1).Result()
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
		keyList, err := Conn.ZRevRange(RequestList, start, end).Result()
		if err != nil {
			logger.Error(err)
			return nil
		}
		detail = queryLogDetail(keyList)
	}

	logs := stream.Just(detail...).Map(func(item any) any {
		return convertLog(item.(*ReqLog[Message]))
	})

	pageResult.Data = stream.ToList[*ReqLog[MessageVO]](logs)

	i, err := Conn.ZCard(RequestList).Result()
	if err == nil {
		pageResult.Total = int(i)
		pageResult.Page = int(i) / param.size
		if pageResult.Page*param.size < pageResult.Total {
			pageResult.Page += 1
		}
	}

	return &pageResult
}

// buildCommandById 目前仅支持常见的后端服务 GET POST请求
// TODO 待扩展,或者寻找成熟的方案
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
