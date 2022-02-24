package cuibase

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

func (s *Set) Add(val interface{}) {
	if val == nil {
		return
	}
	s.cache[val] = struct{}{}
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
