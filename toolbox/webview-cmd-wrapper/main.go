package main

import (
	"flag"
	"github.com/webview/webview"
)

var (
	url    string
	title  string
	width  int
	height int
	debug  bool
)

func init() {
	flag.StringVar(&url, "url", "http://localhost", "")
	flag.StringVar(&title, "title", "", "")
	flag.IntVar(&width, "width", 1024, "")
	flag.IntVar(&height, "height", 768, "")
	flag.BoolVar(&debug, "d", false, "")
}

func main() {
	flag.Parse()

	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(width, height, webview.HintNone)
	w.Navigate(url)
	w.Run()
}
