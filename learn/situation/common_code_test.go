package situation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/pkg/stopwatch"
	"github.com/kuangcp/logger"
	"math/rand"
	"testing"
	"time"
)

// 场景 A 和 B C，A 和 D E ,E 和 F C 三个子组 组内构成替换关系，由于可以发生传递 和无向图的子图合并联通行为类似
type (
	Parts struct {
		parts  string
		first  string
		second string
	}
)

const (
	partsSize = 50
	poolSize  = 500000
)

var (
	partsPool []string
)

func initPool() {
	for i := 0; i < poolSize; i++ {
		partsPool = append(partsPool, uuid.New().String())
	}
}

func initParts() []Parts {
	var list []Parts

	for i := 0; i < partsSize; i++ {
		list = append(list, Parts{
			parts:  partsPool[rand.Intn(poolSize)],
			first:  partsPool[rand.Intn(poolSize)],
			second: partsPool[rand.Intn(poolSize)],
		})
	}
	return list
}

func TestMergeCodeMap(t *testing.T) {
	time.Sleep(time.Second * 6)
	watch := stopwatch.NewWithName("merge")
	watch.Start("init pool")
	initPool()
	watch.Start("init parts")
	parts := initParts()

	watch.Start("group p:" + fmt.Sprint(len(parts)))
	cache := make(map[string]*cuibase.Set)
	for _, p := range parts {
		tmp := cuibase.NewSet(p.parts, p.first, p.second)
		cache[p.parts] = tmp
		cache[p.first] = tmp
		cache[p.second] = tmp
	}
	watch.Start("merge")
	result := make(map[string]*cuibase.Set)
	handled := cuibase.NewSet()
	for k, _ := range cache {

		if handled.Contains(k) {
			continue
		}

		total := cuibase.NewSet()

		sub(cache, total, k)
		total.Loop(func(i interface{}) {
			handled.Add(i)
		})
		result[uuid.New().String()[24:]] = total
	}
	watch.Stop()

	i := 0
	for _, v := range result {
		if v.Len() > 3 {
			i++
		}
	}
	logger.Info("size:", len(result), "merge:", i, watch.PrettyPrint())
	time.Sleep(time.Second * 60)
}

func sub(cache map[string]*cuibase.Set, total *cuibase.Set, code string) {
	total.Add(code)
	block := cache[code]
	if block.Len() == 0 || block == nil {
		return
	}

	block.Loop(func(i interface{}) {
		if total.Contains(i) {
			return
		}
		sub(cache, total, i.(string))
	})
}
