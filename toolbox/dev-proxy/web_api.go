package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func startQueryServer() {
	logger.Info("Start query server on 127.0.0.1:%d", queryPort)

	http.HandleFunc("/list", pageListReqHistory)
	http.HandleFunc("/curlCommand", buildCurlCommand)
	http.HandleFunc("/replayReq", replayRequest)
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
	sort, _ := strconv.Atoi(id)

	commandList := buildCommandBySort(sort, selfProxy)
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
	sort, _ := strconv.Atoi(id)

	res := buildCommandBySort(sort, selfProxy)
	if res == nil {
		return
	}
	writer.Write([]byte(strings.Join(res, "\n\n")))
}

func buildCommandBySort(sort int, selfProxy string) []string {
	result, err := connection.ZRange(TotalReq, int64(sort-1), int64(sort-1)).Result()
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
			cmd += " -x 127.0.0.1:1234 "
		}
		parseUrl, _ := url.Parse(detail.Url)

		query := url.PathEscape(parseUrl.RawQuery)
		query = strings.ReplaceAll(query, "&", "\\&")
		cmd += parseUrl.Scheme + "://" + parseUrl.Host + parseUrl.Path + "\\?" + query

		for k, vl := range detail.Request.Header {
			for _, v := range vl {
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
