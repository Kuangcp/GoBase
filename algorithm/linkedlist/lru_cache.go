package linkedlist

type (
	Cache interface {
		Get(string) interface{}
		Save(string, interface{})
		MaxCachePool() int
	}

	// map & double linked list
	// save or get: move node to head
	// clean: delete tail node
	LRUCache struct {
		maxPool    int
		cachePool  map[string]interface{}
		cacheQueue *DoublyLinkedList
	}
)

func NewLRUCache(maxPool int) *LRUCache {
	cachePool := make(map[string]interface{})
	cacheQueue := NewEmptyDoublyLinkedList()
	return &LRUCache{maxPool: maxPool, cachePool: cachePool, cacheQueue: cacheQueue}
}

func (L LRUCache) Get(s string) interface{} {
	panic("implement me")
}

func (L LRUCache) Save(s string, i interface{}) {
	panic("implement me")
}

func (L LRUCache) MaxCachePool() int {
	panic("implement me")
}
