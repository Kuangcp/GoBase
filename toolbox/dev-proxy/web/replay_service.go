package web

import (
	"bytes"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type (
	BenchStat struct {
		Id       string `json:"id"`
		Total    int    `json:"total"`
		Complete int    `json:"complete"`
		Failed   int    `json:"failed"`
		Con      int    `json:"con"`

		Mills    int64  `json:"mills"`
		Duration string `json:"duration"`

		Qps string `json:"qps"`
		Rt  string `json:"rt"`

		RealMills    int64  `json:"real_mills"` // 实际耗时
		RealDuration string `json:"real_duration"`
		Start        string `json:"start"`
	}
	ReplayCtx struct {
		cmd string
		msg *core.ReqLog[core.Message]
		act func() bool
	}
)

var (
	client     http.Client
	clientList []*http.Client
	cliCnt     = 100
	prf        func(*http.Request) (*url.URL, error)
)

func InitClient() {
	proxy := fmt.Sprintf("http://127.0.0.1:%v", core.Port)
	uri, _ := url.Parse(proxy)
	proxyConf := http.ProxyURL(uri)
	prf = proxyConf
	client = http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: proxyConf,
			//MaxConnsPerHost:     50,
			//MaxIdleConns:        20,
			//MaxIdleConnsPerHost: 50,
			//IdleConnTimeout:     60 * time.Second,
		},
	}

	//clientList = make([]*http.Client, cliCnt)
	//for i := 0; i < cliCnt; i++ {
	//	clientList[i] = &http.Client{
	//		Timeout: time.Second * 30,
	//		Transport: &http.Transport{
	//			Proxy:               proxyConf,
	//			MaxConnsPerHost:     0,
	//			MaxIdleConns:        100,
	//			MaxIdleConnsPerHost: 20,
	//			IdleConnTimeout:     60 * time.Second,
	//		},
	//	}
	//}
}

func getClient() *http.Client {
	return clientList[rand.Int()%cliCnt]
}

// TODO 并发数多的时候 请求的延迟会剧烈上升，curl开启大量文件响应慢？
func curlReplayCtx(id string) *ReplayCtx {
	command := buildCommandById(id, true, false)

	return &ReplayCtx{cmd: command, act: func() bool {
		_, success := core.ExecCommand(command)
		return success
	}}
}

// FIXME 当并发数超过10后当前主机6核12线程 对比 wrk ab 等工具 压测出的QPS误差非常大
func httpReplayCtx(id string) *ReplayCtx {
	debug := false
	detail := core.GetDetailByKey(id)
	return &ReplayCtx{msg: detail, act: func() bool {
		var request *http.Request
		var err error

		u, err := url.Parse(detail.Url)
		if err != nil {
			logger.Error(err)
			return false
		}
		tUrl := u.Scheme + "://" + u.Host + u.Path
		if u.RawQuery != "" {
			tUrl += "?" + url.QueryEscape(u.RawQuery)
		}
		if len(detail.Request.Body) > 0 {
			reader := bytes.NewReader(detail.Request.Body)
			request, err = http.NewRequest(detail.Method, tUrl, reader)
		} else {
			request, err = http.NewRequest(detail.Method, tUrl, nil)
		}

		for k, v := range detail.Request.Header {
			for _, i := range v {
				request.Header.Set(k, i)
			}
		}
		// 直接复制省内存但是有数据竞争问题
		//request.Header = detail.Request.Header
		request.Header[core.HeaderProxyBench] = []string{"1"}

		resp, err := client.Do(request)
		//resp, err := getClient().Do(request)
		if resp == nil || err != nil {
			logger.Error(err)
			return false
		}

		if debug {
			rspBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return false
			}
			logger.Info(string(rspBody))
			logger.Info(resp.Header)
		}

		return true
	}}
}

// TODO 并发低，延迟高，
// https://github.com/tsliwowicz/go-wrk
// https://github.com/adjust/go-wrk
func BenchRequest(request *http.Request) ctool.ResultVO[*BenchStat] {
	var data struct {
		Id    string `form:"id"`
		Con   int    `form:"con"`
		Total int    `form:"total"`
	}
	if err := ctool.Unpack(request, &data); err != nil {
		return ctool.FailedWithMsg[*BenchStat](err.Error())
	}
	if data.Con <= 1 {
		data.Con = 1
	}
	if data.Total <= 1 {
		data.Total = 1
	}

	//ctx := curlReplayCtx(data.Id)
	ctx := httpReplayCtx(data.Id)

	// TODO pool 资源回收
	pool, _ := sizedpool.NewQueuePool(data.Con)
	lock := &sync.Mutex{}
	startTime := time.Now()
	startTimeMs := startTime.UnixMilli()
	stat := &BenchStat{}
	for i := 0; i < data.Total; i++ {
		pool.Submit(func() {
			start := time.Now()
			success := ctx.act()
			end := time.Now()
			waste := end.Sub(start)

			lock.Lock()
			if !success {
				stat.Failed += 1
			} else {
				stat.Complete += 1
			}
			stat.Mills += waste.Milliseconds()
			//fmt.Println(waste)
			stat.Total += 1
			lock.Unlock()
		})
	}
	pool.Wait()
	pool.Close()

	stat.Start = startTime.Format(ctool.YYYY_MM_DD_HH_MM_SS)
	stat.RealMills = time.Now().UnixMilli() - startTimeMs
	stat.RealDuration = (time.Duration(stat.RealMills) * time.Millisecond).String()
	stat.Id = data.Id
	stat.Duration = (time.Duration(stat.Mills) * time.Millisecond).String()
	if stat.RealMills > 0 {
		stat.Qps = fmt.Sprint(int64(stat.Total*1000) / stat.RealMills)
	}
	if stat.Total > 0 {
		stat.Rt = (time.Duration(stat.Mills/int64(stat.Total)) * time.Millisecond).String()
	}
	stat.Con = data.Con

	return ctool.SuccessWith[*BenchStat](stat)
}

func ReplayRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	idx := query.Get("idx")
	id := query.Get("id")
	selfProxy := query.Get("selfProxy")
	if idx != "" && id == "" {
		sortIdx, _ := strconv.Atoi(idx)
		result, err := core.Conn.ZRange(core.RequestList, int64(sortIdx-1), int64(sortIdx-1)).Result()
		if err != nil {
			logger.Error(err)
			return
		}
		if len(result) == 0 {
			return
		}
		id = core.ConvertToDbKey(result[0])
	}

	command := buildCommandById(id, selfProxy == "Y", true)
	if command == "" {
		core.RspStr(writer, id+" not found")
		return
	}
	logger.Info("Replay ", id)
	result, success := core.ExecCommand(command)
	if !success {
		core.RspStr(writer, "ERROR: \n"+command+"\n"+result+"\n")
	} else {
		core.RspStr(writer, result)
	}
}
