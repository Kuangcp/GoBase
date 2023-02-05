package store

import (
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
)

const (
	poolSize = 5
)

var connection *redis.Client

func GetConnection() *redis.Client {
	return connection
}

func InitConnection(option redis.Options, verifyExit bool) {
	option.PoolSize = poolSize
	connection = redis.NewClient(&option)

	assertFunc := func(client *redis.Client) {
		if !isValidConnection(client) {
			if verifyExit {
				os.Exit(1)
			} else {
				logger.Warn("redis connect failed")
			}
		}
	}

	assertFunc(connection)
	go func() {
		for {
			time.Sleep(time.Second * 23)
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
