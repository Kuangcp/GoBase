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

func urlFrequencyApi(writer http.ResponseWriter, request *http.Request) {
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

func urlTimeAnalysis(writer http.ResponseWriter, request *http.Request) {
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

func hostPerformance(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	hostStr := query.Get("host")
	urlStr := query.Get("url")
	if hostStr == "" && urlStr == "" {
		writeJsonRsp(writer, "host or url is required")
		return
	}

	type Val struct {
		Url  string
		Tct  int
		Tall int
		TAvg int
		TP30 int
		TP50 int
		TP90 int
		TP95 int
		TP99 int
	}
	var result []Val

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
		start, err := time.Parse(ctool.YYYY_MM_DD, startStr)
		if err != nil {
			writeJsonRsp(writer, err.Error())
			return
		}
		end, err := time.Parse(ctool.YYYY_MM_DD, endStr)
		if err != nil {
			writeJsonRsp(writer, err.Error())
			return
		}

		zr, err := Conn.ZRangeByScoreWithScores(RequestList, redis.ZRangeBy{Min: fmt.Sprint(start.UnixNano()), Max: fmt.Sprint(end.UnixNano())}).Result()
		if err != nil {
			writeJsonRsp(writer, err.Error())
			return
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

			ele := Val{Url: groupItem.Key.(string)}
			var ts []int
			for _, v := range val {
				msg := v.(*ReqLog[Message])
				ms := msg.ResTime.Sub(msg.ReqTime).Milliseconds()
				ele.Tall += int(ms)
				ts = append(ts, int(ms))
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

			result = append(result, ele)
		})

		writeJsonRsp(writer, result)
		return
	}

}
