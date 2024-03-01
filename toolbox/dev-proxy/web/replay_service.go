package web

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"net/http"
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
	}
	ReplayCtx struct {
		cmd string
		msg *core.ReqLog[core.Message]
		act func() bool
	}
)

func curlReplayCtx(id string) *ReplayCtx {
	command := buildCommandById(id, true, false)

	return &ReplayCtx{cmd: command, act: func() bool {
		_, success := core.ExecCommand(command)
		return success
	}}
}

func httpReplayCtx(id string) *ReplayCtx {
	detail := core.GetDetailByKey(id)

	//client := http.Client{
	//	Transport: &http.Transport{},
	//}
	//_, err := client.Get(detail.Url)

	return &ReplayCtx{msg: detail, act: func() bool {
		_, err := http.Get(detail.Url)
		return err != nil
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

	ctx := curlReplayCtx(data.Id)
	//ctx := httpReplayCtx(data.Id)

	pool, _ := sizedpool.NewQueuePool(data.Con)
	lock := &sync.Mutex{}
	on := time.Now().UnixMilli()
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
			stat.Total += 1
			lock.Unlock()
		})
	}
	pool.Wait()
	stat.RealMills = time.Now().UnixMilli() - on
	stat.RealDuration = (time.Duration(stat.RealMills) * time.Millisecond).String()
	stat.Id = data.Id
	stat.Duration = (time.Duration(stat.Mills) * time.Millisecond).String()
	stat.Qps = fmt.Sprint(int64(stat.Total*1000) / stat.RealMills)
	stat.Rt = (time.Duration(stat.Mills/int64(stat.Total)) * time.Millisecond).String()
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
