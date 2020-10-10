package main

import (
	"crypto/md5"
	"encoding/hex"
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
	QueryParam struct {
		Query     string
		From      string
		To        string
		App       string
		SecretKey string
	}
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

func (t *QueryParam) buildMD5(salt string) string {
	h := md5.New()
	h.Write([]byte(t.App + t.Query + salt + t.SecretKey))
	return hex.EncodeToString(h.Sum(nil))
}

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
				param := QueryParam{
					Query:     fmt.Sprintf("%v", params[4:]),
					From:      "en",
					To:        "zh",
					App:       params[2],
					SecretKey: params[3],
				}
				query(param)
			},
		}, {
			Verb:    "-ze",
			Param:   "appId secretKey query",
			Comment: "Translate zh to en",
			Handler: func(params []string) {
				cuibase.AssertParamCount(4, "param less")
				param := QueryParam{
					Query:     fmt.Sprintf("%v", params[4:]),
					From:      "zh",
					To:        "en",
					App:       params[2],
					SecretKey: params[3],
				}
				query(param)
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

func query(param QueryParam) {
	if anyEmpty(param.Query, param.From, param.To, param.App, param.SecretKey) {
		log.Fatalln(cuibase.Red.Println(" Param error "))
	}

	salt := strconv.Itoa(rand.Intn(65535))

	queryStr := "?from=" + param.From
	queryStr += "&to=" + param.To
	queryStr += "&appid=" + param.App
	queryStr += "&q=" + url.QueryEscape(param.Query)
	queryStr += "&salt=" + salt
	queryStr += "&sign=" + param.buildMD5(salt)

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
