package main

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"net/http"
	"sync"
	"time"
)

type (
	Chart struct {
		Id      string
		Name    string
		mutex   sync.RWMutex
		Group   []Tick
		readIdx int
	}
	Tick struct {
		T time.Time          `json:"t"`
		V map[string]float64 `json:"v"`
	}
)

var (
	cacheSize  = 100
	chartMutex sync.RWMutex
	cache      = make(map[string]*Chart)
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/r", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		name := query.Get("name")
		if name == "" {
			writer.Write(ctool.Failed[string]().JSON())
			return
		}
		chartMutex.Lock()
		chart := &Chart{Id: ctool.RandomAlpha(5), Name: name}
		cache[chart.Id] = chart
		chartMutex.Unlock()
	})

	mux.HandleFunc("/t", func(writer http.ResponseWriter, request *http.Request) {

	})
	http.ListenAndServe(fmt.Sprintf(":%v", 9065), mux)
}
