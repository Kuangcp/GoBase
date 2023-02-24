package main

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestParseUrl(t *testing.T) {
	parse, err := url.Parse("http://127.0.0.1:19011/$1")
	fmt.Println(parse.Host, parse.Path, err)

	//conf := ProxyConf{Name: "sss", Enable: 1, Routers: []string{"http://192.168.16.91:32149", "http://127.0.0.1:19011"}}
	//marshal, _ := json.Marshal(conf)
	//println(string(marshal))
}

func TestMatchPrefix(t *testing.T) {
	compile := regexp.MustCompile("(.*)//192.168.16.90:32009/tg-easy-fetch/(.*)")
	result := compile.ReplaceAllString("http://192.168.16.90:32009/tg-easy-fetch/index.html", "$1//localhost:8081/api/$2")
	fmt.Println("replace: ", result)

	result = compile.ReplaceAllString("http://192.168.16.90:32009/tg-easy-fetch2/index.html", "$1//localhost:8081/api/$2")
	fmt.Println("replace: ", result)

	submatch := compile.FindStringSubmatch("http://192.168.16.90:32009/tg-easy-fetch/index.html")
	fmt.Println(submatch)

	fmt.Println(strings.HasPrefix("xxx", ""))
}

func TestHostName(t *testing.T) {
	name, err := os.Hostname()
	println(name)
	println(err)
}
