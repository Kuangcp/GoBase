package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/kuangcp/gobase/cuibase"
)

const BaiduApi = "https://fanyi-api.baidu.com/api/trans/vip/translate"

type (
	QueryParam struct {
		query     string
		from      string
		to        string
		app       string
		secretKey string
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

func (t *QueryParam) buildFinalURL() string {
	salt := "9527"
	encryptor := md5.New()
	encryptor.Write([]byte(t.app + t.query + salt + t.secretKey))
	sign := hex.EncodeToString(encryptor.Sum(nil))

	values := url.Values{}
	values.Add("from", t.from)
	values.Add("to", t.to)
	values.Add("appid", t.app)
	values.Add("q", t.query)
	values.Add("salt", salt)
	values.Add("sign", sign)

	return BaiduApi + "?" + values.Encode()
}

var info = cuibase.HelpInfo{
	Description: "Translation between Chinese and English By Baidu API",
	Version:     "1.0.3",
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
				cuibase.AssertParamCount(4, "lack of parameters")
				param := QueryParam{
					query:     fmt.Sprintf("%v", params[4:]),
					from:      "en",
					to:        "zh",
					app:       params[2],
					secretKey: params[3],
				}
				doQueryBaidu(param)
			},
		}, {
			Verb:    "-ze",
			Param:   "appId secretKey query",
			Comment: "Translate zh to en",
			Handler: func(params []string) {
				cuibase.AssertParamCount(4, "lack of parameters")
				param := QueryParam{
					query:     fmt.Sprintf("%v", params[4:]),
					from:      "zh",
					to:        "en",
					app:       params[2],
					secretKey: params[3],
				}
				doQueryBaidu(param)
			},
		},
	}}

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

func doQueryBaidu(param QueryParam) {
	if anyStrEmpty(param.query, param.from, param.to, param.app, param.secretKey) {
		log.Fatalln(cuibase.Red.Println(" Param exist empty "))
	}

	resp, err := http.Get(param.buildFinalURL())
	cuibase.CheckIfError(err)
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	//fmt.Printf("resp status code: %d\n", resp.StatusCode)
	content := string(bodyContent)
	if strings.Contains(content, "\"error_code\"") {
		fmt.Printf("response error: %s\n", content)
		return
	}

	var v ResultVO
	err = json.Unmarshal(bodyContent, &v)
	cuibase.CheckIfError(err)

	fmt.Printf("%s %v %s\n", cuibase.LightGreen, strings.Trim(v.TransResult[0].Dst, "[]"), cuibase.End)
}

func main() {
	cuibase.RunActionFromInfo(info, nil)
}
