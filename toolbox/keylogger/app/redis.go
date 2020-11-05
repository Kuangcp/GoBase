package app

import (
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/wonderivan/logger"
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
	return connection
}

func InitConnection(option redis.Options) {
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
