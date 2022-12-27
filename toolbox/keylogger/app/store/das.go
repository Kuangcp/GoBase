package store

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func AddKeyDetail(time time.Time, keyNs int64, keyCode uint16) {
	conn := GetConnection()
	// store us not ns
	num, err := conn.ZAdd(GetDetailKey(time),
		redis.Z{Score: float64(keyCode), Member: keyNs / 1000}).Result()
	if err != nil {
		fmt.Println("detail zadd: ", num, err)
		CloseRedisConnectionThenExit()
	}
}

func IncrRankKey(time time.Time, keyCode uint16) {
	conn := GetConnection()
	result, err := conn.ZIncr(GetRankKey(time),
		redis.Z{Score: 1, Member: keyCode}).Result()
	if err != nil {
		fmt.Println("key zincr: ", result, err)
		CloseRedisConnectionThenExit()
	}
}

func IncrTotalCount(time time.Time) {
	conn := GetConnection()

	result, err := conn.ZIncr(TotalCount,
		redis.Z{Score: 1, Member: time.Format(DateFormat)}).Result()
	if err != nil {
		fmt.Println("total zincr: ", result, err)
		CloseRedisConnectionThenExit()
	}
}

func TotalCountVal(time time.Time) int {
	conn := GetConnection()
	total := conn.ZScore(TotalCount, time.Format(DateFormat)).Val()
	return int(total)
}

func TempKPMVal(time time.Time) string {
	conn := GetConnection()

	tempValue, err := conn.Get(GetTodayTempKPMKey(time)).Result()
	if err != nil {
		tempValue = "0"
	}
	return tempValue
}

func MaxKPMVal(time time.Time) string {
	conn := GetConnection()

	tempValue, err := conn.Get(GetTodayMaxKPMKey(time)).Result()
	if err != nil {
		tempValue = "0"
	}
	return tempValue
}
