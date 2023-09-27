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
	"strconv"
	"strings"
	"time"
)

type (
	PerfPageVo struct {
		Url  string
		Tct  int
		Tall int
		TAvg int
		TP30 int
		TP50 int
		TP90 int
		TP95 int
		TP99 int
		Qps  int // 查询时间段内最大Qps
	}
	DetailVo struct {
		Url   string `json:"url"`
		Start string `json:"start"`
		End   string `json:"end"`
	}
)

func UrlFrequencyApi(writer http.ResponseWriter, request *http.Request) {
	minS := request.URL.Query().Get("min")
	maxS := request.URL.Query().Get("max")
	min, _ := strconv.Atoi(minS)
	if min < 1 {
		min = 50
	}
	max, _ := strconv.Atoi(maxS)
	if max < 1 {
		max = 100
	}

	result, err := Conn.HGetAll(RequestUrlList).Result()
	if err != nil {
		logger.Error(err)
		writeJsonRsp(writer, err.Error())
		return
	}

	type Val struct {
		UrlMap  map[string]int
		HostMap map[string]int
	}

	allUrlMap := make(map[string]int)
	urlMap := make(map[string]int)
	for _, v := range result {
		val, ok := allUrlMap[v]
		if !ok {
			allUrlMap[v] = 1
		} else {
			allUrlMap[v] = val + 1
		}
	}

	//logger.Info(len(allUrlMap))
	for k, v := range allUrlMap {
		if v >= min && v <= max {
			urlMap[k] = v
		}
	}
	allHostMap := make(map[string]int)
	for k, v := range urlMap {
		parse, err := url.Parse(k)
		if err != nil {
			logger.Error(err)
			continue
		}
		h, ok := allHostMap[parse.Host]
		if !ok {
			allHostMap[parse.Host] = v
		} else {
			allHostMap[parse.Host] = h + v
		}
	}
	writeJsonRsp(writer, Val{UrlMap: urlMap, HostMap: allHostMap})
}

func UrlTimeAnalysis(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	timeType := query.Get("type")
	zR, err := Conn.ZRangeWithScores(RequestList, 0, -1).Result()
	if err != nil {
		writeJsonRsp(writer, err.Error())
		return
	}

	// 天 -> 域名 -> 次数
	dayMap := make(map[string]map[string]int)

	ids := parseIdArray(zR)

	// 批量得到url
	// id -> host
	uMap := make(map[string]string)
	array := splitArray(ids, 100)
	for _, ba := range array {
		result, err := Conn.HMGet(RequestUrlList, ba...).Result()
		if err != nil {
			continue
		}
		for i := range ba {
			path := result[i].(string)
			parse, err := url.Parse(path)
			if err != nil {
				continue
			}

			uMap[ba[i]] = parse.Host
		}
	}

	timeS := ctool.YYYY_MM_DD
	if timeType == "month" {
		timeS = ctool.YYYY_MM
	}
	// 按天聚合结果
	for _, val := range zR {
		cols := strings.Split(val.Member.(string), " ")
		id := cols[3]
		hitTime := time.UnixMicro(int64(val.Score) / 1000)
		hitDay := hitTime.Format(timeS)
		m, ok := dayMap[hitDay]
		if !ok {
			m = make(map[string]int)
			dayMap[hitDay] = m
		}

		u := uMap[id]
		c, ok := m[u]
		if !ok {
			m[u] = 1
		} else {
			m[u] = c + 1
		}
	}

	writeJsonRsp(writer, dayMap)
}

func parseIdArray(zR []redis.Z) []string {
	// 得到id
	var ids []string
	for _, val := range zR {
		cols := strings.Split(val.Member.(string), " ")
		id := cols[3]

		ids = append(ids, id)
	}
	return ids
}

func DetailById(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	idStr := query.Get("id")
	if idStr == "" {
		writeJsonParamError(writer, "id is required")
		return
	}
	targetMsg := getDetailByKey(idStr)
	writeJsonRsp(writer, DetailVo{
		Url:   targetMsg.Url,
		Start: targetMsg.ReqTime.Add(-time.Minute * 5).Format(ctool.YYYY_MM_DD_HH_MM),
		End:   targetMsg.ReqTime.Add(time.Minute * 5).Format(ctool.YYYY_MM_DD_HH_MM),
	})
}

