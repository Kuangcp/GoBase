package core

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	PageQueryParam struct {
		page int
		size int
		id   string
		kwd  string
		date string
	}
)

//go:embed static/index.html
var indexHtml string

//go:embed static/monokai-sublime.min.css
var sublimeCss string

//go:embed static/highlight.min.js
var highlightJs string

//go:embed static/json.min.js
var jsonMinJs string

//go:embed static/main.js
var mainJs string

//go:embed static/main.css
var mainCss string

//go:embed static/favicon.ico
var icon string

const (
	jsT   = "text/javascript; charset=utf-8"
	cssT  = "text/css; charset=utf-8"
	htmlT = "text/html; charset=utf-8"
	iconT = "image/vnd.microsoft.icon"
)

func StartQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", QueryPort)

	if Debug {
		http.Handle("/", http.FileServer(http.Dir("./static")))
	} else {
		http.HandleFunc("/", bindStatic(indexHtml, htmlT))
		http.HandleFunc("/favicon.ico", bindStatic(icon, iconT))
		http.HandleFunc("/monokai-sublime.min.css", bindStatic(sublimeCss, cssT))
		http.HandleFunc("/main.css", bindStatic(mainCss, cssT))

		http.HandleFunc("/main.js", bindStatic(mainJs, jsT))
		http.HandleFunc("/highlight.min.js", bindStatic(highlightJs, jsT))
		http.HandleFunc("/json.min.js", bindStatic(jsonMinJs, jsT))

	}

	http.HandleFunc("/list", rtTimeInterceptor(JSONFunc(pageListReqHistory)))
	http.HandleFunc("/curl", rtTimeInterceptor(buildCurlCommandApi))
	http.HandleFunc("/replay", replayRequest)
	http.HandleFunc("/del", rtTimeInterceptor(delRequest))
	http.HandleFunc("/urlFrequency", rtTimeInterceptor(urlFrequencyApi))
	http.HandleFunc("/urlTimeAnalysis", rtTimeInterceptor(urlTimeAnalysis))
	http.HandleFunc("/uploadCache", rtTimeInterceptor(uploadCacheApi))
	http.HandleFunc("/flushAll", rtTimeInterceptor(flushAllData))

	http.ListenAndServe(fmt.Sprintf(":%v", QueryPort), nil)
}

func (p PageQueryParam) buildStartEnd() (int64, int64) {
	return int64((p.page - 1) * p.size), int64(p.page*p.size) - 1
}

func rtTimeInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	result, err := connection.ZRange(RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, key := range result {
		db.Delete([]byte(convertToDbKey(key)), nil)
	}

	connection.Del(RequestList)
	logger.Info("delete: ", len(result))
	RspStr(writer, "OK")
}

func uploadCacheApi(writer http.ResponseWriter, request *http.Request) {
	iterator := db.NewIterator(nil, nil)
	for iterator.Next() {
		bts := iterator.Value()
		var l ReqLog[Message]
		err := json.Unmarshal(bts, &l)
		if err != nil {
			logger.Error("key:["+string(iterator.Key())+"] GET ERROR:", err)
			continue
		}
		connection.ZAdd(RequestList, redis.Z{Member: l.CacheId, Score: float64(l.ReqTime.UnixNano())})
	}
	writeJsonRsp(writer, "OK")
}

// TODO 按域名 天 维度 统计访问频次
func urlTimeAnalysis(writer http.ResponseWriter, request *http.Request) {

}

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

	result, err := connection.ZRevRange(RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		writeJsonRsp(writer, err.Error())
		return
	}

	countMap := make(map[string]int)
	resultMap := make(map[string]int)
	for _, key := range result {
		log := getDetailByKey(convertToDbKey(key))
		if log == nil {
			logger.Warn(key)
			connection.ZRem(RequestList, key)
			continue
		}
		val, ok := countMap[log.Url]
		if !ok {
			countMap[log.Url] = 1
		} else {
			countMap[log.Url] = val + 1
		}
	}
	logger.Info(len(countMap))
	for k, v := range countMap {
		if v >= min && v <= max {
			resultMap[k] = v
		}
	}
	writeJsonRsp(writer, resultMap)
}

func delRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	path := query.Get("path")

	if id != "" {
		deleteById(writer, id)
		return
	}

	if path != "" {
		deleteByPath(writer, path)
		logger.Info("start compact")
		db.CompactRange(util.Range{})
		return
	}

	writeJsonRsp(writer, "invalid param")
}

func deleteByPath(writer http.ResponseWriter, path string) {
	result, err := connection.ZRevRange(RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		writeJsonRsp(writer, err.Error())
		return
	}

	total := 0
	for _, key := range result {
		log := matchDetailByKeyAndKwd(convertToDbKey(key), path)
		if log == nil {
			continue
		}
		total++

		//logger.Info(log.Url)
		connection.ZRem(RequestList, log.CacheId)
		db.Delete([]byte(log.Id), nil)
		if total >= 5000 {
			writeJsonRsp(writer, "out of count")
			return
		}
	}

	writeJsonRsp(writer, Success("OK"))
}

func deleteById(writer http.ResponseWriter, id string) {
	detail := getDetailByKey(id)
	if detail == nil {
		writeJsonRsp(writer, id+" not exist")
		return
	}

	connection.ZRem(RequestList, detail.CacheId)
	db.Delete([]byte(id), nil)
	writeJsonRsp(writer, Success("OK"))
}

func replayRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	idx := query.Get("idx")
	id := query.Get("id")
	selfProxy := query.Get("selfProxy")
	if idx != "" && id == "" {
		sortIdx, _ := strconv.Atoi(idx)
		result, err := connection.ZRange(RequestList, int64(sortIdx-1), int64(sortIdx-1)).Result()
		if err != nil {
			logger.Error(err)
			return
		}
		if len(result) == 0 {
			return
		}
		id = convertToDbKey(result[0])
	}

	command := buildCommandById(id, selfProxy)
	if command == "" {
		RspStr(writer, id+" not found")
		return
	}
	logger.Info("Replay ", id)
	result, success := execCommand(command)
	if !success {
		RspStr(writer, "ERROR: \n"+command+"\n"+result+"\n")
	} else {
		RspStr(writer, result)
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
	RspStr(writer, res)
}

func parseParam(request *http.Request) *PageQueryParam {
	values := request.URL.Query()
	page := values.Get("idx")
	id := values.Get("id")
	size := values.Get("size")
	kwd := values.Get("kwd")
	date := values.Get("date")

	pageI, _ := strconv.Atoi(page)
	sizeI, _ := strconv.Atoi(size)
	if sizeI <= 0 {
		sizeI = 1
	}
	if pageI <= 1 {
		pageI = 1
	}
	if kwd != "" {
		kwd = strings.TrimSpace(kwd)
	}

	if date != "" {
		vals := strings.Split(date, "-")
		date = strings.Join(vals[1:], "-")
	}
	return &PageQueryParam{page: pageI, size: sizeI, id: id, kwd: kwd, date: date}
}

func bindStatic(s, contentType string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", contentType)
		RspStr(writer, s)
	}
}
