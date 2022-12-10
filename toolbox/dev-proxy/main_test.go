package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestParseUrl(t *testing.T) {
	parse, err := url.Parse("http://127.0.0.1:19011/$1")
	fmt.Println(parse.Host, parse.Path, err)

	conf := ProxyConf{Name: "sss", Enable: true, Routers: []string{"http://192.168.16.91:32149", "http://127.0.0.1:19011"}}
	marshal, _ := json.Marshal(conf)
	println(string(marshal))
}

func TestMatchPrefix(t *testing.T) {
	fmt.Println(strings.HasPrefix("xxx", ""))
}
