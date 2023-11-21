package ctool

import (
	"testing"
)

func TestNewMap(t *testing.T) {
	newMap := NewMap(MapEntry[int, string]{3, "object"})
	s := newMap.Get(2)
	println(s, newMap.Contain(2))
	println(newMap.Get(3), newMap.Contain(3))
}

func TestMapsPut(t *testing.T) {
	maps := NewMaps(MapsEntry[int, string]{Key: 5, Val: []string{"sss"}})
	maps.Put(5, "sss", "fffffff")
	for k, v := range maps.cache {
		println(k, v)
	}
}
