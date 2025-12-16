package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	http.HandleFunc("/list", pageListReqHistory)
	http.HandleFunc("/curl", buildCurlCommand)
	http.HandleFunc("/replay", replayRequest)
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

func replayRequest(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	selfProxy := query.Get("selfProxy")
	sortIdx, _ := strconv.Atoi(id)

	commandList := buildCommandBySort(sortIdx, selfProxy)
	if commandList == nil {
		writer.Write([]byte(id + " not found"))
		return
	}
	for _, c := range commandList {
		result, success := execCommand(c)
		if !success {
			writer.Write([]byte("ERROR: \n" + c + "\n" + result + "\n"))
		} else {
			writer.Write([]byte(result))
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

		if detail.Request.Body != "" {
			cmd += fmt.Sprintf(" --data-raw $'%s'", detail.Request.Body)
		}
		//logger.Info(cmd)

		cmdList = append(cmdList, cmd)
	}
	return cmdList
}

func pageListReqHistory(writer http.ResponseWriter, request *http.Request) {
	values := request.URL.Query()
	page := values.Get("page")
	size := values.Get("size")
	pageResult := pageQueryReqLog(page, size)
	result := ResultVO[*PageVO[*ReqLog[MessageVO]]]{}
	if pageResult == nil {
		result.Code = 101
		result.Msg = "invalid data"
	} else {
		result.Code = 0
		result.Data = pageResult

		hiddenHeaderEachLog(pageResult)
	}

	writer.Header().Set("Content-Type", "application/json")
	buffer := toJSONBuffer(result)
	writer.Write(buffer.Bytes())
}

func hiddenHeaderEachLog(pageResult *PageVO[*ReqLog[MessageVO]]) {
	if pageResult.Data == nil {
		return
	}
	for _, v := range pageResult.Data {
		hiddenHeader(v.Request.Header)
		hiddenHeader(v.Response.Header)
	}
}

func hiddenHeader(header http.Header) {
	delete(header, "User-Agent")
	delete(header, "Accept-Encoding")
	delete(header, "Referer")
	delete(header, "Cache-Control")
	delete(header, "Accept-Language")
	delete(header, "Pragma")
}
