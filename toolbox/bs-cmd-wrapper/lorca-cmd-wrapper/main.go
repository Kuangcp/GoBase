package main

import (
	"flag"
	"log"
	"strings"

	"github.com/zserge/lorca"
)

var (
	url    string
	title  string
	width  int
	height int
)

func init() {
	flag.StringVar(&url, "url", "http://localhost", "")
	flag.StringVar(&title, "title", "", "")
	flag.IntVar(&width, "width", 1024, "")
	flag.IntVar(&height, "height", 768, "")
}

func main() {
	flag.Parse()

	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	ui, err := lorca.New(url, "", width, height)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	<-ui.Done()
}
