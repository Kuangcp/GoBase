package main

import (
	"flag"
	"github.com/webview/webview"
	"github.com/zserge/lorca"
	"log"
)

const (
	width  = 998
	height = 795
	url    = "https://wx2.qq.com/"
)

var (
	useWebview = false
)

func init() {
	flag.BoolVar(&useWebview, "w", false, "use webview not lorca")
}
func main() {
	flag.Parse()
	if useWebview {
		w := webview.New(false)
		defer w.Destroy()
		w.SetTitle("Wechat")
		w.SetSize(width, height, webview.HintNone)
		w.Navigate(url)
		w.Run()
	} else {
		// Create UI with basic HTML passed via data URI
		ui, err := lorca.New(url, "", width, height)
		if err != nil {
			log.Fatal(err)
		}
		defer ui.Close()
		// Wait until UI window is closed
		<-ui.Done()
	}
}
