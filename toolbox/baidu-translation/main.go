package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/cuibase"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	Description: "Format markdown file, generate catalog",
	Version:     "1.0.0",
	VerbLen:     -3,
	ParamLen:    -5,
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

func query(query string, fromLang string, toLang string, appId string, secretKey string) {
	if len(query) == 0 || len(appId) == 0 || len(secretKey) == 0 || len(fromLang) == 0 || len(toLang) == 0 {
		log.Fatalf("%v Param error %v", cuibase.Red, cuibase.End)
	}

	urls := BaiduApi + "?from=" + fromLang + "&to=" + toLang
	urls += "&appid=" + appId
	urls += "&q=" + url.QueryEscape(query)
	salt := strconv.Itoa(rand.Intn(65535))
	urls += "&salt=" + salt
	urls += "&sign=" + fmt.Sprintf("%x", md5.Sum([]byte(appId+query+salt+secretKey)))

	resp, err := http.Get(urls)
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
