package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

func TestParseUrl(t *testing.T) {
	parse, err := url.Parse("http://127.0.0.1:19011/$1")
	fmt.Println(parse.Host, parse.Path, err)

	conf := ProxyConf{Name: "sss", Use: true, Routers: []string{"http://192.168.16.91:32149", "http://127.0.0.1:19011"}}
	marshal, _ := json.Marshal(conf)
	println(string(marshal))
}
