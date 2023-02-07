package main

import (
	_ "embed"
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"sort"
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

func buildCommandById(id, selfProxy string) string {
	detail := getDetailByKey(id)
	if detail == nil {
		return ""
	}
	cmd := "curl "
	if selfProxy == "Y" {
		cmd += fmt.Sprintf(" -x 127.0.0.1:%v ", port)
	}
	parseUrl, _ := url.Parse(detail.Url)
	cmd += parseUrl.Scheme + "://" + parseUrl.Host + parseUrl.Path
	if parseUrl.RawQuery != "" {
		query := url.PathEscape(parseUrl.RawQuery)
		query = strings.ReplaceAll(query, "&", "\\&")
		cmd += "\\?" + query
	}
	var key []string
	for k := range detail.Request.Header {
		key = append(key, k)
	}
	sort.Strings(key)
	for _, k := range key {
		val := detail.Request.Header.Values(k)
		for _, v := range val {
			cmd += fmt.Sprintf(" -H '%s: %s'", k, v)
		}
	}

	if len(detail.Request.Body) > 0 {
		cmd += fmt.Sprintf(" --data-raw $'%s'", string(detail.Request.Body))
	}
	//logger.Info(cmd)

	return cmd
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
