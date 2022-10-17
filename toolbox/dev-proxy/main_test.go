package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestParseUrl(t *testing.T) {
	parse, err := url.Parse("http://127.0.0.1:19011/$1")
	fmt.Println(parse.Host, parse.Path, err)
}
