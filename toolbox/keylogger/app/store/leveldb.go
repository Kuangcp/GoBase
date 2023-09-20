package store

import (
	"github.com/kuangcp/gobase/toolbox/keylogger/app/conf"
	"github.com/kuangcp/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
	"strings"
)

var newDB *leveldb.DB

func GetDb() *leveldb.DB {
	return newDB
}

// InitDb 目前Leveldb仅用于web查询和处理逻辑
func InitDb() {
	db, err := leveldb.OpenFile(conf.DbPath, nil)
	if err != nil {
		logger.Error(err)
		panic("leveldb init failed")
	}

	newDB = db
}

func QueryDetailByDay(day string) []DetailVO {
	key := GetDetailKeyByString(day)
	return QueryDetailByKey(key)
}

func QueryDetailByKey(key string) []DetailVO {
	conn := GetConnection()
	result, _ := conn.Exists(key).Result()
	//logger.Info(key, result)
	if result == 1 {
		result, err := conn.ZRangeWithScores(key, 0, -1).Result()
		if err != nil {
			logger.Error(key, err)
			return nil
		}

		var list []DetailVO
		for _, v := range result {

			hitTime, err := strconv.ParseInt(v.Member.(string), 0, 64)
			if err != nil {
				logger.Warn(v, err)
				continue
			}
			list = append(list, DetailVO{Code: int(v.Score), HitTime: hitTime})
		}
		return list
	} else {
		db := GetDb()
		ret, err := db.Has([]byte(key), nil)
		if err != nil {
			logger.Error(key, err)
			return nil
		}
		if !ret {
			return nil
		}
		value, err := db.Get([]byte(key), nil)
		if err != nil {
			logger.Error(key, err)
			return nil
		}
		var list []DetailVO
		rows := strings.Split(string(value), "\n")
		for _, row := range rows {
			if row == "" {
				continue
			}
			cols := strings.Split(row, ",")
			if len(cols) < 2 {
				logger.Warn(cols)
				continue
			}
			code, _ := strconv.Atoi(cols[0])
			hitTime, _ := strconv.ParseInt(cols[1], 0, 64)
			list = append(list, DetailVO{Code: code, HitTime: hitTime})
		}
		return list
	}
}
