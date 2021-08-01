package main

import (
	"github.com/webview/webview"
	"os"
)

func main() {
	var url string
	if len(os.Args) > 1 {
		url = os.Args[1]
	} else {
		url = "http://localhost:8066/"
	}

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("Hosts Group")
	w.SetSize(810, 700, webview.HintNone)
	w.Navigate(url)
	w.Run()
}
