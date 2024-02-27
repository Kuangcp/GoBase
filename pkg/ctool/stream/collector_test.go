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
	userStream := Just(1, 8, 10, 11, 20, 21).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   "user-" + fmt.Sprint(v),
			areaId: v / 3,
		}
	})

	userMap := ToMap[string, User](userStream, func(a any) string {
		u := a.(User)
		return u.name
	}, Self[User])
	for k, v := range userMap {
		fmt.Println(k, v)
	}

	aMap := ToMap[int, User](userStream, func(a any) int {
		u := a.(User)
		return u.areaId
	}, Self[User])
	for k, v := range aMap {
		fmt.Println("merge", k, v)
	}
}

func TestCollectToMaps(t *testing.T) {
	userStream := JustN(13).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   "user-" + fmt.Sprint(v),
			areaId: v % 3,
		}
	})

	origin, fork := userStream.Fork()
	aMap := ToMap[int, User](origin, func(a any) int {
		u := a.(User)
		return u.areaId
	}, Self[User])
	for k, v := range aMap {
		fmt.Println("over", k, v)
	}

	aMap = ToMaps[int, User](fork, func(a any) int {
		u := a.(User)
		return u.areaId
	}, Self[User], func(old User, cur User) User {
		return old
	})
	for k, v := range aMap {
		fmt.Println("merge", k, v)
	}
}

func TestCollectToJoin(t *testing.T) {
	join := ToJoin(Just(1, 2, 43).MapStr())
	fmt.Println("Join result:", join)

	ss := Just(1, 2, 43).MapStr()
	fmt.Println("Joins result:", ToJoins(ss, "|"))
}

func TestToSet(t *testing.T) {
	set := ToSet[int](Just(1, 4, 2, 2, 1))
	fmt.Println(set)

	set = ToSetFunc[int](Just(1, 4, 2, 2, 1), func(s any) int {
		return s.(int) + 10
	})
	fmt.Println(set)
}

// Join 场景简单，导致了并发协程多没有提效反而更多情况是降低效率，而且结果也是乱序
func TestJoinPerformance(t *testing.T) {
	start := time.Now().UnixMicro()
	_ = ToJoin(JustN(10000).MapStr())
	fmt.Println("sin", time.Now().UnixMicro()-start, "us")

	start = time.Now().UnixMicro()
	_ = ToJoin(JustN(10000).Map(ToString, WithWorkers(100)))
	fmt.Println("con", time.Now().UnixMicro()-start, "us")
}

func TestMultiGet(t *testing.T) {
	start := time.Now().UnixMicro()
	_ = ToJoin(JustN(70).Map(func(item any) any {
		http.Get("https://jd.com")
		return "xx"
	}, WithWorkers(10)))
	fmt.Println(time.Now().UnixMicro()-start, "us")

}

func TestToSum(t *testing.T) {
	sum := ToSum[int](JustN(4))
	fmt.Println(sum)
}

func TestToGroupBy(t *testing.T) {
	result := ToGroupBy[int, int](JustN(10), func(item int) int {
		return item % 3
	})
	fmt.Println(result)
	for k, v := range result {
		fmt.Println(k, v)
	}
	userStream := JustN(17).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   "user-" + fmt.Sprint(v),
			areaId: v % 7,
		}
	})
	userMap := ToGroupBy[int, User](userStream, func(user User) int {
		return user.areaId
	})
	for k, v := range userMap {
		fmt.Println(k, v)
	}
}

func TestGroupByMap(t *testing.T) {
	userStream := JustN(17).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   "user-" + fmt.Sprint(v),
			areaId: v % 7,
		}
	})
	userMap := ToGroupByMap[User, int, string](userStream, func(user User) int {
		return user.areaId
	}, func(user User) string {
		return user.name
	})
	for k, v := range userMap {
		fmt.Println(k, v)
	}
}
