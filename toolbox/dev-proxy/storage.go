package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	PoolSize    = 5
	Prefix      = "proxy:"
	RequestList = Prefix + "request-list"
)

var (
	connection *redis.Client
	db         *leveldb.DB
)

type (
	// storage in leveldb
	Message struct {
		Header http.Header `json:"header"`
		Body   []byte      `json:"body"`
	}
	// use in rest api
	MessageVO struct {
		Header  http.Header `json:"header"`
		Body    any         `json:"body"`
		BodyStr *string     `json:"bodyStr,omitempty"`
	}

	ReqLog[T any] struct {
		Id          string    `json:"id"`
		Method      string    `json:"method"`
		Url         string    `json:"url"`
		Status      string    `json:"status"`
		StatusCode  int       `json:"statusCode"`
		ReqTime     time.Time `json:"reqTime"`
		ResTime     time.Time `json:"resTime"`
		ElapsedTime string    `json:"useTime"`
		Request     T         `json:"request"`
		Response    T         `json:"response"`
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

	option := redis.Options{Addr: "192.168.9.155" + ":6667", Password: "", DB: 1}

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

func convertLog(v *ReqLog[Message]) *ReqLog[MessageVO] {
	reqLog := copyObj[*ReqLog[Message], ReqLog[MessageVO]](v)

	reqLog.Request = MessageVO{Header: v.Request.Header, Body: strToAny(v.Request.Body)}
	reqLog.Response = MessageVO{Header: v.Response.Header, Body: strToAny(v.Response.Body)}

	fillDefault(v.Request, &reqLog.Request)
	fillDefault(v.Response, &reqLog.Response)

	return reqLog
}

func fillDefault(src Message, vo *MessageVO) {
	reqLen := len(src.Body)
	if vo.Body != nil || reqLen <= 0 {
		return
	}

	var str string
	if reqLen > 0 && reqLen < 100 {
		str = string(src.Body)
	} else {
		str = "请求体过大"
	}
	vo.BodyStr = &str
}

func strToAny(body []byte) any {
	if body == nil || len(body) == 0 {
		return nil
	}
	var d any
	err := json.Unmarshal(body, &d)
	if err != nil {
		//logger.Error(err)
		return nil
	}
	return d
}

func queryLogDetail(keyList []string) []*ReqLog[Message] {
	var list []*ReqLog[Message]
	for i := range keyList {
		key := keyList[i]
		l := getDetailByKey(convertToDbKey(key))
		list = append(list, l)
	}
	return list
}

func matchDetailByKeyAndKwd(key, kwd string) *ReqLog[Message] {
	value, err := db.Get([]byte(key), nil)
	if err != nil {
		//logger.Error("key:["+key+"] GET ERROR:", err)
		return nil
	}

	if kwd != "" && !strings.Contains(string(value), kwd) {
		return nil
	}

	var l ReqLog[Message]
	err = json.Unmarshal(value, &l)
	if err != nil {
		logger.Error("key:["+key+"] GET ERROR:", err, len(value))
		return nil
	}
	return &l
}

func getDetailByKey(key string) *ReqLog[Message] {
	value, err := db.Get([]byte(key), nil)
	if err != nil {
		logger.Error("key:["+key+"] GET ERROR:", err, len(value))
		return nil
	}

	var l ReqLog[Message]
	err = json.Unmarshal(value, &l)
	if err != nil {
		logger.Error("key:["+key+"] GET ERROR:", err, len(value))
		return nil
	}
	return &l
}
