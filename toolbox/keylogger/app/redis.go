package app

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

//Prefix redis prefix
const Prefix = "keyboard:"

//LastInputEvent last use event
const LastInputEvent = Prefix + "last-event"

//TotalCount total count ZSET
const TotalCount = Prefix + "total"
const DateFormat = "2006:01:02"

//KeyMap cache key code map HASH
const KeyMap = Prefix + "key-map"

var connection *redis.Client

//GetRankKey by time
func GetRankKey(time time.Time) string {
	return GetRankKeyByString(time.Format(DateFormat))
}

//GetRankKeyByString
func GetRankKeyByString(time string) string {
	return Prefix + time + ":rank"
}

//GetDetailKey by time
func GetDetailKey(time time.Time) string {
	return GetDetailKeyByString(time.Format(DateFormat))
}

func GetDetailKeyByString(time string) string {
	return Prefix + time + ":detail"
}

func GetConnection() *redis.Client {
	if time.Now().Second()%7 == 0 {
		_, err := connection.Ping().Result()
		if err != nil {
			fmt.Println("ping redis failed:", err)
			os.Exit(1)
		}
	}
	return connection
}

func InitConnection(option redis.Options) {
	connection = redis.NewClient(&option)
	_, err := connection.Ping().Result()
	if err != nil {
		fmt.Println("ping redis failed:", err)
		os.Exit(1)
	}
}

func CloseConnection() {
	if connection == nil {
		return
	}
	err := connection.Close()
	if err != nil {
		fmt.Println("close redis connection error: ", err)
	}
}
