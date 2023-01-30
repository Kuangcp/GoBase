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
	PoolSize = 5
	Prefix   = "dev-proxy:"
	TotalReq = Prefix + "total-req"
)

var (
	connection *redis.Client
	db         *leveldb.DB
)

type (
	Message struct {
		Header http.Header `json:"header"`
		Body   string      `json:"body"`
	}
	ReqLog struct {
		Id          string    `json:"id"`
		Url         string    `json:"url"`
		ReqTime     time.Time `json:"reqTime"`
		ResTime     time.Time `json:"resTime"`
		ElapsedTime string    `json:"useTime"`
		Request     Message   `json:"request"`
		Response    Message   `json:"response"`
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

	option.PoolSize = PoolSize
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

func saveRequest(log *ReqLog) {
	if log == nil {
		return
	}
	now := time.Now()
	key := now.Format("01-02 15:04:05.000") + " " + log.Id[0:6]
	db.Put([]byte(key), toJSONBuffer(log).Bytes(), nil)
	connection.ZAdd(TotalReq, redis.Z{Member: key, Score: float64(now.UnixNano())})
}

// page start with 1
func pageQueryReqLog(page, size string) *PageVO[ReqLog] {
	pageI, _ := strconv.Atoi(page)
	sizeI, _ := strconv.Atoi(size)
	if sizeI <= 0 || pageI < 0 {
		return nil
	}

	result, err := connection.ZRange(TotalReq, int64((pageI-1)*sizeI), int64(pageI*sizeI)-1).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}

	pageResult := PageVO[ReqLog]{}
	pageResult.Data = queryLogDetail(result)

	i, err := connection.ZCard(TotalReq).Result()
	if err == nil {
		pageResult.Total = int(i)
		pageResult.Page = int(i) / sizeI
		if pageResult.Page*sizeI < pageResult.Total {
			pageResult.Page += 1
		}
	}

	return &pageResult
}

func queryLogDetail(result []string) []ReqLog {
	var list []ReqLog
	for i := range result {
		key := result[i]

		value, err := db.Get([]byte(key), nil)
		if err != nil {
			logger.Error("key:["+key+"] GET ERROR:", err)
			continue
		}

		var l ReqLog
		err = json.Unmarshal(value, &l)
		if err != nil {
			logger.Error("key:["+key+"] GET ERROR:", err)
			continue
		}
		list = append(list, l)
	}
	return list
}
