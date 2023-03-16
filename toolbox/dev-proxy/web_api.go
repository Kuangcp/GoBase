package main

import (
	_ "embed"
	"fmt"
	"github.com/kuangcp/logger"
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
var indexPage string

//go:embed static/monokai-sublime.min.css
var sublimeStyle string

//go:embed static/highlight.min.js
var highlightJs string

//go:embed static/json.min.js
var jsonJs string

const (
	jsT   = "text/javascript; charset=utf-8"
	cssT  = "text/css; charset=utf-8"
	htmlT = "text/html; charset=utf-8"
)

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	if debug {
		http.Handle("/", http.FileServer(http.Dir("./static")))
	} else {
		http.HandleFunc("/", bindStatic(indexPage, htmlT))
		http.HandleFunc("/monokai-sublime.min.css", bindStatic(sublimeStyle, cssT))
		http.HandleFunc("/highlight.min.js", bindStatic(highlightJs, jsT))
		http.HandleFunc("/json.min.js", bindStatic(jsonJs, jsT))
	}

	http.HandleFunc("/list", handleInterceptor(JSONFunc(pageListReqHistory)))
	http.HandleFunc("/curl", handleInterceptor(buildCurlCommandApi))
	http.HandleFunc("/replay", replayRequest)
	http.HandleFunc("/del", delRequest)
	http.HandleFunc("/flushAll", handleInterceptor(flushAllData))

	http.ListenAndServe(fmt.Sprintf(":%v", queryPort), nil)
}

func (p PageQueryParam) buildStartEnd() (int64, int64) {
	return int64((p.page - 1) * p.size), int64(p.page*p.size) - 1
}

func handleInterceptor(h http.HandlerFunc) http.HandlerFunc {
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
