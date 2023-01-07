package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"net/http"
	"os"
	"time"
)

const (
	poolSize  = 5
	PREFIX    = "dev-proxy:"
	TOTAL_REQ = PREFIX + "total-req"
)

var (
	connection *redis.Client
)

type (
	ReqLog struct {
		Id     string      `json:"id"`
		Url    string      `json:"url"`
		Header http.Header `json:"header"`
		Body   string      `json:"body"`
	}
)

func InitConnection() {
	option := redis.Options{Addr: "127.0.0.1" + ":6667", Password: "", DB: 1}
	fmt.Println(option)

	option.PoolSize = poolSize
	connection = redis.NewClient(&option)
	if !isValidConnection(connection) {
		os.Exit(1)
	}
	go func() {
		for {
			time.Sleep(time.Second * 17)
			if !isValidConnection(connection) {
				os.Exit(1)
			}
		}
	}()
}

func isValidConnection(client *redis.Client) bool {
	_, err := client.Ping().Result()
	if err != nil {
		logger.Error("ping redis failed:", err)
		return false
	}
	return true
}

func CloseConnection() {
	if connection == nil {
		return
	}
	err := connection.Close()
	if err != nil {
		logger.Error("close redis connection error: ", err)
	}
}

func saveRequest(log ReqLog) {
	now := time.Now()
	t := now.Format("01:02:15:04_05.000")
	marshal, err := json.Marshal(log)
	if err != nil {
		logger.Error(err)
	}
	key := t + "_" + log.Id[0:5]
	connection.Set(PREFIX+key, string(marshal), -1)
	connection.ZAdd(TOTAL_REQ, redis.Z{Member: key, Score: float64(now.UnixNano())})
}

func queryRequest() {
	result, err := connection.ZRevRange(TOTAL_REQ, 0, 3).Result()
	if err != nil {
		logger.Error(err)
		return
	}
	for i := range result {
		key := result[i]
		rs, err := connection.Get(PREFIX + key).Result()
		if err != nil {
			logger.Error(key, err)
			continue
		}

		var l ReqLog
		json.Unmarshal([]byte(rs), &l)
		fmt.Println(key, l.Url)
	}
}
