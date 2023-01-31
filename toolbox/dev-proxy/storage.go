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
	// storage in leveldb
	Message struct {
		Header http.Header `json:"header"`
		Body   string      `json:"body"`
	}
	// use in rest api
	MessageVO struct {
		Header  http.Header `json:"header"`
		Body    any         `json:"body"`
		BodyStr *string     `json:"bodyStr,omitempty"`
	}

	ReqLog[T any] struct {
		Id          string    `json:"id"`
		Url         string    `json:"url"`
		ReqTime     time.Time `json:"reqTime"`
		ResTime     time.Time `json:"resTime"`
		ElapsedTime string    `json:"useTime"`
		Request     T         `json:"request"`
		Response    T         `json:"response"`
		Status      string    `json:"status"`
		StatusCode  int       `json:"statusCode"`
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

func saveReqLog(log *ReqLog[Message]) {
	if log == nil {
		return
	}

	db.Put([]byte(log.Id), toJSONBuffer(log).Bytes(), nil)
}

// page start with 1
func pageQueryReqLog(page, size string) *PageVO[*ReqLog[MessageVO]] {
	pageI, _ := strconv.Atoi(page)
	sizeI, _ := strconv.Atoi(size)
	if sizeI <= 0 || pageI < 0 {
		return nil
	}

	keyList, err := connection.ZRange(TotalReq, int64((pageI-1)*sizeI), int64(pageI*sizeI)-1).Result()
	if err != nil {
		logger.Error(err)
		return nil
	}

	pageResult := PageVO[*ReqLog[MessageVO]]{}
	detail := queryLogDetail(keyList)
	pageResult.Data = convertList(detail, convertLog, nil)

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

func convertLog(v *ReqLog[Message]) *ReqLog[MessageVO] {
	reqLog := &ReqLog[MessageVO]{
		Id: v.Id, Url: v.Url, ReqTime: v.ReqTime, ResTime: v.ResTime, ElapsedTime: v.ElapsedTime,
		Status: v.Status, StatusCode: v.StatusCode,
		Request:  MessageVO{Header: v.Request.Header, Body: strToAny(v.Request.Body)},
		Response: MessageVO{Header: v.Response.Header, Body: strToAny(v.Response.Body)}}
	if reqLog.Request.Body == nil && v.Request.Body != "" {
		reqLog.Request.BodyStr = &v.Request.Body
	}
	if reqLog.Response.Body == nil && v.Response.Body != "" {
		reqLog.Response.BodyStr = &v.Response.Body
	}
	return reqLog
}

func strToAny(body string) any {
	if body == "" {
		return nil
	}
	var d any
	err := json.Unmarshal([]byte(body), &d)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return d
}

func queryLogDetail(result []string) []*ReqLog[Message] {
	var list []*ReqLog[Message]
	for i := range result {
		key := result[i]

		value, err := db.Get([]byte(key), nil)
		if err != nil {
			logger.Error("key:["+key+"] GET ERROR:", err)
			continue
		}

		var l ReqLog[Message]
		err = json.Unmarshal(value, &l)
		if err != nil {
			logger.Error("key:["+key+"] GET ERROR:", err)
			continue
		}
		list = append(list, &l)
	}
	return list
}
