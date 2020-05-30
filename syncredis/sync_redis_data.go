package main

import (
	"log"
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

func initConnection() (*redis.Client, *redis.Client) {
	origin := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "",
		DB:       3,
	})

	target := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "",
		DB:       4,
	})

	_, err := origin.Ping().Result()
	if err != nil {
		log.Fatal("origin can not connection ", err)
	}

	_, err = target.Ping().Result()
	if err != nil {
		log.Fatal("target can not connection ", err)
	}
	return origin, target
}

func main() {
	origin, target := initConnection()

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
				// TODO complete
				// var val map[string]interface{}
				val, _ := origin.HGetAll(key).Result()
				logInfo(workerId, "value", val)
				// target.HMSet(key, val)
			}
		})
		if err != nil {
			logWarn("", err)
		}
	}
	latch.Wait()
}
