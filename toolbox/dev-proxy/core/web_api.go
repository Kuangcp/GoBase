package core

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb/util"
	"net/http"
	"net/url"
	"os"
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

// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Proxy_servers_and_tunneling/Proxy_Auto-Configuration_PAC_file
//
//go:embed static/proxy.pac
var pacFile string

const (
	jsT   = "text/javascript; charset=utf-8"
	cssT  = "text/css; charset=utf-8"
	htmlT = "text/html; charset=utf-8"
	iconT = "image/vnd.microsoft.icon"
	pacT  = "application/x-ns-proxy-autoconfig"
)

func StartQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", ApiPort)

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
		http.HandleFunc(PacUrl, PacFileApi)
	}

	http.HandleFunc("/list", rtTimeInterceptor(JSONFunc(pageListReqHistory)))
	http.HandleFunc("/curl", rtTimeInterceptor(buildCurlCommandApi))
	http.HandleFunc("/replay", replayRequest)
	http.HandleFunc("/del", rtTimeInterceptor(delRequest))
	http.HandleFunc("/urlFrequency", rtTimeInterceptor(urlFrequencyApi))
	http.HandleFunc("/urlTimeAnalysis", rtTimeInterceptor(urlTimeAnalysis))
	http.HandleFunc("/uploadCache", rtTimeInterceptor(uploadCacheApi))
	http.HandleFunc("/flushAll", rtTimeInterceptor(flushAllData))

	http.ListenAndServe(fmt.Sprintf(":%v", ApiPort), nil)
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
	result, err := Conn.ZRange(RequestList, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, key := range result {
		db.Delete([]byte(convertToDbKey(key)), nil)
	}

	Conn.Del(RequestList)
	logger.Info("delete: ", len(result))
	RspStr(writer, "OK")
}

// upload leveldb data to redis
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
		Conn.ZAdd(RequestList, redis.Z{Member: l.CacheId, Score: float64(l.ReqTime.UnixNano())})
		Conn.HSet(RequestUrlList, l.Id, l.Url)
	}
	writeJsonRsp(writer, "OK")
}

func splitArray(l []string, batch int) [][]string {
	var result [][]string
	var s = []string{}
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

	// 得到id
	var ids []string
	for _, val := range zR {
		cols := strings.Split(val.Member.(string), " ")
		id := cols[3]

		ids = append(ids, id)
	}

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
		db.CompactRange(util.Range{})
		logger.Info("finish compact")
		return
	}

	writeJsonRsp(writer, "invalid param")
}

func deleteByPath(writer http.ResponseWriter, path string, size int) {
	result, err := Conn.ZRevRange(RequestList, 0, -1).Result()
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

		logger.Info(log.Url, log.CacheId, log.Id)
		Conn.ZRem(RequestList, log.CacheId)
		Conn.HDel(RequestUrlList, log.Id)
		db.Delete([]byte(log.Id), nil)
		if total >= size {
			writeJsonRsp(writer, "out of count")
			return
		}
	}

	writeJsonRsp(writer, Success(fmt.Sprintf("Finish delete: %v", total)))
}

func deleteById(writer http.ResponseWriter, id string) {
	detail := getDetailByKey(id)
	if detail == nil {
		writeJsonRsp(writer, id+" not exist")
		return
	}

	Conn.ZRem(RequestList, detail.CacheId)
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
		result, err := Conn.ZRange(RequestList, int64(sortIdx-1), int64(sortIdx-1)).Result()
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

// PacFileApi 默认使用缺省文件，优先使用独立配置文件
func PacFileApi(writer http.ResponseWriter, request *http.Request) {
	fileBt, err := os.ReadFile(pacFilePath)
	if err != nil || fileBt == nil || len(fileBt) == 0 {
		logger.Error(err)
		bindStatic(pacFile, pacT)(writer, request)
	} else {
		RspStr(writer, string(fileBt))
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
