package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var printf = fmt.Printf
var println = fmt.Println

func handlerArgs(verb string) {
	sendYoudaoEnToZh(verb)
}
func sendGet(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	return body
}

func readConfig() string {
	b, err := ioutil.ReadFile("/home/kcp/.config/kuangcp/baiduTrans/youdao.ini")
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

// TODO 同样的参数做正则检验
// https://ai.youdao.com/docs/doc-trans-api.s#p01
func sendYoudaoEnToZh(word string) {
	var url = readConfig()
	body := sendGet(url + word)
	type web struct {
		Key   string
		Value []string
	}
	type content struct {
		Translation []string
		Query       string
		Web         []web
	}
	var result content
	json.Unmarshal(body, &result)
	printf("原句: \033[0;32m %v\033[0m\n", result.Query)
	printf("翻译: \033[0;32m %v\033[0m\n", result.Translation)
	for i, web := range result.Web {
		printf("拓展%v:\033[0;32m %v\033[0m -> %v\n", i+1, web.Key, web.Value)
	}
}

func main() {
	var argLen int = len(os.Args)
	if argLen < 2 {
		println("Please input word")
		return
	}
	handlerArgs(os.Args[1])
}
