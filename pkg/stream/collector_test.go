package stream

import (
	"fmt"
	"testing"
)

func TestCollectList(t *testing.T) {
	result := Just(1, 8, 10, 11, 20, 21).Map(func(item any) any {
		v := item.(int)
		return User{
			id:     v,
			name:   fmt.Sprint(v),
			areaId: v / 3,
		}
	}).Collect(ToListSelf())
	fmt.Println(result)
	ls := result.([]any)
	for _, i := range ls {
		u := i.(User)
		fmt.Println(u.id, u.areaId)
	}
}
