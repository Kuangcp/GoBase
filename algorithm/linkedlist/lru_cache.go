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
		maxSize    int
		cachePool  map[string]*DoublyLinkedNode
		cacheQueue *DoublyLinkedList
		mutex      sync.RWMutex
	}
	Entry struct {
		key   string
		value interface{}
	}
)

func NewLRUCache(maxSize int) *LRUCache {
	cachePool := make(map[string]*DoublyLinkedNode)
	cacheQueue := NewEmptyDoublyLinkedList()
	return &LRUCache{maxSize: maxSize, cachePool: cachePool, cacheQueue: cacheQueue}
}

func (L *LRUCache) Get(s string) interface{} {
	L.mutex.RLock()
	defer L.mutex.RUnlock()
	node := L.cachePool[s]
	if node == nil {
		return nil
	}
	// TODO put head
	return node.data.(Entry).value
}

func (L *LRUCache) Save(s string, i interface{}) {
	L.mutex.Lock()
	defer L.mutex.Unlock()

	//fmt.Println(s, L.Size(), L.MaxSize())
	if L.Size() >= L.MaxSize() {
		node := L.cacheQueue.tail
		L.cacheQueue.RemoveTail()
		entry := node.data.(Entry)
		delete(L.cachePool, entry.key)
		//fmt.Println("remove", entry.key)
	}
	L.cacheQueue.AddToHead(Entry{key: s, value: i})
	L.cachePool[s] = L.cacheQueue.head
}

func (L *LRUCache) MaxSize() int {
	return L.maxSize
}

func (L *LRUCache) Size() int {
	return len(L.cachePool)
}
