package web

import (
	"bytes"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
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

		RealMills    int64  `json:"real_mills"`
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
	client http.Client
)

func InitClient() {
	proxy := fmt.Sprintf("http://127.0.0.1:%v", core.Port)
	uri, _ := url.Parse(proxy)
	client = http.Client{
		Timeout: time.Second * 30,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(uri),
		},
	}
}

// TODO 并发数多的时候 请求的延迟会快速上升，curl开启大量文件响应慢？
func curlReplayCtx(id string) *ReplayCtx {
	command := buildCommandById(id, true, false)

	return &ReplayCtx{cmd: command, act: func() bool {
		_, success := core.ExecCommand(command)
		return success
	}}
}

func httpReplayCtx(id string) *ReplayCtx {
	detail := core.GetDetailByKey(id)
	return &ReplayCtx{msg: detail, act: func() bool {
		var request *http.Request
		var err error
		if len(detail.Request.Body) > 0 {
			reader := bytes.NewReader(detail.Request.Body)
			request, err = http.NewRequest(detail.Method, detail.Url, reader)
		} else {
			request, err = http.NewRequest(detail.Method, detail.Url, nil)
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
		if resp == nil || err != nil {
			logger.Error(err)
			return false
		}

		//rspBody, err := io.ReadAll(resp.Body)
		//if err != nil {
		//	return false
		//}
		//logger.Info(string(rspBody))
		//logger.Info(resp.Header)

		return true
	}}
}

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

	pool, _ := sizedpool.NewQueuePool(data.Con)
	lock := &sync.Mutex{}
	startTime := time.Now()
	startTimeMs := startTime.UnixMilli()
	stat := &BenchStat{}
	for i := 0; i < data.Total; i++ {
		pool.Submit(func() {
			start := time.Now().UnixMilli()
			success := ctx.act()
			end := time.Now().UnixMilli()
			waste := end - start

			lock.Lock()
			if !success {
				stat.Failed += 1
			} else {
				stat.Complete += 1
			}
			stat.Mills += waste
			//fmt.Println(waste)
			stat.Total += 1
			lock.Unlock()
		})
	}
	pool.Wait()
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
