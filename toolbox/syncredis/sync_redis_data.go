package main

import (
	"flag"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/sizedwaitgroup"
)

var debugFlag = false

var (
	fromAddr string
	fromPwd  string
	fromDB   int

	toAddr string
	toPwd  string
	toDB   int
)

func init() {
	flag.StringVar(&fromAddr, "f.addr", "127.0.0.1:6379", "from redis address")
	flag.StringVar(&fromPwd, "f.pwd", "", "from redis password")
	flag.IntVar(&fromDB, "f.db", 2, "from redis db")
	flag.StringVar(&toAddr, "t.addr", "127.0.0.1:6379", "to redis address")
	flag.StringVar(&toPwd, "t.pwd", "", "to redis password")
	flag.IntVar(&toDB, "t.db", 3, "to redis db")
}

func main() {
	flag.Parse()

	Action(false,
		&redis.Options{
			Addr:     fromAddr,
			Password: fromPwd,
			DB:       fromDB,
		},
		&redis.Options{
			Addr:     toAddr,
			Password: toPwd,
			DB:       toDB,
		}, SyncAllKey)
}

func logDebug(msg string, v ...interface{}) {
	if debugFlag {
		log.Println(msg, v)
	}
}

func logInfo(msg string, v ...interface{}) {
	log.Println(cuibase.Green, msg, v, cuibase.End)
}

func logWarn(msg string, v ...interface{}) {
	log.Println(cuibase.Yellow, msg, v, cuibase.End)
}

func Action(debug bool, originO *redis.Options, targetO *redis.Options,
	action func(client *redis.Client, client2 *redis.Client)) {

	debugFlag = debug
	if originO == nil || targetO == nil {
		log.Fatal("origin or target option is nil")
	}

	origin := redis.NewClient(originO)
	target := redis.NewClient(targetO)

	_, err := origin.Ping().Result()
	if err != nil {
		log.Fatal("origin can not connection ", err)
	}

	_, err = target.Ping().Result()
	if err != nil {
		log.Fatal("target can not connection ", err)
	}
	action(origin, target)
}

func SyncKeyRecord(origin *redis.Client, target *redis.Client) {
	logInfo("start sync")
	result, _ := origin.Keys("*").Result()
	//logInfo("total key: ", result)

	swg := sizedwaitgroup.New(12)

	for i := range result {
		swg.Add()
		key := result[i]
		keyType, _ := origin.Type(key).Result()
		go func() {

			logWarn(key, keyType)
			switch keyType {
			case STRING:
				val, _ := origin.Get(key).Result()
				logDebug("value: ", val)

				if strings.HasPrefix(key, "all-") {
					atoi, _ := strconv.Atoi(val)
					finalKey := strings.ReplaceAll(strings.TrimPrefix(key, "all-"), "-", ":")
					target.ZAdd("keyboard:total", redis.Z{Member: finalKey, Score: float64(atoi)})
				} else {
					target.Set(key, val, -1)
				}
			case ZSET:
				val, _ := origin.ZRangeWithScores(key, 0, -1).Result()
				logDebug("value", val)

				newKey := key
				if strings.HasPrefix(key, "detail-") {
					day := strings.TrimPrefix(key, "detail-")
					newKey = "keyboard:" + strings.ReplaceAll(day, "-", ":") + ":detail"
					var result []redis.Z
					for _, ele := range val {
						float, _ := strconv.ParseFloat(ele.Member.(string), 64)

						result = append(result, redis.Z{Member: int64(float * 1000000), Score: ele.Score})
					}
					// []string -> ...string
					target.ZAdd(newKey, result...)
				}

				if strings.HasPrefix(key, "2019") || strings.HasPrefix(key, "2020") {
					newKey = "keyboard:" + strings.ReplaceAll(key, "-", ":") + ":rank"
					// []string -> ...string
					target.ZAdd(newKey, val...)
				}
			}
		}()
	}
	swg.Wait()
}

func convert(data map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range data {
		result[k] = v
	}
	return result
}
func SyncAllKey(origin *redis.Client, target *redis.Client) {
	logInfo("start sync")
	result, _ := origin.Keys("*").Result()

	swg := sizedwaitgroup.New(12)
	logInfo("total key: ", len(result))
	for i := range result {
		swg.Add()
		key := result[i]
		keyType, _ := origin.Type(key).Result()

		go func() {
			defer swg.Done()

			logWarn(key, keyType)
			switch keyType {
			case STRING:
				val, _ := origin.Get(key).Result()
				logDebug("value: ", val)
				target.Set(key, val, -1)
			case LIST:
				val := origin.LRange(key, 0, -1)
				logDebug("value", val.Val())
				target.LPush(key, val.Val())
			case SET:
				val, _ := origin.SMembers(key).Result()
				logDebug("value", val)
				target.SAdd(key, val)
			case ZSET:
				val, _ := origin.ZRangeWithScores(key, 0, -1).Result()
				logDebug("value", val)
				// []string -> ...string
				target.ZAdd(key, val...)
			case HASH:
				val, _ := origin.HGetAll(key).Result()
				logDebug("value", val)
				target.HMSet(key, convert(val))
			}
		}()
	}
	swg.Wait()
}
