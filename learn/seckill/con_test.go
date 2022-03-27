package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var skuId string
var secKillTime = 8
var delaySec = 5
var startSec, endSec int

func init() {
	_ = logger.SetLoggerConfig(&logger.LogConfig{
		Console: &logger.ConsoleLogger{
			Level:    logger.InformationalDesc,
			Colorful: true,
		},
		TimeFormat: logger.LogTimeDetailFormat,
	})
	skuId = fmt.Sprint(9999999999999 - time.Now().UnixMilli())

	http.Get("http://localhost:8099/reg?skuId=" + skuId)
	second := time.Now().Second()
	startSec = second + delaySec
	endSec = second + delaySec + secKillTime
}

func TestBuy(t *testing.T) {
	for i := 0; i < 100; i++ {
		for j := 0; j < 500; j++ {
			go secBuy(i*1000 + j)
			go cancel(i*1000+j, 7)
		}
		//time.Sleep(time.Millisecond * 500)
	}

	fmt.Println(startSec, endSec)
	time.Sleep(time.Second * time.Duration(delaySec+secKillTime+30))
}

func cancel(id, cancelRate int) {
	if rand.Intn(1000)%cancelRate != 0 {
		return
	}

	//fmt.Println(id)
	for {
		second := time.Now().Second()
		if second > endSec {
			logger.Error("timeout")
			break
		}

		if second < startSec {
			time.Sleep(time.Millisecond * 300)
			continue
		}

		time.Sleep(time.Millisecond * time.Duration((200+id)%360+200))
		userId := fmt.Sprint("user", id)
		rsp, err := http.Get("http://localhost:8099/cancel?skuId=" + skuId + "&userId=" + userId)
		if err != nil {
			return
		}
		rsp.Body.Close()
		break
	}
}

func secBuy(id int) {
	userId := fmt.Sprint("user", id)
	for {
		second := time.Now().Second()
		if second > endSec {
			logger.Error("timeout")
			break
		}
		if second < startSec {
			time.Sleep(time.Millisecond * 500)
			//logger.Info("waiting")
			continue
		}
		rsp, err := http.Get("http://localhost:8099/buy?skuId=" + skuId + "&userId=" + userId)
		if err != nil {
			continue
		}

		//time.Sleep(time.Millisecond * time.Duration((id+333)%200))
		all, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			continue
		}
		rsp.Body.Close()
		rspStr := string(all)
		if rspStr == "OK" {
			//logger.Info("ok", userId)
			//fmt.Print(id, ",")
			break
		}
		if rspStr == "SALE_OUT" {
			//logger.Info("sale_out", userId)
			break
		}
	}
}

func TestRand(t *testing.T) {
	for i := 0; i < 299; i++ {
		fmt.Println(rand.Intn(100))
	}
}

func TestSet(t *testing.T) {
	option := redis.Options{Addr: "127.0.0.1:8834", DB: 0}
	connection = redis.NewClient(&option)

	connection.Set("ab", 3, 0)
	connection.Incr("ab")
}
