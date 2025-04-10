package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	host   string
	client bool
)

func main() {
	flag.StringVar(&host, "h", "localhost:80", "host:port")
	flag.BoolVar(&client, "c", false, "client mode")
	flag.Parse()

	if client {
		clientReq()
		return
	}

	listenServer()
}

func listenServer() {
	http.HandleFunc("/s", func(w http.ResponseWriter, r *http.Request) {
		req := time.Now().UnixMilli()
		val := r.URL.Query().Get("c")
		if val == "" {
			w.Write([]byte("Error"))
			return
		}
		ct, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			w.Write([]byte("Error"))
			return
		}

		w.Write([]byte(fmt.Sprint(req - ct)))
	})

	http.ListenAndServe(":8343", nil)
}

func clientReq() {
	var tt int64 = 0
	cnt := 10
	for i := 0; i < cnt; i++ {
		tt += queryDiff()
	}
	fmt.Printf("%vms %v avg: %vms\n", tt, cnt, tt/int64(cnt))
}

func queryDiff() int64 {
	req := time.Now().UnixMilli()
	resp, err := http.Get("http://" + host + "/s?c=" + fmt.Sprintf("%d", req))
	if err != nil {
		log.Fatal(err)
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	delta := string(all)
	ct, err := strconv.ParseInt(delta, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return ct
}
