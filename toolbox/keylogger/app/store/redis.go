package store

import (
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
)

const (
	poolSize          = 5
	checkRedisTimeout = 23
)

var connection *redis.Client

func GetConnection() *redis.Client {
	return connection
}

func InitConnection(option redis.Options, coreApp bool) {
	option.PoolSize = poolSize
	connection = redis.NewClient(&option)

	assertFunc := func(client *redis.Client) {
		if isValidConnection(client) {
			// TODO heart beat and mark color
			//if coreApp {
			//	client.Set(CoreLive, 1, time.Minute*1)
			//}
			return
		}

		if coreApp {
			os.Exit(1)
		} else {
			logger.Warn("redis connect failed")
		}
	}

	assertFunc(connection)
	go func() {
		for range time.NewTicker(time.Duration(checkRedisTimeout) * time.Second).C {
			assertFunc(connection)
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

func CloseRedisConnectionThenExit() {
	CloseConnection()
	os.Exit(1)
}
