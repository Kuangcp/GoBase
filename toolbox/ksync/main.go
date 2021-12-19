package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"io/ioutil"
	"log"
	"net/http"
)

//go:embed sync.png
var iconImg string

var serverList []string
var port int
var version bool

func init() {
	flag.IntVar(&port, "p", 8000, "")
	flag.BoolVar(&version, "v", false, "")
}

func checkModFile() {
	dir, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, info := range dir {
		if info.IsDir() {
			fmt.Println("dir", info)
			continue
		}
		fmt.Println(info.Name(), info.ModTime())
	}
	fmt.Println()
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.0")
		return
	}

	//go func() {
	//	ticker := time.NewTicker(time.Second * 3)
	//	for range ticker.C {
	//		checkModFile()
	//	}
	//}()

	go func() {
		http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
			client := request.Header.Get("self")
			serverList = append(serverList, client)
			writer.Write([]byte("OK"))
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatal("error: ", err)
		}
	}()

	systray.Run(OnReady, OnExit)
}
