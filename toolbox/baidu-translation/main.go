package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/cuibase"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type (
	queryParam struct {
		query     string
		from      string
		to        string
		appId     string
		secretKey string
	}
	resultVO struct {
		From        string     `json:"from"`
		To          string     `json:"to"`
		TransResult []transMap `json:"trans_result"`
	}
	transMap struct {
		Src string `json:"src"` // 原文
		Dst string `json:"dst"` // 译文
	}
)

const baiduFanYiApi = "https://fanyi-api.baidu.com/api/trans/vip/translate"

func (t *queryParam) buildFinalURL() string {
	salt := "9527"
	encryptor := md5.New()
	encryptor.Write([]byte(t.appId + t.query + salt + t.secretKey))
	sign := hex.EncodeToString(encryptor.Sum(nil))

	values := url.Values{
		"from":  {t.from},
		"to":    {t.to},
		"appid": {t.appId},
		"q":     {t.query},
		"salt":  {salt},
		"sign":  {sign},
	}

	return baiduFanYiApi + "?" + values.Encode()
}

var info = cuibase.HelpInfo{
	Description: "Translation between Chinese and English By Baidu API",
	Version:     "1.0.5",
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
			Handler: handleToZh,
		}, {
			Verb:    "-ze",
			Param:   "appId secretKey query",
			Comment: "Translate zh to en",
			Handler: handleToEn,
		},
	}}

func handleToZh(params []string) {
	handleTranslation(params, "en", "zh")
}

func handleToEn(params []string) {
	handleTranslation(params, "zh", "en")
}

func handleTranslation(params []string, from, to string) {
	cuibase.AssertParamCount(4, "lack of parameters")
	param := queryParam{
		query:     fmt.Sprintf("%v", params[4:]),
		from:      from,
		to:        to,
		appId:     params[2],
		secretKey: params[3],
	}
	doQueryBaidu(param)
}

func anyStrEmpty(value ...string) bool {
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

func doQueryBaidu(param queryParam) {
	if anyStrEmpty(param.query, param.from, param.to, param.appId, param.secretKey) {
		log.Fatalln(cuibase.Red.Println(" Param exist empty string"))
	}

	finalURL := param.buildFinalURL()
	resp, err := http.Get(finalURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("resp status code: %d\n", resp.StatusCode)
	content := string(bodyContent)
	if strings.Contains(content, "\"error_code\"") {
		fmt.Printf("response error: %s\n", content)
		return
	}

	var v resultVO
	err = json.Unmarshal(bodyContent, &v)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s %v %s\n", cuibase.LightGreen, strings.Trim(v.TransResult[0].Dst, "[]"), cuibase.End)
}

func main() {
	cuibase.RunActionFromInfo(info, nil)
}
