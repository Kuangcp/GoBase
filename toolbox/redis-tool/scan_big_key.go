package main

import (
	"context"
	"fmt"
	"github.com/kuangcp/logger"
	"github.com/redis/go-redis/v9"
	"sort"
)

type Key struct {
	size int64
	name string
}

func (k *Key) String() string {
	return fmt.Sprintf("%v => %v %vKib ", k.name, k.size, k.size/1024)
}

func queryKeyDetail(originOpt *redis.Options) {
	origin := redis.NewClient(originOpt)

	size, err := origin.MemoryUsage(context.Background(), queryKey, 0).Result()
	if err != nil {
		logger.Error(err)
		return
	}
	// TODO 支持多数据类型的统计信息
	fmt.Printf("Key:\t%v\nBytes:\t%v\t%vKib\t%vMib\n",
		queryKey, size, size>>10, size>>20)
}

func scanBigKey(originOpt *redis.Options) {
	var result []*Key
	origin := redis.NewClient(originOpt)
	var cursor uint64 = 0
	var batch int64 = 50
	counter := 0
	total := 5000
	for {
		keys, cursors, err := origin.Scan(context.Background(), cursor, "*", batch).Result()
		if err != nil {
			logger.Error(err)
			break
		}
		counter += len(keys)
		if counter > total {
			break
		}
		logger.Info("scan progress:", counter, "cursor:", cursor)
		for _, key := range keys {
			btSize, err := origin.MemoryUsage(context.Background(), key, 0).Result()
			if err != nil {
				logger.Error(err)
				continue
			}
			result = append(result, &Key{size: btSize, name: key})
		}
		if cursors == 0 {
			break
		}
		cursor = cursors
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].size > result[j].size
	})

	var totalBt int64 = 0
	for _, k := range result {
		fmt.Println(k.String())
		totalBt += k.size
	}

	logger.Info("total:", counter, "size:",
		fmt.Sprintf("%vbytes %vKib %vMib", totalBt, totalBt/1024, totalBt/1024/1024))
}
