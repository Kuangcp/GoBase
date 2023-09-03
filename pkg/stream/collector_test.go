package stream

import (
	"fmt"
	"testing"
)

func TestCollectList(t *testing.T) {
	result := ToListSelf[User](Just(1, 8, 10, 11, 20, 21).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   fmt.Sprint(v),
			areaId: v / 3,
		}
	}))
	fmt.Println(result)
	for _, u := range result {
		fmt.Println(u.id, u.areaId)
	}
}
