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

func InitConnection(option redis.Options) {
	option.PoolSize = poolSize
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

func CloseRedisConnectionThenExit() {
	CloseConnection()
	os.Exit(1)
}
