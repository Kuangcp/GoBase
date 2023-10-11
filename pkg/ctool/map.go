package ctool

import "sort"
import "golang.org/x/exp/constraints"

type (
	MapEntry[K, V comparable] struct {
		Key K
		Val V
	}
	Map[K, V comparable] struct {
		cache map[K]V
	}
)

func NewMap[K, V comparable](k ...MapEntry[K, V]) *Map[K, V] {
	obj := &Map[K, V]{cache: make(map[K]V)}
	if k != nil {
		for _, m := range k {
			obj.cache[m.Key] = m.Val
		}
	}
	return obj
}

func (m *Map[K, V]) Contain(k K) bool {
	_, ok := m.cache[k]
	return ok
}

func (m *Map[K, V]) Get(k K) V {
	return m.cache[k]
}

type (
	MapsEntry[K, V comparable] struct {
		Key K
		Val []V
	}
	Maps[K, V comparable] struct {
		cache map[K][]V
	}
)

func NewMaps[K, V comparable](k ...MapsEntry[K, V]) *Maps[K, V] {
	obj := &Maps[K, V]{cache: make(map[K][]V)}
	if k != nil {
		for _, m := range k {
			obj.cache[m.Key] = m.Val
		}
	}
	return obj
}

func (m *Maps[K, V]) Contain(k K) bool {
	_, ok := m.cache[k]
	return ok
}

func (m *Maps[K, V]) Get(k K) []V {
	return m.cache[k]
}

func (m *Maps[K, V]) Put(k K, v ...V) bool {
	if len(v) == 0 {
		return false
	}
	vs, ok := m.cache[k]
	if !ok {
		m.cache[k] = v
	}
	var vss []string
	vss = append(vss, "ss")

	vs = append(vs, v...)
	m.cache[k] = vs
	return true
}

func NewSortMap[K comparable, V constraints.Ordered](origin map[K]V) []MapEntry[K, V] {
	var result []MapEntry[K, V]
	for k, v := range origin {
		result = append(result, MapEntry[K, V]{Key: k, Val: v})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Val < result[j].Val
	})
	return result
}
