package main

import (
	"context"
	"github.com/kuangcp/gobase/pkg/ctool"
	"log"
	"strconv"
	"strings"

	"github.com/kuangcp/sizedwaitgroup"
	"github.com/redis/go-redis/v9"
)

func logDebug(msg string, v ...interface{}) {
	if debugFlag {
		log.Println(msg, v)
	}
}

func logInfo(msg string, v ...interface{}) {
	log.Println(ctool.Green, msg, v, ctool.End)
}

func logWarn(msg string, v ...interface{}) {
	log.Println(ctool.Yellow, msg, v, ctool.End)
}

func scanAllKey(origin redis.Cmdable) []string {
	cursor := origin.Scan(context.Background(), 0, "*", 1000)

	result, c, err := cursor.Result()
	if err != nil {
		return result
	}
	var totalKey = result
	for c != 0 {
		cursor := origin.Scan(context.Background(), c, "*", 1000)
		result, c, err = cursor.Result()
		if err != nil {
			return totalKey
		}
		totalKey = append(totalKey, result...)
	}
	return totalKey
}

func convert(data map[string]string) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range data {
		result[k] = v
	}
	return result
}

func Action(action func(client redis.Cmdable, client2 redis.Cmdable),
	originO *redis.Options,
	targetO *redis.Options,
	debug bool) {

	debugFlag = debug
	if originO == nil || targetO == nil {
		log.Fatal("origin or target option is nil")
	}

	origin := redis.NewClient(originO)
	target := redis.NewClient(targetO)

	_, err := origin.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("origin can not connection ", err)
	}

	_, err = target.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("target can not connection ", err)
	}
	action(origin, target)
}

func SyncKeyRecord(origin redis.Cmdable, target redis.Cmdable) {
	logInfo("start sync keylogger data")

	//result, _ := origin.Keys("*").Result()
	result := scanAllKey(origin)
	logInfo("total key: ", result)

	swg := sizedwaitgroup.New(12)

	for i := range result {
		swg.Add()
		key := result[i]
		keyType, _ := origin.Type(context.Background(), key).Result()
		go func() {

			logWarn(key, keyType)
			switch keyType {
			case STRING:
				val, _ := origin.Get(context.Background(), key).Result()
				logDebug("value: ", val)

				if strings.HasPrefix(key, "all-") {
					atoi, _ := strconv.Atoi(val)
					finalKey := strings.ReplaceAll(strings.TrimPrefix(key, "all-"), "-", ":")
					target.ZAdd(context.Background(), "keyboard:total", redis.Z{Member: finalKey, Score: float64(atoi)})
				} else {
					target.Set(context.Background(), key, val, -1)
				}
			case ZSET:
				val, _ := origin.ZRangeWithScores(context.Background(), key, 0, -1).Result()
				logDebug("value", val)

				newKey := key
				if strings.HasPrefix(key, "detail-") {
					day := strings.TrimPrefix(key, "detail-")
					newKey = "keyboard:" + strings.ReplaceAll(day, "-", ":") + ":detail"
					var result []redis.Z
					for _, ele := range val {
						float, _ := strconv.ParseFloat(ele.Member.(string), 64)

						result = append(result, redis.Z{Member: int64(float * 1000000), Score: ele.Score})
					}
					// []string -> ...string
					target.ZAdd(context.Background(), newKey, result...)
				}

				if strings.HasPrefix(key, "2019") || strings.HasPrefix(key, "2020") {
					newKey = "keyboard:" + strings.ReplaceAll(key, "-", ":") + ":rank"
					// []string -> ...string
					target.ZAdd(context.Background(), newKey, val...)
				}
			}
		}()
	}
	swg.Wait()
}

func SyncAllKey(origin redis.Cmdable, target redis.Cmdable) {
	logInfo("start sync all key")
	//result, _ := origin.Keys("*").Result()
	result := scanAllKey(origin)

	swg := sizedwaitgroup.New(12)
	logInfo("total key: ", len(result))
	for i := range result {
		swg.Add()
		key := result[i]
		keyType, _ := origin.Type(context.Background(), key).Result()

		go func() {
			defer swg.Done()

			logWarn(key, keyType)
			switch keyType {
			case STRING:

				val, _ := origin.Get(context.Background(), key).Result()
				logDebug("value: ", val)
				target.Set(context.Background(), key, val, -1)
			case LIST:
				val := origin.LRange(context.Background(), key, 0, -1)
				logDebug("value", val.Val())
				target.LPush(context.Background(), key, val.Val())
			case SET:
				val, _ := origin.SMembers(context.Background(), key).Result()
				logDebug("value", val)
				target.SAdd(context.Background(), key, val)
			case ZSET:
				val, _ := origin.ZRangeWithScores(context.Background(), key, 0, -1).Result()
				logDebug("value", val)
				// []string -> ...string
				target.ZAdd(context.Background(), key, val...)
			case HASH:
				val, _ := origin.HGetAll(context.Background(), key).Result()
				logDebug("value", val)
				target.HMSet(context.Background(), key, convert(val))
			}
		}()
	}
	swg.Wait()
}
