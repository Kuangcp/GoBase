package core

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func PageListReqHistory(request *http.Request) ResultVO[*PageVO[*ReqLog[MessageVO]]] {
	result := ResultVO[*PageVO[*ReqLog[MessageVO]]]{}
	param := &PageQueryParam{}
	if err := ctool.Unpack(request, param); err != nil {
		logger.Error(err)
		result.Code = 101
		result.Msg = "invalid param"
		return result
	}

	if param.Size <= 0 {
		param.Size = 1
	}
	if param.Page <= 1 {
		param.Page = 1
	}
	if param.Kwd != "" {
		param.Kwd = strings.TrimSpace(param.Kwd)
	}

	var pageResult *PageVO[*ReqLog[MessageVO]]
	if param.Kwd != "" {
		result.Msg = "kwd"
		list, total := pageQueryReqByUrlKwd(param)
		pageResult = &PageVO[*ReqLog[MessageVO]]{}
		pageResult.Data = list
		pageResult.Total = total
		page := total / param.Size
		if page*param.Size < total {
			page++
		}
		pageResult.Page = page
	} else if param.Date != nil && param.Id == "" {
		result.Msg = "time"
		pageResult = pageQueryLogByTime(param)
	} else {
		result.Msg = "index"
		pageResult = pageQueryLogByIdOrIndex(param)
	}

	if pageResult == nil {
		result.Code = 101
		result.Msg += " no data"
	} else {
		result.Code = 0
		result.Data = pageResult
		hiddenHeaderEachLog(pageResult)
	}
	return result
}

// search url
func pageQueryReqByUrlKwd(param *PageQueryParam) ([]*ReqLog[MessageVO], int) {
	var cursor uint64 = 0
	const fetchSize int64 = 100
	const maxPage = 5

	keys, cursor, err := Conn.HScan(RequestUrlList, cursor, "", fetchSize).Result()
	if err != nil {
		logger.Error(err)
		return nil, 0
	}

	startIdx := (param.Page - 1) * param.Size
	maxIdx := (param.Page + maxPage) * param.Size
	total := 0
	var list []*ReqLog[MessageVO]
	overflow := false
	for len(keys) > 0 {
		for i := 0; i < len(keys); i += 2 {
			key := keys[i]
			val := keys[i+1]

			if !strings.Contains(url.QueryEscape(val), url.QueryEscape(param.Kwd)) {
				continue
			}

			total++
			if total > startIdx && len(list) < param.Size {
				log := getDetailByKey(key)
				if log == nil {
					RemoveReqUrlKey(key)
				} else {
					list = append(list, convertLog(log))
				}
			}
			if total >= maxIdx {
				logger.Warn("reach max scan row", total, param.Size, len(list))
				overflow = true
				break
			}
		}
		if overflow {
			break
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

func pageQueryLogByTime(param *PageQueryParam) *PageVO[*ReqLog[MessageVO]] {
	zr, err := Conn.ZRevRangeByScoreWithScores(RequestList, redis.ZRangeBy{
		Min: fmt.Sprint(param.Date.UnixNano()),
		Max: fmt.Sprint(param.Date.Add(time.Hour * 24).UnixNano()),
	}).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}

	start, end := param.buildStartEnd()
	var keyList []string
	for i, val := range zr {
		if i < int(start) {
			continue
		}
		if i > int(end) {
			break
		}
		keyList = append(keyList, val.Member.(string))
	}

	detail := queryLogDetail(keyList)
	return detailToPage(param, detail, len(zr))
}

// page start with 1
func pageQueryLogByIdOrIndex(param *PageQueryParam) *PageVO[*ReqLog[MessageVO]] {
	var detail []*ReqLog[Message]
	if param.Id != "" {
		val := getDetailByKey(param.Id)
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

	total, _ := Conn.ZCard(RequestList).Result()
	return detailToPage(param, detail, int(total))
}

func detailToPage(param *PageQueryParam, detail []*ReqLog[Message], total int) *PageVO[*ReqLog[MessageVO]] {
	logs := stream.Just(detail...).Map(func(item any) any {
		return convertLog(item.(*ReqLog[Message]))
	})
	pageResult := PageVO[*ReqLog[MessageVO]]{}
	pageResult.Data = stream.ToList[*ReqLog[MessageVO]](logs)

	pageResult.Total = total
	pageResult.Page = total / param.Size
	if pageResult.Page*param.Size < pageResult.Total {
		pageResult.Page += 1
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
		hiddenHeader(v.Request.Header, v.Response.Header)
	}
}

func hiddenHeader(header ...http.Header) {
	for _, head := range header {
		for _, h := range HiddenHeader {
			delete(head, h)
		}
	}
}
