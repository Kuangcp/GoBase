package situation

import (
	"encoding/json"
	"fmt"
	"github.com/kuangcp/gobase/pkg/sizedpool"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
	"time"

	_ "net/http/pprof"

	"github.com/google/uuid"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/gobase/pkg/stopwatch"
	"github.com/kuangcp/logger"
)

// 场景 A 和 B C，A 和 D E ,E 和 F C 三个子组 组内构成替换关系，由于可以发生传递 和无向图的子图合并联通行为类似
type (
	Parts struct {
		Parts  string `json:"parts"`
		First  string `json:"first"`
		Second string `json:"second"`
	}
)

const (
	partsSize = 40000
	poolSize  = 90000
)

var (
	partsPool []string
)

func init() {
	go func() {
		http.ListenAndServe("0.0.0.0:8897", nil)
	}()
}

func TestGenParts(t *testing.T) {
	parts := initParts()
	marshal, _ := json.Marshal(parts)
	writer, _ := ctool.NewWriter("b.json", true)
	defer writer.Close()
	writer.Write(marshal)
}

func TestMergeCodeMap(t *testing.T) {
	parts := readParts()
	mergeCodeMap(parts)
}

func TestMergeCodeMapBench(t *testing.T) {
	//time.Sleep(time.Second * 6)
	//parts := initParts()
	parts := readParts()

	group, _ := sizedpool.New(sizedpool.PoolOption{})
	for i := 0; i < 100; i++ {
		group.Run(func() {
			mergeCodeMap(parts)
		})
		//time.Sleep(time.Second * 5)
	}

	//time.Sleep(time.Second * 60)
}

func readParts() []Parts {
	file, err := ioutil.ReadFile("b.json")
	if err != nil {
		return nil
	}

	var p []Parts
	err = json.Unmarshal(file, &p)
	if err != nil {
		return nil
	}
	return p
}

func initPool() {
	for i := 0; i < poolSize; i++ {
		partsPool = append(partsPool, uuid.New().String()[:16])
	}
}

func initParts() []Parts {
	var list []Parts

	initPool()
	for i := 0; i < partsSize; i++ {
		list = append(list, Parts{
			Parts:  partsPool[rand.Intn(poolSize)],
			First:  partsPool[rand.Intn(poolSize)],
			Second: partsPool[rand.Intn(poolSize)],
		})
	}
	return list
}

func appendMap(cache map[string]*ctool.Set[string], tmp *ctool.Set[string], key string) {
	set, ok := cache[key]
	if !ok {
		cache[key] = tmp
	} else {
		set.Adds(tmp)
	}
}

// 将水平关联的配件数据转换为层次通用数据
func mergeCodeMap(parts []Parts) map[string]*ctool.Set[string] {
	logger.Info("parts:", len(parts))
	watch := stopwatch.NewWithName("relation")
	watch.Start("init")
	cache := make(map[string]*ctool.Set[string])
	for _, p := range parts {
		tmp := ctool.NewSet(p.Parts, p.First, p.Second)
		appendMap(cache, tmp, p.Parts)
		appendMap(cache, tmp, p.First)
		appendMap(cache, tmp, p.Second)
	}
	watch.Start("merge")
	result := make(map[string]*ctool.Set[string])
	handled := ctool.NewSet[string]()
	for k := range cache {
		if handled.Contains(k) {
			continue
		}

		total := ctool.NewSet[string]()
		recursiveFind(cache, total, k)
		handled.Adds(total)
		result[uuid.New().String()[24:]] = total
	}
	watch.Stop()

	i := 0
	c := 0
	for _, v := range result {
		if v.Len() > 3 {
			i++
		}
		c += v.Len()
	}

	logger.Info("配件数:", len(parts), "去重总编码数:", len(cache), "通用码块:", len(result),
		"块内总数:", c, "合并次数:", i, watch.PrettyPrint())

	//for k, v := range result {
	//	logger.Info(k, "->", v)
	//}
	return result
}

// 由于go的栈设计能容纳深度很深，只受限于内存 goroutine stack exceeds 1000000000-byte limit
// 在编码数量级和关联程度上来说 很难超出栈的最大内存限制, 尝试100w数据 达到151w递归次数
// 这是Java无法实现的
func recursiveFind(cache map[string]*ctool.Set[string], total *ctool.Set[string], code string) {
	total.Add(code)
	block := cache[code]
	if block.Len() == 0 || block == nil {
		return
	}

	block.Loop(func(i string) {
		if total.Contains(i) {
			return
		}
		recursiveFind(cache, total, i)
	})
}

func TestParse(t *testing.T) {

	milli := time.UnixMilli(585327600000)
	fmt.Println(milli)

	fmt.Println(time.UnixMilli(653414400000))
}
