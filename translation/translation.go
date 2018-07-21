package main

import(
	"os"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var printf = fmt.Printf
var println = fmt.Println

func handlerArgs(verb string, param []string){
	switch verb {
		case "-h":
			var format string = "\033[0;32m %-5v \033[0;33m%-10v \033[0m%v \n"
			printf(format, "-h", "", "帮助")
			printf(format, "-i", "ze/ez", "使用爱词霸翻译 zh-en/en-zh")
			os.Exit(0)
		case "y":
			// 单词或者一句话 "how are you "
			sendYoudaoEntoZh(param[0])
			break
	}
}
func sendGet(url string)[]byte{
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

func readConfig() string{
	b, err := ioutil.ReadFile("/home/kcp/.config/kuangcp/baiduTrans/youdao.ini")
    if err != nil {
        fmt.Print(err)
    }
	return string(b)
}
// TODO 将本地化数据 放在用户目录下
// TODO 同样的参数做正则检验
// https://ai.youdao.com/docs/doc-trans-api.s#p01
func sendYoudaoEntoZh(word string){
	var url = readConfig()
	println("_",url,"_")
	body := sendGet(url+word)
	type web struct{
		Key string
		Value []string
	}
	type content struct{
		Translation []string
		Query string
		Web []web
	}
	var result content
	json.Unmarshal(body, &result)
	printf("原句: \033[0;32m %v\033[0m\n", result.Query)
	printf("翻译: \033[0;32m %v\033[0m\n", result.Translation)
	for i, web := range result.Web{
		printf("拓展%v:\033[0;32m %v\033[0m -> %v\n", i+1, web.Key, web.Value)
	}
}

func main(){
	var argLen int = len(os.Args)
	if argLen == 2 {
		handlerArgs(os.Args[1], nil)
	}
	if argLen > 2 {
		handlerArgs(os.Args[1], os.Args[2:])
	}
}
