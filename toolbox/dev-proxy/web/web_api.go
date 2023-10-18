package web

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/arl/statsviz"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ratelimiter"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	PageQueryParam struct {
		Page int `form:"idx"`
		Size int
		Id   string
		Kwd  string
		Date *time.Time `fmt:"2006-01-02"`
	}
)

// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
//
//go:embed static/proxy.pac
var pacFile string

//go:embed static
var static embed.FS

const (
	pacT = "application/x-ns-proxy-autoconfig"
)

func StartQueryServer() {
	mux := http.NewServeMux()

	limiter := ratelimiter.CreateLeakyLimiter(3)
	logger.Info("Start query server on 127.0.0.1:%d", core.ApiPort)

	if core.Debug {
		logger.Warn("debug mode")
		//fs := http.FileServer(http.Dir("./core/static"))
		//http.Handle("/", http.StripPrefix("/", fs))

		mux.Handle("/", http.FileServer(http.Dir("./web/static")))
	} else {
		sub, err := fs.Sub(static, "static")
		if err != nil {
			panic(err)
		}
		mux.Handle("/", http.FileServer(http.FS(sub)))
		mux.HandleFunc(core.PacUrl, PacFileApi)
	}

	err := statsviz.Register(mux)
	if err != nil {
		logger.Error(err)
	}

	//mux.HandleFunc("/list", rtRateInt(Json(PageListReqHistory), limiter))
	mux.HandleFunc("/list", core.Json(PageListReqHistory))
	mux.HandleFunc("/curl", rtInt(buildCurlCommandApi))
	mux.HandleFunc("/replay", rtRateInt(replayRequest, limiter))

	mux.HandleFunc("/del", rtInt(delRequest))
	mux.HandleFunc("/uploadCache", rtInt(uploadCacheApi))
	mux.HandleFunc("/flushAll", rtInt(flushAllData))

	mux.HandleFunc("/queryConfig", rtInt(core.Json(QueryConfig)))
	mux.HandleFunc("/saveConfig", rtInt(SaveConfig))

	mux.HandleFunc("/urlFrequency", rtInt(UrlFrequencyApi))
	mux.HandleFunc("/urlTimeAnalysis", rtInt(UrlTimeAnalysis))
	mux.HandleFunc("/detailById", rtInt(DetailById))
	mux.HandleFunc("/hostPerf", rtInt(core.Json(HostPerformance)))
	mux.HandleFunc("/connNum", rtInt(ConnNum))
	mux.HandleFunc("/exit", rtInt(Exit))

	http.ListenAndServe(fmt.Sprintf(":%v", core.ApiPort), mux)
}

func (p PageQueryParam) buildStartEnd() (int64, int64) {
	return int64((p.Page - 1) * p.Size), int64(p.Page*p.Size) - 1
}

func rtInt(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UnixMilli()
		h(w, r)
		end := time.Now().UnixMilli()
		logger.Info(r.URL.Path, end-start, "ms")
	}
}

func rtRateInt(h http.HandlerFunc, limiter ratelimiter.RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.TryAcquire() {
			buffer := ctool.ToJSONBuffer(ctool.ResultVO[string]{Code: 101, Msg: "频繁请求"})
			w.Write(buffer.Bytes())
			return
		}
		limiter.Acquire()

		start := time.Now().UnixMilli()
		h(w, r)
		end := time.Now().UnixMilli()
		logger.Info(r.URL.Path, end-start, "ms")
	}
}

func convertToDbKey(key string) string {
	return strings.Split(key, "  ")[1]
}

