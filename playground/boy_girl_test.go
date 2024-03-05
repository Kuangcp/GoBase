package playground

import (
	"fmt"
	"math/rand"
	"testing"
)

const (
	evil_cnt  = 3
	evil_rate = 19
)

type (
	Family struct {
		boy      int
		girl     int
		girlDead int
	}
)

func (f *Family) oneBoyGirl() bool {
	return f.boy > 0 && f.girl > 0
}
func (f *Family) oneBoy() bool {
	return f.boy > 0
}
func (f *Family) String() string {
	return fmt.Sprintf("boy:%v girl:%v dead:%v", f.boy, f.girl, f.girlDead)
}
func evil() bool {
	return rand.Int()%100 < evil_rate
}
func bornBoy() bool {
	return rand.Int()%2 == 0
}

// 重男轻女 只要生到男孩就不生了 性别比 1:1
func TestWantBoy(t *testing.T) {
	cnt := 10000_0000
	total := 0
	for i := 0; i < cnt; i++ {
		var cs int

		for !bornBoy() {
			cs++
			total++
		}
		cs++
		total++
	}
	girl := total - cnt
	fmt.Printf("total:%v boy:%v girl:%v rate:%v", total, cnt, girl, float64(cnt)/float64(girl))
}

// 儿女双全 性别比 1:1
func TestBoyAndGirl(t *testing.T) {
	cnt := 10000_0000
	var fs = make([]*Family, cnt)
	for i := 0; i < cnt; i++ {
		f := &Family{boy: 0, girl: 0}
		for !f.oneBoyGirl() {
			if bornBoy() {
				f.boy++
			} else {
				f.girl++
			}
		}
		fs[i] = f
	}

	boy, girl := 0, 0
	for _, f := range fs {
		boy += f.boy
		girl += f.girl
	}
	fmt.Printf("total:%v boy:%v girl:%v rate:%v", boy+girl, boy, girl, float64(boy)/float64(girl))
}

// 遗弃或堕胎 女婴
// 所有家庭 5% 概率遗弃女婴 得到性别比 105:100
func TestBoyAndAbandonedGirl(t *testing.T) {
	cnt := 20000_0000
	var fs = make([]*Family, cnt)
	for i := 0; i < cnt; i++ {
		f := &Family{boy: 0, girl: 0}
		for !f.oneBoy() {
			if bornBoy() {
				f.boy++
			} else {
				if evil() && f.girlDead < evil_cnt {
					f.girlDead++
				} else {
					f.girl++
				}
			}
		}
		fs[i] = f
	}

	maxChild := 0
	var maxF *Family
	maxBorn := 0
	var maxBornF *Family

	boy, girl := 0, 0
	girlDead := 0
	for _, f := range fs {
		boy += f.boy
		girl += f.girl
		girlDead += f.girlDead

		if f.boy+f.girl > maxChild {
			maxF = f
			maxChild = f.boy + f.girl
		}
		if f.boy+f.girl+f.girlDead > maxBorn {
			maxBornF = f
			maxBorn = f.boy + f.girl + f.girlDead
		}
	}
	fmt.Printf("total:%v boy:%v girl:%v %v rate:%v \n", boy+girl, boy, girl,
		girlDead, float64(boy)/float64(girl)*100)
	fmt.Println(maxF)
	fmt.Println(maxBornF)
}
