package stream

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestCollectList(t *testing.T) {
	result := ToList[User](
		Just(1, 8, 10, 11, 20, 21).
			Map(func(item any) any {
				v := item.(int)
				return User{
					id:     v,
					name:   fmt.Sprint(v),
					areaId: v / 3,
				}
			}),
	)
	fmt.Println(result)
	for _, u := range result {
		fmt.Println(u.id, u.areaId)
	}
}

func TestCollectToMap(t *testing.T) {
	userMap := ToMap[string, User](
		Just(1, 8, 10, 11, 20, 21).
			Map(func(item any) any {
				v := item.(int)
				return User{
					id:     v,
					name:   fmt.Sprint(v),
					areaId: v / 3,
				}
			}),
		func(a any) string {
			u := a.(User)
			return u.name
		}, Self[User],
	)
	for k, v := range userMap {
		fmt.Println(k, v)
	}
}

func TestCollectToJoin(t *testing.T) {
	join := ToJoin(Just(1, 2, 43).Map(ToString))
	fmt.Println("Join result:", join)

	ss := Just(1, 2, 43).Map(ToString)
	fmt.Println("Joins result:", ToJoins(ss, "|"))
}

func TestFork(t *testing.T) {
	list := Just(1, 4, 2, 9, 1, 1, 6).Map(ToString)
	ss := list.ForkN(10)
	fmt.Println("origin:", ToJoin(ss[0]))
	fmt.Println("fork:", ToJoins(ss[1], ","))
	fmt.Println("fork:", ToJoins(ss[2], "|"))
	fmt.Println("fork:", ToJoins(ss[3], "-"))
}

// Join 场景简单，导致了并发协程多没有提效反而更多情况是降低效率，而且结果也是乱序
func TestJoinPerformance(t *testing.T) {
	start := time.Now().UnixMicro()
	_ = ToJoin(JustN(10000).Map(ToString))
	fmt.Println(time.Now().UnixMicro()-start, "us")

	start = time.Now().UnixMicro()
	_ = ToJoin(JustN(10000).Map(ToString, WithWorkers(100)))
	fmt.Println(time.Now().UnixMicro()-start, "us")
}

func TestMultiGet(t *testing.T) {
	start := time.Now().UnixMicro()
	_ = ToJoin(JustN(70).Map(func(item any) any {
		http.Get("https://jd.com")
		return "xx"
	}, WithWorkers(10)))
	fmt.Println(time.Now().UnixMicro()-start, "us")

}
