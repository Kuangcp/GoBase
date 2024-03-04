package playground

import (
	"fmt"
	"math/rand"
	"testing"
)

type (
	Family struct {
		child []bool
	}
)

func randBool() bool {
	return rand.Int()%2 == 0
}

// 重男轻女 只要生到男孩就不生了
func TestWantBoy(t *testing.T) {
	cnt := 10000_0000
	total := 0
	for i := 0; i < cnt; i++ {
		var cs int

		for !randBool() {
			cs++
			total++
		}
		cs++
		total++
	}
	girl := total - cnt
	fmt.Printf("total:%v boy:%v girl:%v rate:%v", total, girl, cnt, float64(cnt)/float64(girl))
}

// 儿女双全
func TestBoyAndGirl(t *testing.T) {
	cnt := 10000_0000
	total := 0
	for i := 0; i < cnt; i++ {
		var cs int

		for !randBool() {
			cs++
			total++
		}
		cs++
		total++
	}
	girl := total - cnt
	fmt.Printf("total:%v boy:%v girl:%v rate:%v", total, girl, cnt, float64(cnt)/float64(girl))
}
