package core

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/ctool/stream"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strings"
	"time"
)

const (
	PoolSize = 5
	Prefix   = "proxy"
)

var (
	Conn           *redis.Client
	Leveldb        *leveldb.DB
	RequestList    = "" // ZSet redis key (member: 03-16 18:27:45.653 80b85e3c653, score: request nanoTime), leveldb key (80b85e3c653)
	RequestUrlList = "" // Hash: id <-> url
	listFmt        = "%s:%s:request-list"
	urlListFmt     = "%s:%s:request-url-list"
)

type (
	// Message storage in leveldb
	Message struct {
		Header http.Header `json:"header"`
		Body   []byte      `json:"body"`
	}
	// MessageVO use in rest api
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
)

func InitConnection() {
	logger.Info("loading leveldb", dbDirPath)
	newDB, err := leveldb.OpenFile(dbDirPath, nil)
	if err != nil {
		if strings.Contains(err.Error(), "resource temporarily unavailable") {
			panic("其他进程正在占用LevelDB数据库")
		}
		logger.Painc(err)
	}
	Leveldb = newDB

	InitRedisConn(ProxyConfVar.Redis)
	//if !IsValidConnection(Conn) {
	//	os.Exit(1)
	//}
	//go func() {
	//	for {
	//		time.Sleep(time.Second * 17)
	//		if !IsValidConnection(Conn) {
	//			os.Exit(1)
	//		}
	//	}
	//}()
}

func InitRedisConn(redisConf *RedisConf) {
	var opt redis.Options
	if redisConf != nil {
		if redisConf.PoolSize < 1 {
			redisConf.PoolSize = PoolSize
		}
		opt = redis.Options{Addr: redisConf.Addr, Password: redisConf.Password, DB: redisConf.DB, PoolSize: redisConf.PoolSize}

		Conn = redis.NewClient(&opt)
		if !IsValidConnection(Conn) {
			Conn = nil
		} else {
			logger.Info("Resume connect redis")
		}
	}
}

func IsValidConnection(client *redis.Client) bool {
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
	if Leveldb != nil {
		err := Leveldb.Close()
		if err != nil {
			logger.Error("close leveldb error", err)
		}
	}
}

// TrySaveLog 尝试保存，忽略静态资源及无类型标记的接口
func TrySaveLog(reqLog *ReqLog[Message], res *http.Response) {
	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		return
	}

	staticType := stream.Just(DirectContentType...).AnyMatch(func(item any) bool {
		return strings.Contains(contentType, item.(string))
	})
	if staticType {
		return
	}

	jsonType := strings.Contains(contentType, "application/json")
	if TrackAllType || jsonType && reqLog != nil {
		// TODO https://cloud.tencent.com/developer/article/1532122 优化buffer
		bodyBts, body := ctool.CopyStream(res.Body)
		res.Body = body
		resMes := Message{Header: res.Header, Body: bodyBts}
		reqLog.Response = resMes
		encoding := res.Header.Get("Content-Encoding")

		// 注意协程内不可以读写 res的Header， goproxy库里在请求结束后会将status写入header
		// github.com/kuangcp/goproxy@v1.3.5/proxy.go:353
		Go(func() {
			FillReqLogResponse(reqLog, res)
			HandleCompressed(&reqLog.Response, encoding)
			SaveReqLog(reqLog)
		})
	}
}

func IsJsonResponse(header http.Header) bool {
	contentType := header.Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

func SaveReqLog(log *ReqLog[Message]) {
	// redis cache
	Conn.ZAdd(RequestList, redis.Z{Member: log.CacheId, Score: float64(log.ReqTime.UnixNano())})
	Conn.HSet(RequestUrlList, log.Id, log.Url)

	Leveldb.Put([]byte(log.Id), ctool.ToJSONBuffer(log).Bytes(), nil)
}

func RemoveReqMember(member any) {
	Conn.ZRem(RequestList, member)
}
func RemoveReqUrlKey(key string) {
	Conn.HDel(RequestUrlList, key)
}

func ConvertLog(v *ReqLog[Message]) *ReqLog[MessageVO] {
	if v == nil {
		return nil
	}
	reqLog := ctool.CopyObj[*ReqLog[Message], ReqLog[MessageVO]](v)

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

func ConvertToDbKey(key string) string {
	return strings.Split(key, "  ")[1]
}

func QueryLogDetail(keyList []string) []*ReqLog[Message] {
	var list []*ReqLog[Message]
	for i := range keyList {
		key := keyList[i]
		dbKey := ConvertToDbKey(key)
		l := GetDetailByKey(dbKey)
		if l != nil {
			list = append(list, l)
		} else {
			logger.Warn("Delete not in leveldb key: ", key)
			RemoveReqMember(key)
			RemoveReqUrlKey(dbKey)
		}
	}
	return list
}

// MatchDetailByKeyAndKwd
// key: redis key
// kwd: 搜索关键字 (url 以及header等所有字符串)
func MatchDetailByKeyAndKwd(key, kwd string) *ReqLog[Message] {
	value, err := Leveldb.Get([]byte(key), nil)
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

func GetDetailByKey(key string) *ReqLog[Message] {
	value, err := Leveldb.Get([]byte(key), nil)
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
