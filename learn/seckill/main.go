package main

import (
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"net/http"
	"strconv"
	"time"
)

const (
	goodsKey    = "sec_goods_key:"
	goodsRegKey = "sec_goods_key_reg"
	amount      = 500
)

type (
	BuyParam struct {
		SkuId  string
		UserId string
	}
)

func getFailedKey(skuId string) string {
	return goodsKey + skuId + ":failed"
}
func getCancelKey(skuId string) string {
	return goodsKey + skuId + ":cancel"
}
func getCountKey(skuId string) string {
	return goodsKey + skuId + ":count"
}

func getHisKey(skuId string) string {
	return goodsKey + skuId + ":his"
}

var connection *redis.Client

func init() {
	option := redis.Options{Addr: "127.0.0.1:8834", DB: 0}
	connection = redis.NewClient(&option)
}

func sync() {
	for range time.NewTicker(time.Millisecond * 200).C {
		latestSku, err := connection.LIndex(goodsRegKey, 0).Result()
		if err != nil {
			continue
		}
		skuCountKey := getCountKey(latestSku)
		i, err := connection.Get(skuCountKey).Result()
		if i == "" {
			logger.Warn("not exist")
			continue
		}
		curAmount, _ := strconv.Atoi(i)
		//if err != nil || curAmount < amount*1.2 {
		//	//logger.Warn("not full")
		//	continue
		//}
		result, err := connection.HLen(getHisKey(latestSku)).Result()
		if err != nil {
			logger.Error(latestSku, err)
			continue
		}
		if result < amount {
			logger.Warn(skuCountKey, "reset to", result)
			connection.IncrBy(skuCountKey, -1*(int64(curAmount)-result))
			//connection.Set(skuCountKey, result, 0)
		}
	}
}

func main() {
	http.HandleFunc("/reg", func(writer http.ResponseWriter, request *http.Request) {
		param := request.URL.Query()
		skuId := param.Get("skuId")
		if skuId == "" {
			return
		}
		connection.LPush(goodsRegKey, skuId)
	})

	http.HandleFunc("/cancel", func(writer http.ResponseWriter, request *http.Request) {
		param := parseParam(request)
		if param == nil {
			writer.Write([]byte("error"))
			return
		}

		hisKey := getHisKey(param.SkuId)
		result, err := connection.HExists(hisKey, param.UserId).Result()
		if err != nil {
			return
		}
		if !result {
			return
		}

		skuKey := getCountKey(param.SkuId)
		i, err := connection.IncrBy(skuKey, -1).Result()
		if err != nil {
			return
		}
		connection.HDel(hisKey, param.UserId)
		connection.SAdd(getCancelKey(param.SkuId), param.UserId)
		logger.Info("cancel:", i)
	})

	http.HandleFunc("/buy", func(writer http.ResponseWriter, request *http.Request) {
		param := parseParam(request)
		if param == nil {
			writer.Write([]byte("error"))
			return
		}

		result, err := connection.Incr(getCountKey(param.SkuId)).Result()
		if err != nil {
			logger.Error(err)
			writer.Write([]byte("error"))
			return
		}

		if result > amount {
			writer.Write([]byte("SALE_OUT"))
			connection.SAdd(getFailedKey(param.SkuId), param.UserId)
			return
		}

		userAlreadyBuy, err := connection.HExists(getHisKey(param.SkuId), param.UserId).Result()
		if err != nil {
			logger.Error(err)
			writer.Write([]byte("error"))
			return
		}
		// 恢复库存
		if userAlreadyBuy {
			_, err := connection.IncrBy(getCountKey(param.SkuId), -1).Result()
			if err != nil {
				logger.Error(err)
				writer.Write([]byte("error"))
				return
			}
			//logger.Info(result)
			writer.Write([]byte("already buy"))
			return
		}

		connection.HSet(getHisKey(param.SkuId), param.UserId, 0)
		writer.Write([]byte("OK"))
	})

	go sync()
	http.ListenAndServe(":8099", nil)
}

func parseParam(request *http.Request) *BuyParam {
	param := request.URL.Query()
	skuId := param.Get("skuId")
	userId := param.Get("userId")

	if skuId == "" || userId == "" {
		return nil
	}

	buyParam := BuyParam{
		SkuId:  skuId,
		UserId: userId,
	}
	return &buyParam
}
