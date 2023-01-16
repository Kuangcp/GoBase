package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	poolSize  = 5
	PREFIX    = "dev-proxy:"
	TOTAL_REQ = PREFIX + "total-req"
)

var (
	connection *redis.Client
	db         *leveldb.DB
)

type (
	ReqLog struct {
		Id      string      `json:"id"`
		Url     string      `json:"url"`
		Header  http.Header `json:"header"`
		Body    string      `json:"body"`
		Time    time.Time   `json:"time"`
		ResTime time.Time   `json:"resTime"`
	}
	ResultVO[T any] struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data T      `json:"data"`
	}
	PageVO[T any] struct {
		Total int `json:"total"`
		Page  int `json:"page"`
		Data  []T `json:"data"`
	}
)

func InitConnection() {
	newDB, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		logger.Painc(err)
	}
	db = newDB

	option := redis.Options{Addr: "127.0.0.1" + ":6667", Password: "", DB: 1}

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
	db.Put([]byte(key), marshal, nil)
	connection.ZAdd(TOTAL_REQ, redis.Z{Member: key, Score: float64(now.UnixNano())})
}

// page start with 1
func pageQueryReqLog(page, size string) *PageVO[ReqLog] {
	pageI, _ := strconv.Atoi(page)
	sizeI, _ := strconv.Atoi(size)
	result, err := connection.ZRange(TOTAL_REQ, int64((pageI-1)*sizeI), int64(pageI*sizeI)).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}
	pageResult := PageVO[ReqLog]{}
	var list []ReqLog
	for i := range result {
		key := result[i]

		value, err := db.Get([]byte(key), nil)
		if err != nil {
			logger.Error(key, err)
			continue
		}

		var l ReqLog
		err = json.Unmarshal(value, &l)
		if err != nil {
			logger.Error(key, err)
			continue
		}
		list = append(list, l)
	}
	pageResult.Data = list
	i, err := connection.ZCard(TOTAL_REQ).Result()
	if err == nil {
		pageResult.Total = int(i)
		pageResult.Page = int(i) / sizeI
		if pageResult.Page*sizeI < pageResult.Total {
			pageResult.Page += 1
		}
	}

	return &pageResult
}
