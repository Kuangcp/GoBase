package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
)

const BaiduApi = "https://fanyi-api.baidu.com/api/trans/vip/translate"

type (
	ResultVO struct {
		From        string     `json:"from"`
		To          string     `json:"to"`
		TransResult []TransMap `json:"trans_result"`
	}
	TransMap struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	}
)

var info = cuibase.HelpInfo{
	Description: "Translation between Chinese and English By Baidu API",
	Version:     "1.0.1",
	VerbLen:     -3,
	ParamLen:    -21,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help info",
		}, {
			Verb:    "-ez",
			Param:   "appId secretKey query",
			Comment: "Translate en to zh",
			Handler: func(params []string) {
				cuibase.AssertParamCount(4, "param less")
				query(fmt.Sprintf("%v", params[4:]), "en", "zh", params[2], params[3])
			},
		}, {
			Verb:    "-ze",
			Param:   "appId secretKey query",
			Comment: "Translate zh to en",
			Handler: func(params []string) {
				cuibase.AssertParamCount(4, "param less")
				query(fmt.Sprintf("%v", params[4:]), "zh", "en", params[2], params[3])
			},
		},
	}}

func anyEmpty(value ...string) bool {
	if len(value) == 0 {
		return true
	}
	for _, s := range value {
		if len(s) == 0 {
			return true
		}
	}
	return false
}

func query(query string, fromLang string, toLang string, appId string, secretKey string) {
	if anyEmpty(query, fromLang, toLang, appId, secretKey) {
		log.Fatalln(cuibase.Red.Println(" Param error "))
	}

	salt := strconv.Itoa(rand.Intn(65535))

	queryStr := "?from=" + fromLang
	queryStr += "&to=" + toLang
	queryStr += "&appid=" + appId
	queryStr += "&q=" + url.QueryEscape(query)
	queryStr += "&salt=" + salt
	queryStr += "&sign=" + fmt.Sprintf("%x", md5.Sum([]byte(appId+query+salt+secretKey)))

	resp, err := http.Get(BaiduApi + queryStr)
	cuibase.CheckIfError(err)
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("resp status code:[%d]\n", resp.StatusCode)
	//fmt.Printf("resp body data:[%s]\n", string(bodyContent))

	var v ResultVO
	err = json.Unmarshal(bodyContent, &v)
	cuibase.CheckIfError(err)

	fmt.Printf("%s %v %s\n", cuibase.LightGreen, strings.Trim(v.TransResult[0].Dst, "[]"), cuibase.End)
}

func main() {
	cuibase.RunActionFromInfo(info, nil)
}
