package linkedlist

import (
	"sync"
)

type (
	Cache interface {
		Get(string) interface{}
		Save(string, interface{})
		MaxSize() int
		Size() int
	}

	LRUCache struct {
		maxSize     int
		pool        map[string]*DoublyLinkedNode
		list        *DoublyLinkedList
		mutex       sync.RWMutex
		removeCount int
	}
	Entry struct {
		key   string
		value interface{}
	}
)

func NewLRUCache(maxSize int) *LRUCache {
	cachePool := make(map[string]*DoublyLinkedNode)
	cacheQueue := NewEmptyDoublyLinkedList()
	return &LRUCache{maxSize: maxSize, pool: cachePool, list: cacheQueue}
}

func (L *LRUCache) Get(key string) interface{} {
	L.mutex.RLock()
	defer L.mutex.RUnlock()
	node := L.pool[key]
	if node == nil {
		return nil
	}

	L.resetHead(key)
	return node.data.(Entry).value
}

func (L *LRUCache) resetHead(key string) *DoublyLinkedNode {
	node := L.pool[key]
	if node == nil {
		return node
	}
	L.list.RemoveNode(node)
	L.list.AddToHead(node.data)
	return node
}

func (L *LRUCache) Save(key string, val interface{}) {
	L.mutex.Lock()
	defer L.mutex.Unlock()

	node := L.resetHead(key)
	if node != nil {
		entry := node.data.(Entry)
		entry.value = val
		L.pool[key] = L.list.head
		return
	}

	//fmt.Println(s, L.Size(), L.MaxSize())
	if L.Size() >= L.MaxSize() {
		node := L.list.tail
		L.list.RemoveTail()
		entry := node.data.(Entry)
		L.remove(entry.key)
	}
	L.list.AddToHead(Entry{key: key, value: val})
	L.pool[key] = L.list.head
}

func (L *LRUCache) MaxSize() int {
	return L.maxSize
}

func (L *LRUCache) Size() int {
	return len(L.pool)
}

func (L *LRUCache) remove(key string) {
	delete(L.pool, key)
}

// map not gc?
func (L *LRUCache) remove2(key string) {
	L.removeCount++
	delete(L.pool, key)

	if L.removeCount > L.maxSize*7 {
		oldPool := L.pool
		L.pool = make(map[string]*DoublyLinkedNode)
		for k, node := range oldPool {
			L.pool[k] = node
		}
		L.removeCount = 0
	}
}
