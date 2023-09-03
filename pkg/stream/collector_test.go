package stream

import (
	"fmt"
	"testing"
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
		Just(1, 8, 10, 11, 20, 21).Map(func(item any) any {
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
		}, Self[User](),
	)
	for k, v := range userMap {
		fmt.Println(k, v)
	}
}
