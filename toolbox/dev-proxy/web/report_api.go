package web

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/app"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type (
	PerfPage struct {
		list []PerfPageVo
	}
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
	var param struct {
		Min   int
		Max   int
		Depth int
	}
	if err := ctool.Unpack(request, &param); err != nil {
		core.WriteJsonParamError(writer, err.Error())
		return
	}

	if param.Min < 1 {
		param.Min = 50
	}
	if param.Max < 1 {
		param.Max = 100
	}

	result, err := core.Conn.HGetAll(core.RequestUrlList).Result()
	if err != nil {
		logger.Error(err)
		core.WriteJsonRsp(writer, err.Error())
		return
	}

	type Val struct {
		UrlMap  map[string]int
		HostMap map[string]int
	}

	allUrlMap := make(map[string]int)
	urlMap := make(map[string]int)
	for _, v := range result {
		if param.Depth != 0 {

			parse, err := url.Parse(v)
			if err != nil {
				logger.Error(err)
				continue
			}
			path := parse.Path
			parts := strings.Split(path, "/")
			depth := param.Depth
			depth += 1
			if len(parts) < depth {
				depth = len(parts)
			}
			newPath := strings.Join(parts[:depth], "/")
			v = parse.Scheme + "://" + parse.Host + newPath
		}

		val, ok := allUrlMap[v]
		if !ok {
			allUrlMap[v] = 1
		} else {
			allUrlMap[v] = val + 1
		}
	}

	//logger.Info(len(allUrlMap))
	for k, v := range allUrlMap {
		if v >= param.Min && v <= param.Max {
			urlMap[k] = v
		}
	}
	hostMap := make(map[string]int)
	for k, v := range urlMap {
		parse, err := url.Parse(k)
		if err != nil {
			logger.Error(err)
			continue
		}
		h, ok := hostMap[parse.Host]
		if !ok {
			hostMap[parse.Host] = v
		} else {
			hostMap[parse.Host] = h + v
		}
	}
	core.WriteJsonRsp(writer, Val{UrlMap: urlMap, HostMap: hostMap})
}

