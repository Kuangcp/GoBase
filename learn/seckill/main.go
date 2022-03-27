package main

import (
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"net/http"
	"time"
)

const (
	goodsKey    = "sec_goods_key:"
	goodsRegKey = "sec_goods_key_reg"
	amount      = 3000
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
	for range time.NewTicker(time.Millisecond * 1000).C {
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

		successCount, err := connection.HLen(getHisKey(latestSku)).Result()
		if err != nil {
			logger.Error(latestSku, err)
			continue
		}

		if successCount < amount {
			logger.Warn(skuCountKey, "reset to", successCount)
			connection.Set(skuCountKey, successCount, 0)
		}

		// validate
		cancelCount, err := connection.SCard(getCancelKey(latestSku)).Result()
		if err != nil {
			continue
		}

		failedCount, err := connection.SCard(getFailedKey(latestSku)).Result()
		if err != nil {
			continue
		}

		logger.Debug(failedCount + cancelCount + successCount)
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
		logger.Info("reg", skuId)
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
		_, err = connection.IncrBy(skuKey, -1).Result()
		if err != nil {
			return
		}
		connection.HDel(hisKey, param.UserId)
		connection.SAdd(getCancelKey(param.SkuId), param.UserId)
		//logger.Info("cancel:", i)
	})

	http.HandleFunc("/buy", func(writer http.ResponseWriter, request *http.Request) {
		param := parseParam(request)
		if param == nil {
			writer.Write([]byte("error"))
			return
		}

		// 进入先减库存
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
		// 已购买情况，恢复库存
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
