package heap

type (
	Item interface {
		Value() int
	}

	MinHeap struct {
		Values []Item
	}
	IntItem struct {
		value int
	}
)

func NewIntItem(value int) *IntItem {
	return &IntItem{value: value}
}
func (t *IntItem) Value() int {
	return t.value
}

func (t *MinHeap) insertValue(value Item) {
	if value == nil {
		return
	}

	t.Values = append(t.Values, value)
	lens := len(t.Values)
	parentIdx := (lens - 1) / 2
	insertIdx := lens - 1
	for insertIdx >= 1 {
		if t.Values[parentIdx].Value() > value.Value() && parentIdx >= 1 {
			t.Values[insertIdx] = t.Values[parentIdx]
			insertIdx = parentIdx
			parentIdx = insertIdx / 2
		} else {
			t.Values[insertIdx] = value
			break
		}
	}
}
