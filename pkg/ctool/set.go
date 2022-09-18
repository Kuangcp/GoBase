package ctool

type Set[T comparable] struct {
	cache map[T]struct{}
}

func NewSet[T comparable](s ...T) *Set[T] {
	result := Set[T]{cache: make(map[T]struct{})}
	if s != nil {
		for _, v := range s {
			result.Add(v)
		}
	}
	return &result
}

func (s *Set[T]) Add(val ...T) {
	if val == nil {
		return
	}
	for _, v := range val {
		s.cache[v] = struct{}{}
	}
}

func (s *Set[T]) Adds(set *Set[T]) {
	if set == nil {
		return
	}
	for k := range set.cache {
		s.cache[k] = struct{}{}
	}
}

// Intersect 交集
func (s *Set[T]) Intersect(set *Set[T]) *Set[T] {
	if s == nil || set == nil {
		return nil
	}

	result := NewSet[T]()
	for k := range s.cache {
		if set.Contains(k) {
			result.Add(k)
		}
	}

	return result
}

// Difference 差集: 在当前集合 但是 不在参数集合内
func (s *Set[T]) Difference(set *Set[T]) *Set[T] {
	result := NewSet[T]()
	for k := range s.cache {
		if !set.Contains(k) {
			result.Add(k)
		}
	}
	return result
}

// Union 并集
func (s *Set[T]) Union(set *Set[T]) *Set[T] {
	result := NewSet[T]()
	result.Adds(s)
	result.Adds(set)
	return result
}

// Supplementary 补集 余集
func (s *Set[T]) Supplementary(set *Set[T]) *Set[T] {
	if s == nil || set == nil {
		return nil
	}

	result := NewSet[T]()
	result.Adds(s)
	for k := range set.cache {
		result.Remove(k)
	}
	return result
}

func (s *Set[T]) Len() int {
	return len(s.cache)
}

func (s *Set[T]) IsEmpty() bool {
	return len(s.cache) == 0
}

func (s *Set[T]) Contains(val T) bool {
	_, ok := s.cache[val]
	return ok
}

func (s *Set[T]) Remove(val T) {
	delete(s.cache, val)
}

func (s *Set[T]) Clear() {
	s.cache = make(map[T]struct{})
}

func (s *Set[T]) Loop(action func(T)) {
	for k, _ := range s.cache {
		action(k)
	}
}
