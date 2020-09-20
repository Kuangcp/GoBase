package main

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/kuangcp/threadpool"
)

func logInfo(msg string, v ...interface{}) {
	log.Println(cuibase.Green, msg, v, cuibase.End)
}

func logWarn(msg string, v ...interface{}) {
	log.Println(cuibase.Yellow, msg, v, cuibase.End)
}

func Action(originO *redis.Options, targetO *redis.Options,
	action func(client *redis.Client, client2 *redis.Client)) {

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

	pool := threadpool.NewThreadPoolWithPrefix(10, 10000, "sync-")
	var latch sync.WaitGroup
	latch.Add(len(result))
	for i := range result {
		key := result[i]
		keyType, _ := origin.Type(key).Result()
		err := pool.ExecuteFunc(func(workerId string) {
			defer latch.Done()

			logWarn(workerId, key, keyType)
			switch keyType {
			case STRING:
				val, _ := origin.Get(key).Result()
				//logInfo(workerId, "value: ", val)

				if strings.HasPrefix(key, "all-") {
					atoi, _ := strconv.Atoi(val)
					finalKey := strings.ReplaceAll(strings.TrimPrefix(key, "all-"), "-", ":")
					target.ZAdd("keyboard:total", redis.Z{Member: finalKey, Score: float64(atoi)})
				} else {
					target.Set(key, val, -1)
				}
			case ZSET:
				val, _ := origin.ZRangeWithScores(key, 0, -1).Result()
				//logInfo(workerId, "value", val)

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
		})
		if err != nil {
			logWarn("", err)
		}
	}
	latch.Wait()
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
	//logInfo("total key: ", result)

	pool := threadpool.NewThreadPoolWithPrefix(10, 10000, "sync-")
	var latch sync.WaitGroup
	latch.Add(len(result))
	for i := range result {
		key := result[i]
		keyType, _ := origin.Type(key).Result()
		err := pool.ExecuteFunc(func(workerId string) {
			defer latch.Done()

			logWarn(workerId, key, keyType)
			switch keyType {
			case STRING:
				val, _ := origin.Get(key).Result()
				logInfo(workerId, "value: ", val)
				target.Set(key, val, -1)
			case LIST:
				val := origin.LRange(key, 0, -1)
				logInfo(workerId, "value", val.Val())
				target.LPush(key, val.Val())
			case SET:
				val, _ := origin.SMembers(key).Result()
				logInfo(workerId, "value", val)
				target.SAdd(key, val)
			case ZSET:
				val, _ := origin.ZRangeWithScores(key, 0, -1).Result()
				logInfo(workerId, "value", val)
				// []string -> ...string
				target.ZAdd(key, val...)
			case HASH:
				val, _ := origin.HGetAll(key).Result()
				logInfo(workerId, "value", val)
				target.HMSet(key, convert(val))
			}
		})
		if err != nil {
			logWarn("", err)
		}
	}
	latch.Wait()
}
