package store

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/conf"
	"github.com/kuangcp/logger"
	"time"
)

type (
	DetailVO struct {
		Code    int
		HitTime int64
	}
)

func (d *DetailVO) ToString() string {
	return fmt.Sprintf("%v,%v", d.Code, d.HitTime)
}

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

func ExportDetailToCsv(day time.Time) {
	list := QueryDetailByDay(day.Format(DateFormat))
	if list == nil {
		logger.Warn("no detail data")
		return
	}

	dayFmt := day.Format("2006-01-02")
	writer, err := ctool.NewWriter(conf.LogDir+"/"+dayFmt+"-detail.csv", true)
	if err != nil {
		logger.Error(err)
		return
	}
	defer writer.Close()
	for _, v := range list {
		writer.WriteLine(v.ToString())
	}
}