func UrlTimeAnalysis(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	timeType := query.Get("type")
	zR, err := core.Conn.ZRangeWithScores(core.RequestList, 0, -1).Result()
	if err != nil {
		core.WriteJsonRsp(writer, err.Error())
		return
	}

	// 天 -> 域名 -> 次数
	dayMap := make(map[string]map[string]int)

	ids := parseIdArray(zR)

	// 批量得到url
	// id -> host
	uMap := make(map[string]string)
	batch := 100
	array := splitArray(ids, batch)
	for bi, ba := range array {
		result, err := core.Conn.HMGet(core.RequestUrlList, ba...).Result()
		if err != nil {
			continue
		}
		for i := range ba {
			tmp := result[i]
			if tmp == nil {
				logger.Warn("NOT MATCH", zR[bi*batch+i], ba[i])
				core.RemoveReqMember(zR[bi*batch+i].Member)
				continue
			}
			path := tmp.(string)
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

	core.WriteJsonRsp(writer, dayMap)
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

func ConnNum(writer http.ResponseWriter, _ *http.Request) {
	var num int32 = 0
	if app.ProxyHandler != nil {
		num = app.ProxyHandler.ClientConnNum()
	}
	writer.Write([]byte(fmt.Sprint(num)))
}

func DetailById(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	idStr := query.Get("id")
	if idStr == "" {
		core.WriteJsonParamError(writer, "id is required")
		return
	}
	targetMsg := core.GetDetailByKey(idStr)
	core.WriteJsonRsp(writer, DetailVo{
		Url:   targetMsg.Url,
		Start: targetMsg.ReqTime.Add(-time.Minute * 5).Format(ctool.YYYY_MM_DD_HH_MM),
		End:   targetMsg.ReqTime.Add(time.Minute * 5).Format(ctool.YYYY_MM_DD_HH_MM),
	})
}

func HostPerformance(request *http.Request) ctool.ResultVO[[]PerfPageVo] {
	if core.Conn == nil {
		return ctool.FailedWithMsg[[]PerfPageVo]("empty")
	}
	var param struct {
		Id    string     `form:"id"`
		Host  string     `form:"host"`
		Url   string     `form:"url"`
		Start *time.Time `form:"start" fmt:"2006-01-02T15:04"`
		End   *time.Time `form:"end" fmt:"2006-01-02T15:04"`
		Min   int
	}

	if err := ctool.Unpack(request, &param); err != nil {
		return ctool.FailedWithMsg[[]PerfPageVo](err.Error())
	}

	if param.Host == "" && param.Url == "" && param.Id == "" {
		return ctool.FailedWithMsg[[]PerfPageVo]("host or url is required")
	}

	if param.Id != "" {
		targetMsg := core.GetDetailByKey(param.Id)
		parse, err := url.Parse(targetMsg.Url)
		if err != nil {
			return ctool.FailedWithMsg[[]PerfPageVo](err.Error())
		}

		fullUrl := parse.Scheme + "://" + parse.Host + parse.Path
		start := targetMsg.ReqTime.Add(-time.Minute * 5)
		end := targetMsg.ResTime.Add(time.Minute * 5)
		return queryData(start, end, "", fullUrl, 1)
	}

	rsp := ctool.ResultVO[[]PerfPageVo]{}
	if param.Start != nil && param.End != nil {
		rsp = queryData(*param.Start, *param.End, param.Host, param.Url, param.Min)
	} else {
		rsp.Msg = "invalid param"
		rsp.Code = 400
	}
	return rsp
}

func queryData(start time.Time, end time.Time, hostStr string, urlStr string, min int) ctool.ResultVO[[]PerfPageVo] {
	var result []PerfPageVo

	rsp := ctool.ResultVO[[]PerfPageVo]{}

	//logger.Info("query data")
	zr, err := core.Conn.ZRangeByScoreWithScores(core.RequestList, redis.ZRangeBy{
		Min: fmt.Sprint(start.UnixNano()),
		Max: fmt.Sprint(end.UnixNano()),
	}).Result()
	if err != nil {
		return ctool.FailedWithMsg[[]PerfPageVo](err.Error())
	}

	ids := parseIdArray(zr)
	batchSize := 300
	array := splitArray(ids, batchSize)
	var cache []*core.ReqLog[core.Message]

	// TODO memory leak?
	pool, _ := sizedpool.NewTmpFuturePool(sizedpool.PoolOption{Size: 20, Timeout: time.Second * 20})

	//logger.Info("query from leveldb")
	var tasks []*sizedpool.FutureTask
	for bi, batch := range array {
		result, err := core.Conn.HMGet(core.RequestUrlList, batch...).Result()
		if err != nil {
			logger.Error(err)
			continue
		}

		for i := range batch {
			tmp := result[i]
			if tmp == nil {
				logger.Warn("NOT MATCH", zr[bi*batchSize+i])
				core.RemoveReqMember(zr[bi*batchSize+i].Member)
				continue
			}
			fullURL := tmp.(string)
			matchHost := hostStr != "" && (strings.HasPrefix(fullURL, "http://"+hostStr) ||
				strings.HasPrefix(fullURL, "https://"+hostStr))
			matchUrl := urlStr != "" && (strings.HasPrefix(fullURL, urlStr) ||
				strings.HasPrefix(fullURL, "http://"+urlStr) ||
				strings.HasPrefix(fullURL, "https://"+urlStr))
			if matchHost || matchUrl {
				key := batch[i]
				future := pool.SubmitFuture(sizedpool.Callable{
					ActionFunc: func(ctx context.Context) (interface{}, error) {
						msg := core.GetDetailByKey(key)
						// 提早释放内存
						msg.Request = core.Message{}
						msg.Response = core.Message{}
						return msg, nil
					},
				})
				tasks = append(tasks, future)
			}
		}
	}
	pool.Wait()
	for _, f := range tasks {
		data, err := f.GetData()
		if err != nil {
			logger.Error(err)
			continue
		}
		cache = append(cache, data.(*core.ReqLog[core.Message]))
	}
	pool.Close()

	stream.Just(cache...).Filter(func(item any) bool {
		msg := item.(*core.ReqLog[core.Message])
		if msg == nil {
			return false
		}
		_, err := url.Parse(msg.Url)
		if err != nil {
			return false
		}
		return true
	}).Group(func(item any) any {
		msg := item.(*core.ReqLog[core.Message])
		parseUrl, _ := url.Parse(msg.Url)
		return parseUrl.Host + "" + parseUrl.Path
	}).ForEach(func(item any) {
		groupItem := item.(stream.GroupItem)
		val := groupItem.Val
		if len(val) < min {
			return
		}

		showUrl := groupItem.Key.(string)
		idx := strings.Index(showUrl, "/")
		if idx != -1 {
			showUrl = showUrl[idx:]
		}
		ele := PerfPageVo{Url: showUrl}
		var ts []int
		var reqList []int64
		for _, v := range val {
			msg := v.(*core.ReqLog[core.Message])
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
	sort.Slice(result, func(i, j int) bool {
		return result[i].Url < result[j].Url
	})

	//logger.Info("group by")

	rsp.Data = result
	rsp.Code = 0
	return rsp
}
