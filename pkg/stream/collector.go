package stream

type (
	Collector struct {
		s      Stream
		action func() any
	}
)

func (c *Collector) collect() any {
	return c.action()
}

func ToListSelf() *Collector {
	return ToList(nil)
}
func ToList(fn MapFunc) *Collector {
	c := &Collector{}
	c.action = func() any {
		var result []any
		for item := range c.s.source {
			if fn == nil {
				result = append(result, item)
			} else {
				result = append(result, fn(item))
			}
		}
		return result
	}
	return c
}