func flushAllData(writer http.ResponseWriter, _ *http.Request) {
	result, err := core.Conn.ZRange(core.RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, key := range result {
		core.Leveldb.Delete([]byte(core.ConvertToDbKey(key)), nil)
	}

	core.Conn.Del(core.RequestList)
	logger.Info("delete: ", len(result))
	core.RspStr(writer, "OK")
}

// upload leveldb data to redis
func uploadCacheApi(writer http.ResponseWriter, request *http.Request) {
	iterator := core.Leveldb.NewIterator(nil, nil)
	for iterator.Next() {
		bts := iterator.Value()
		var l core.ReqLog[core.Message]
		err := json.Unmarshal(bts, &l)
		if err != nil {
			logger.Error("key:["+string(iterator.Key())+"] GET ERROR:", err)
			continue
		}
		core.Conn.ZAdd(core.RequestList, redis.Z{Member: l.CacheId, Score: float64(l.ReqTime.UnixNano())})
		core.Conn.HSet(core.RequestUrlList, l.Id, l.Url)
	}
	core.WriteJsonRsp(writer, "OK")
}

func splitArray(l []string, batch int) [][]string {
	var result [][]string
	var s []string
	for _, l := range l {
		if len(s) == batch {
			result = append(result, s)
			s = []string{}
		}
		s = append(s, l)
	}
	if len(s) > 0 {
		result = append(result, s)
	}
	return result
}

func delRequest(writer http.ResponseWriter, request *http.Request) {
	// id 精准删除
	query := request.URL.Query()
	id := query.Get("id")
	if id != "" {
		deleteById(writer, id)
		return
	}

	// 按路径模糊删除
	size := query.Get("size")
	sizeI, err := strconv.Atoi(size)
	if err != nil {
		sizeI = 1
	}
	path := query.Get("path")
	if path != "" {
		deleteByPath(writer, path, sizeI)
		logger.Info("start compact leveldb")
		core.Leveldb.CompactRange(util.Range{})
		logger.Info("finish compact")
		return
	}

	core.WriteJsonRsp(writer, "invalid param")
}

func deleteByPath(writer http.ResponseWriter, path string, size int) {
	result, err := core.Conn.ZRevRange(core.RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		core.WriteJsonRsp(writer, err.Error())
		return
	}

	total := 0
	for _, key := range result {
		log := core.MatchDetailByKeyAndKwd(core.ConvertToDbKey(key), path)
		if log == nil {
			continue
		}
		total++

		logger.Info(log.Url, log.CacheId, log.Id)
		core.RemoveReqMember(log.CacheId)
		core.RemoveReqUrlKey(log.Id)
		core.Leveldb.Delete([]byte(log.Id), nil)
		if total >= size {
			core.WriteJsonRsp(writer, fmt.Sprintf("out of count %v", size))
			return
		}
	}

	core.WriteJsonRsp(writer, ctool.SuccessWith(fmt.Sprintf("Finish delete: %v", total)))
}

func deleteById(writer http.ResponseWriter, id string) {
	detail := core.GetDetailByKey(id)
	if detail == nil {
		core.WriteJsonRsp(writer, id+" not exist")
		return
	}

	core.RemoveReqMember(detail.CacheId)
	core.RemoveReqUrlKey(core.ConvertToDbKey(detail.CacheId))
	core.Leveldb.Delete([]byte(id), nil)
	core.WriteJsonRsp(writer, ctool.SuccessWith("OK"))
}

func replayRequest(writer http.ResponseWriter, request *http.Request) {
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

	command := buildCommandById(id, selfProxy)
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

// PacFileApi 默认使用缺省文件，优先使用独立配置文件
func PacFileApi(writer http.ResponseWriter, request *http.Request) {
	fileBt, err := os.ReadFile(core.PacFilePath)
	if err != nil || fileBt == nil || len(fileBt) == 0 {
		logger.Error(err)
		bindStatic(pacFile, pacT)(writer, request)
	} else {
		core.RspStr(writer, string(fileBt))
	}
}

func buildCurlCommandApi(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	selfProxy := query.Get("selfProxy")

	res := buildCommandById(id, selfProxy)
	if res == "" {
		return
	}
	core.RspStr(writer, res)
}

func bindStatic(s, contentType string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", contentType)
		core.RspStr(writer, s)
	}
}
