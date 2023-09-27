package core

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
	PoolSize = 5
	Prefix   = "proxy"
)

var (
	Conn           *redis.Client
	db             *leveldb.DB
	RequestList    = "" // ZSet redis key (member: 03-16 18:27:45.653 80b85e3c653, score: nanoTime), leveldb key (80b85e3c653)
	RequestUrlList = "" // Hash key: id, value: url
	listFmt        = "%s:%s:request-list"
	urlListFmt     = "%s:%s:request-url-list"
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
		CacheId     string    `json:"cacheId"`
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

func Success[T any](data T) ResultVO[T] {
	return ResultVO[T]{
		Code: 0,
		Data: data,
	}
}
func InitConnection() {
	newDB, err := leveldb.OpenFile(dbDirPath, nil)
	if err != nil {
		if strings.Contains(err.Error(), "resource temporarily unavailable") {
			panic("其他进程正在占用LevelDB数据库")
		}
		logger.Painc(err)
	}
	db = newDB

	var opt redis.Options
	redisConf := ProxyConfVar.Redis
	if redisConf != nil {
		poolSize := PoolSize
		if redisConf.PoolSize != 0 {
			poolSize = redisConf.PoolSize
		}
		opt = redis.Options{Addr: redisConf.Addr, Password: redisConf.Password, DB: redisConf.DB, PoolSize: poolSize}
	} else {
		opt = redis.Options{Addr: "192.168.9.155" + ":6667", Password: "", DB: 1, PoolSize: PoolSize}
	}

	Conn = redis.NewClient(&opt)
	if !isValidConnection(Conn) {
		os.Exit(1)
	}
	go func() {
		for {
			time.Sleep(time.Second * 17)
			if !isValidConnection(Conn) {
				os.Exit(1)
			}
		}
	}()
}

func isValidConnection(client *redis.Client) bool {
	_, err := client.Ping().Result()
	if err != nil {
		logger.Error("ping redis failed:", client.Options(), err)
		return false
	}
	return true
}

func CloseConnection() {
	if Conn == nil {
		return
	}
	err := Conn.Close()
	if err != nil {
		logger.Error("close redis Conn error: ", err)
	}
	if db != nil {
		err := db.Close()
		if err != nil {
			logger.Error("close leveldb error", err)
		}
	}
}
func TrySaveLog(reqLog *ReqLog[Message], res *http.Response) {
	if !TrackAllType && !IsJsonResponse(res.Header) {
		return
	}

	FillReqLogResponse(reqLog, res)
	SaveReqLog(reqLog)
}

func IsJsonResponse(header http.Header) bool {
	contentType := header.Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

func SaveReqLog(log *ReqLog[Message]) {
	if log == nil {
		return
	}

	// redis cache
	Conn.ZAdd(RequestList, redis.Z{Member: log.CacheId, Score: float64(log.ReqTime.UnixNano())})
	Conn.HSet(RequestUrlList, log.Id, log.Url)

	db.Put([]byte(log.Id), toJSONBuffer(log).Bytes(), nil)
}

func convertLog(v *ReqLog[Message]) *ReqLog[MessageVO] {
	if v == nil {
		return nil
	}
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
		dbKey := convertToDbKey(key)
		l := getDetailByKey(dbKey)
		if l != nil {
			list = append(list, l)
		} else {
			logger.Warn("Delete not in leveldb key: ", key)
			Conn.ZRem(RequestList, key)
			Conn.HDel(RequestUrlList, dbKey)
		}
	}
	return list
}

// key: redis key
// kwd 关键字
// url 以及header等所有字符串
func matchDetailByKeyAndKwd(key, kwd string) *ReqLog[Message] {
	value, err := db.Get([]byte(key), nil)
	if err != nil {
		//logger.Error("key:["+key+"] GET ERROR:", err)
		return nil
	}

	tr := string(value)
	var l ReqLog[Message]
	err = json.Unmarshal(value, &l)

	if kwd != "" &&
		!strings.Contains(tr, kwd) &&
		!strings.Contains(string(l.Request.Body), kwd) &&
		!strings.Contains(string(l.Response.Body), kwd) {
		return nil
	}

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
