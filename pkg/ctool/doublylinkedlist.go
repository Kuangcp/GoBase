package ctool

import (
	"log"
)

type DoublyLinkedList[T comparable] struct {
	head *DoublyLinkedNode[T]
	tail *DoublyLinkedNode[T]
	len  int
}

type DoublyLinkedNode[T comparable] struct {
	pre  *DoublyLinkedNode[T]
	data T
	next *DoublyLinkedNode[T]
}

func NewEmptyDoublyLinkedList[T comparable]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{head: nil, tail: nil}
}

func NewDoublyLinkedList[T comparable](head, tail *DoublyLinkedNode[T]) *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{head: head, tail: tail}
}

func NewEmptyDoublyLinkedNode[T comparable](data T) *DoublyLinkedNode[T] {
	return &DoublyLinkedNode[T]{pre: nil, data: data, next: nil}
}

func NewDoublyLinkedNode[T comparable](pre, next *DoublyLinkedNode[T], data T) *DoublyLinkedNode[T] {
	return &DoublyLinkedNode[T]{pre: pre, data: data, next: next}
}

func (list *DoublyLinkedList[T]) IsEmpty() bool {
	return list.head == nil
}

func (list *DoublyLinkedList[T]) AddToHead(data T) {
	node := NewEmptyDoublyLinkedNode(data)
	if list.head == nil {
		list.head = node
		list.tail = node
	} else {
		list.head.pre = node
		node.next = list.head
		list.head = node
	}
	list.len++
}

// Add to tail
func (list *DoublyLinkedList[T]) Add(data T) {
	node := NewEmptyDoublyLinkedNode[T](data)
	if list.head == nil {
		list.head = node
		list.tail = node
	} else {
		list.tail.next = node
		node.pre = list.tail
		list.tail = node
	}
	list.len++
}

func (list *DoublyLinkedList[T]) Find(data T) *DoublyLinkedNode[T] {
	if list.IsEmpty() {
		return nil
	}

	node := list.head
	for {
		if node.data == data {
			return node
		}
		if node.next != nil {
			node = node.next
		} else {
			break
		}
	}
	return nil
}

func (list *DoublyLinkedList[T]) Clear() {
	list.head = nil
	list.tail = nil
	list.len = 0
}

func (list *DoublyLinkedList[T]) RemoveTail() {
	if list.IsEmpty() {
		return
	}

	if list.head == list.tail {
		list.Clear()
		return
	}
	sec := list.tail.pre
	list.tail = sec
	sec.next = nil
	list.len--
}

func (list *DoublyLinkedList[T]) RemoveNode(node *DoublyLinkedNode[T]) {
	if node == nil {
		return
	}
	pre := node.pre
	next := node.next

	// isHead
	if pre == nil {
		list.head = next
		if list.head != nil {
			list.head.pre = nil
		}
		list.len--
		return
	}

	// isTail
	if next == nil {
		list.tail = pre
		pre.next = nil
		list.len--
		return
	}

	// remove current node
	pre.next = next
	next.pre = pre
	list.len--
}

func (list *DoublyLinkedList[T]) Remove(data T) {
	if list.IsEmpty() {
		return
	}

	node := list.Find(data)
	list.RemoveNode(node)
}

// 单链表反转 三个指针前进
func (list *DoublyLinkedList[T]) ReverseBySingle() *DoublyLinkedList[T] {
	if list.IsEmpty() || list.len == 1 {
		return list
	}

	first := list.head
	second := list.head.next

	if list.len == 2 {
		list.head = second
		second.next = first
		first.next = nil
		list.tail = first
		return list
	}

	first.next = nil
	list.tail = first
	third := second.next
	for {
		second.next = first

		if third.next == nil {
			third.next = second
			list.head = third
			return list
		}
		first = second
		second = third
		third = third.next
	}
}

func (list *DoublyLinkedList[T]) PrintList() {
	list.PrintListWithDetail(false)
}

func (list *DoublyLinkedList[T]) PrintListWithDetail(needDetail bool) {
	if list.IsEmpty() {
		log.Println("list is empty")
		return
	}

	node := list.head
	log.Println("len:", list.len, ",head:", list.head, ",tail:", list.tail)

	for {
		if needDetail {
			log.Println(
				Green, "pre=", node.pre, End,
				"data=", node.data,
				Green, "next=", node.next, End)
		} else {
			log.Println(node.data)
		}
		if node.next != nil {
			node = node.next
		} else {
			break
		}
	}
	println()
}
