package ctk

type Set struct {
	cache map[interface{}]struct{}
}

func NewSet(s ...interface{}) *Set {
	result := Set{cache: make(map[interface{}]struct{})}
	if s != nil {
		for _, v := range s {
			result.Add(v)
		}
	}
	return &result
}

func (s *Set) Add(val ...interface{}) {
	if val == nil {
		return
	}
	for _, v := range val {
		s.cache[v] = struct{}{}
	}
}

func (s *Set) Adds(set *Set) {
	if set == nil {
		return
	}
	for k := range set.cache {
		s.cache[k] = struct{}{}
	}
}

// 交集
func (s *Set) Intersect(set *Set) *Set {
	if s == nil || set == nil {
		return nil
	}

	result := NewSet()
	for k := range s.cache {
		if set.Contains(k) {
			result.Add(k)
		}
	}

	return result
}

// 差集
func (s *Set) Difference(set *Set) *Set {
	return nil
}

// 并集
func (s *Set) Union(set *Set) *Set {
	return nil
}

// 补集 余集
func (s *Set) Supplementary(set *Set) *Set {
	if s == nil || set == nil {
		return nil
	}

	for k := range set.cache {
		s.Remove(k)
	}
	return s
}

func (s *Set) Len() int {
	return len(s.cache)
}

func (s *Set) IsEmpty() bool {
	return len(s.cache) == 0
}

func (s *Set) Contains(val interface{}) bool {
	_, ok := s.cache[val]
	return ok
}

func (s *Set) Remove(val interface{}) {
	delete(s.cache, val)
}

func (s *Set) Clear() {
	s.cache = make(map[interface{}]struct{})
}

func (s *Set) Loop(action func(interface{})) {
	for k, _ := range s.cache {
		action(k)
	}
}

func (s *Set) Intersection(set *Set) *Set {
	if set == nil {
		return nil
	}
	result := NewSet()
	for k := range s.cache {
		_, ok := set.cache[k]
		if ok {
			result.Add(k)
		}
	}
	return result
}