func HostPerformance(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	idStr := query.Get("id")
	hostStr := query.Get("host")
	urlStr := query.Get("url")
	if hostStr == "" && urlStr == "" && idStr == "" {
		writeJsonParamError(writer, "host or url is required")
		return
	}

	if idStr != "" {
		targetMsg := getDetailByKey(idStr)
		parse, err := url.Parse(targetMsg.Url)
		if err != nil {
			writeJsonParamError(writer, err.Error())
			return
		}

		path := parse.Scheme + "://" + parse.Host + parse.Path
		rsp := queryData(targetMsg.ReqTime.Add(-time.Minute*5), targetMsg.ResTime.Add(time.Minute*5), "", path, 1)
		writeJsonRsp(writer, rsp)
		return
	}

	rsp := ctool.ResultVO[[]PerfPageVo]{}

	startStr := query.Get("start")
	endStr := query.Get("end")
	minStr := query.Get("min")
	min := 1
	if minStr != "" {
		atoi, err := strconv.Atoi(minStr)
		if err == nil {
			min = atoi
		}
	}
	if startStr != "" && endStr != "" {
		startStr = strings.Replace(startStr, "T", " ", 1)
		endStr = strings.Replace(endStr, "T", " ", 1)
		start, err := time.Parse(ctool.YYYY_MM_DD_HH_MM, startStr)
		if err != nil {
			writeJsonParamError(writer, err.Error())
			return
		}
		end, err := time.Parse(ctool.YYYY_MM_DD_HH_MM, endStr)
		if err != nil {
			writeJsonParamError(writer, err.Error())
			return
		}
		rsp = queryData(start.Add(-time.Hour*8), end.Add(-time.Hour*8), hostStr, urlStr, min)
	} else {
		rsp.Msg = "invalid param"
		rsp.Code = 400
	}
	writeJsonRsp(writer, rsp)
}

func queryData(start time.Time, end time.Time, hostStr string, urlStr string, min int) ctool.ResultVO[[]PerfPageVo] {
	var result []PerfPageVo

	rsp := ctool.ResultVO[[]PerfPageVo]{}

	// 为什么会有时区问题
	zr, err := Conn.ZRangeByScoreWithScores(RequestList, redis.ZRangeBy{
		Min: fmt.Sprint(start.UnixNano()),
		Max: fmt.Sprint(end.UnixNano()),
	}).Result()
	if err != nil {
		return ctool.FailedWithMsg[[]PerfPageVo](err.Error())
	}

	ids := parseIdArray(zr)
	array := splitArray(ids, 300)
	var cache []*ReqLog[Message]
	for _, batch := range array {
		result, err := Conn.HMGet(RequestUrlList, batch...).Result()
		if err != nil {
			continue
		}
		for i := range batch {
			path := result[i].(string)
			if (hostStr != "" && (strings.HasPrefix(path, "http://"+hostStr) || strings.HasPrefix(path, "https://"+hostStr))) ||
				(urlStr != "" && strings.HasPrefix(path, urlStr)) {
				msg := getDetailByKey(batch[i])
				cache = append(cache, msg)
			}
		}
	}

	stream.Just(cache...).Filter(func(item any) bool {
		msg := item.(*ReqLog[Message])
		if msg == nil {
			return false
		}
		_, err := url.Parse(msg.Url)
		if err != nil {
			return false
		}
		return true
	}).Group(func(item any) any {
		msg := item.(*ReqLog[Message])
		parseUrl, _ := url.Parse(msg.Url)
		return parseUrl.Host + "" + parseUrl.Path
	}).ForEach(func(item any) {
		groupItem := item.(stream.GroupItem)
		val := groupItem.Val
		if len(val) < min {
			return
		}

		ele := PerfPageVo{Url: groupItem.Key.(string)}
		var ts []int
		var reqList []int64
		for _, v := range val {
			msg := v.(*ReqLog[Message])
			ms := msg.ResTime.Sub(msg.ReqTime).Milliseconds()
			ele.Tall += int(ms)
			ts = append(ts, int(ms))
			reqList = append(reqList, msg.ReqTime.Unix())
		}
		sort.Slice(ts, func(i, j int) bool {
			return ts[i] < ts[j]
		})
		ele.TP30 = ts[int(float32(len(ts))*0.3)]
		ele.TP50 = ts[int(float32(len(ts))*0.5)]
		ele.TP90 = ts[int(float32(len(ts))*0.9)]
		ele.TP95 = ts[int(float32(len(ts))*0.95)]
		ele.TP99 = ts[int(float32(len(ts))*0.99)]

		ele.TAvg = ele.Tall / len(val)
		ele.Tct = len(val)

		stream.Just(reqList...).Group(func(item any) any {
			return item.(int64)
		}).ForEach(func(item any) {
			groupItem := item.(stream.GroupItem)
			qps := len(groupItem.Val)
			if ele.Qps < qps {
				ele.Qps = qps
			}
		})
		result = append(result, ele)
	})

	rsp.Data = result
	rsp.Code = 0
	return rsp
}
