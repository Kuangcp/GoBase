package main

import (
	"fmt"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var skuId string

func init() {
	_ = logger.SetLoggerConfig(&logger.LogConfig{
		Console: &logger.ConsoleLogger{
			Level:    logger.InformationalDesc,
			Colorful: true,
		},
		TimeFormat: logger.LogTimeDetailFormat,
	})
	skuId = fmt.Sprint(time.Now().UnixMilli())

	http.Get("http://localhost:8099/reg?skuId=" + skuId)
}

func cancel(id int) {
	if rand.Intn(100)%3 != 0 {
		return
	}

	time.Sleep(time.Millisecond * time.Duration((200+id)%60))
	userId := fmt.Sprint("user", id)
	rsp, err := http.Get("http://localhost:8099/cancel?skuId=" + skuId + "&userId=" + userId)
	if err != nil {
		return
	}
	rsp.Body.Close()
}

func secBuy(id int) {
	userId := fmt.Sprint("user", id)
	for {
		rsp, err := http.Get("http://localhost:8099/buy?skuId=" + skuId + "&userId=" + userId)
		if err != nil {
			continue
		}

		time.Sleep(time.Millisecond * time.Duration((id+333)%295))
		all, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			continue
		}
		rsp.Body.Close()
		rspStr := string(all)
		if rspStr == "OK" {
			//logger.Info("ok", userId)
			fmt.Print(id, ",")
			break
		}
		if rspStr == "SALE_OUT" {
			//logger.Info("sale_out", userId)
			break
		}
	}
}

func TestBuy(t *testing.T) {
	logger.Warn("start new")
	for i := 0; i < 1000; i++ {
		go secBuy(i)
		go cancel(i)
	}

	time.Sleep(time.Second * 4)
	logger.Warn("start new")
	for i := 1000; i < 2500; i++ {
		go secBuy(i)
	}
	time.Sleep(time.Second * 4)
}

func TestRand(t *testing.T) {
	for i := 0; i < 299; i++ {
		fmt.Println(rand.Intn(100))
	}
}

func TestRename(t *testing.T) {

}
