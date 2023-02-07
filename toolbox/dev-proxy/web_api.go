package main

import (
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

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	http.HandleFunc("/", searchPage)
	http.HandleFunc("/list", handleInterceptor(JSONFunc(pageListReqHistory)))
	http.HandleFunc("/curl", buildCurlCommand)
	http.HandleFunc("/replay", replayRequest)
	http.HandleFunc("/del", delRequest)
	http.HandleFunc("/flushAll", flushAllData)

	http.ListenAndServe(fmt.Sprintf(":%v", queryPort), nil)
}

func flushAllData(_ http.ResponseWriter, _ *http.Request) {
	result, err := connection.ZRange(TotalReq, 0, -1).Result()
	if err != nil {
		logger.Error(err)
		return
	}

	for _, key := range result {
		db.Delete([]byte(key), nil)
	}

	connection.Del(TotalReq)
	logger.Info("delete: ", len(result))
}

func delRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	if id != "" {
		connection.ZRem(TotalReq, id)
		db.Delete([]byte(id), nil)
		writeJsonRsp(writer, "OK")
	} else {
		writeJsonRsp(writer, id+" not exist")
	}
}

func replayRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	selfProxy := query.Get("selfProxy")
	sortIdx, _ := strconv.Atoi(id)

	commandList := buildCommandBySort(sortIdx, selfProxy)
	if commandList == nil {
		RspStr(writer, id+" not found")
		return
	}
	for _, c := range commandList {
		result, success := execCommand(c)
		if !success {
			RspStr(writer, "ERROR: \n"+c+"\n"+result+"\n")
		} else {
			RspStr(writer, result)
		}
	}
}

func buildCurlCommand(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	selfProxy := query.Get("selfProxy")
	sortIdx, _ := strconv.Atoi(id)

	res := buildCommandBySort(sortIdx, selfProxy)
	if res == nil {
		return
	}
	writer.Write([]byte(strings.Join(res, "\n\n")))
	RspStr(writer, strings.Join(res, "\n\n"))
}

func buildCommandBySort(sortIdx int, selfProxy string) []string {
	result, err := connection.ZRange(TotalReq, int64(sortIdx-1), int64(sortIdx-1)).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}
	var cmdList []string
	for _, key := range result {
		detail := getDetailByKey(key)
		if detail == nil {
			continue
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

		cmdList = append(cmdList, cmd)
	}
	return cmdList
}

func parseParam(request *http.Request) *PageQueryParam {
	values := request.URL.Query()
	page := values.Get("page")
	size := values.Get("size")
	kwd := values.Get("kwd")
	prefix := values.Get("prefix")

	pageI, _ := strconv.Atoi(page)
	sizeI, _ := strconv.Atoi(size)
	if sizeI <= 0 {
		sizeI = 1
	}
	if pageI < 0 {
		return nil
	}

	return &PageQueryParam{page: pageI, size: sizeI, kwd: kwd, prefix: prefix}
}

func searchPage(writer http.ResponseWriter, request *http.Request) {
	RspStr(writer, "<!DOCTYPE html>\n<html lang=\"en\"><body>")
	RspStr(writer, "<form action=\"list\" style=\"text-align:center;\">"+
		"<input name=\"page\" style=\"width:60px;\" type=\"number\"/>"+
		"<input name=\"prefix\" style=\"width:100px;\"/>"+
		"<input name=\"kwd\" style=\"width:280px;\"></input><button>Search</button> </form>")
	RspStr(writer, "</body></html>")
}
