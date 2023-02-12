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
		page   int
		size   int
		kwd    string
		prefix string
	}
)

//go:embed index.html
var indexPage string

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	if debug {
		http.Handle("/", http.FileServer(http.Dir(".")))
	} else {
		http.HandleFunc("/", searchPage)

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

func flushAllData(_ http.ResponseWriter, _ *http.Request) {
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
}

func delRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	if id != "" {
		connection.ZRem(RequestList, id)
		db.Delete([]byte(id), nil)
		writeJsonRsp(writer, "OK")
	} else {
		writeJsonRsp(writer, id+" not exist")
	}
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
	size := values.Get("size")
	kwd := values.Get("kwd")
	prefix := values.Get("prefix")

	pageI, _ := strconv.Atoi(page)
	sizeI, _ := strconv.Atoi(size)
	if sizeI <= 0 {
		sizeI = 1
	}
	if pageI <= 1 {
		pageI = 1
	}

	return &PageQueryParam{page: pageI, size: sizeI, kwd: kwd, prefix: prefix}
}

func searchPage(writer http.ResponseWriter, request *http.Request) {
	RspStr(writer, indexPage)
}
