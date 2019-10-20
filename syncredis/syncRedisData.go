package main

import (
	"log"

	"github.com/go-redis/redis"
)

var green = "\033[0;32m"
var yellow = "\033[0;33m"
var purple = "\033[0;35m"
var end = "\033[0m"

func logInfo(msg string, v ...interface{}) {
	log.Println(green+msg, v, end)
}

func logWarn(msg string, v ...interface{}) {
	log.Println(yellow+msg, v, end)
}

func initConnection() (*redis.Client, *redis.Client) {
	origin := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "",
		DB:       4,
	})

	target := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "",
		DB:       5,
	})

	_, err := origin.Ping().Result()
	if err != nil {
		log.Fatal("origin can not connetction ", err)
	}

	_, err = target.Ping().Result()
	if err != nil {
		log.Fatal("target can not connetction ", err)
	}
	return origin, target
}

func main() {
	origin, target := initConnection()

	logInfo("start sync")
	result, _ := origin.Keys("*").Result()
	logInfo("total key: ", result)

	for i := range result {
		key := result[i]
		keyType, _ := origin.Type(key).Result()
		logWarn(key, keyType)
		switch keyType {
		case STRING:
			val, _ := origin.Get(key).Result()
			logInfo("value: ", val)
			target.Set(key, val, -1)
		case LIST:
			val := origin.LRange(key, 0, -1)
			logInfo("value", val.Val())
			target.LPush(key, val.Val())
		case SET:
			val, _ := origin.SMembers(key).Result()
			logInfo("value", val)
			target.SAdd(key, val)
		case ZSET:
			val, _ := origin.ZRangeWithScores(key, 0, -1).Result()
			logInfo("value", val)
			// []string -> ...string
			target.ZAdd(key, val...)
		case HASH:
			// TODO complete
			// var val map[string]interface{}
			val, _ := origin.HGetAll(key).Result()
			logInfo("value", val)
			// target.HMSet(key, val)
		}
	}
}
