package ctool

import (
	"sync"
)

// 双向链表维护使用情况：每次访问过的数据都置于链表头，容量满时清除链表尾节点数据

type (
	Cache[T any] interface {
		Get(string) T
		Save(string, T)
		MaxSize() int
		Size() int
	}

	LRUCache[T any] struct {
		maxSize int
		pool    map[string]*DoublyLinkedNode[Entry[T]]
		list    *DoublyLinkedList[Entry[T]]
		mutex   sync.RWMutex
	}
	Entry[T any] struct {
		key   string
		value T
	}
)

func NewLRUCache[T any](maxSize int) *LRUCache[T] {
	cachePool := make(map[string]*DoublyLinkedNode[Entry[T]])
	cacheQueue := NewEmptyDoublyLinkedList[Entry[T]]()
	return &LRUCache[T]{maxSize: maxSize, pool: cachePool, list: cacheQueue}
}

func (L *LRUCache[T]) Get(key string) (t T) {
	L.mutex.RLock()
	defer L.mutex.RUnlock()
	node := L.pool[key]
	if node == nil {
		// 使用泛型导致无法返回nil，只能返回类型的零值，对上层业务使用带来了限制
		return
	}

	L.resetHead(key)
	return node.data.value
}

func (L *LRUCache[T]) resetHead(key string) *DoublyLinkedNode[Entry[T]] {
	node := L.pool[key]
	if node == nil {
		return node
	}
	L.list.RemoveNode(node)
	L.list.AddToHead(node.data)
	return node
}

func (L *LRUCache[T]) Save(key string, val T) {
	L.mutex.Lock()
	defer L.mutex.Unlock()

	node := L.resetHead(key)
	if node != nil {
		entry := node.data
		entry.value = val
		L.pool[key] = L.list.head
		return
	}

	//fmt.Println(s, L.Size(), L.MaxSize())
	if L.Size() >= L.MaxSize() {
		node := L.list.tail
		L.list.RemoveTail()
		entry := node.data
		L.remove(entry.key)
	}
	L.list.AddToHead(Entry[T]{key: key, value: val})
	L.pool[key] = L.list.head
}

func (L *LRUCache[T]) MaxSize() int {
	return L.maxSize
}

func (L *LRUCache[T]) Size() int {
	return len(L.pool)
}

func (L *LRUCache[T]) remove(key string) {
	delete(L.pool, key)
}
