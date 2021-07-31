package linkedlist

import (
	"log"

	"github.com/kuangcp/gobase/cuibase"
)

type DoublyLinkedList struct {
	head *DoublyLinkedNode
	tail *DoublyLinkedNode
	len  int
}

type DoublyLinkedNode struct {
	pre  *DoublyLinkedNode
	data interface{}
	next *DoublyLinkedNode
}

func NewEmptyDoublyLinkedList() *DoublyLinkedList {
	return &DoublyLinkedList{head: nil, tail: nil}
}

func NewDoublyLinkedList(head, tail *DoublyLinkedNode) *DoublyLinkedList {
	return &DoublyLinkedList{head: head, tail: tail}
}

func NewEmptyDoublyLinkedNode(data interface{}) *DoublyLinkedNode {
	return &DoublyLinkedNode{pre: nil, data: data, next: nil}
}

func NewDoublyLinkedNode(pre, next *DoublyLinkedNode, data interface{}) *DoublyLinkedNode {
	return &DoublyLinkedNode{pre: pre, data: data, next: next}
}

func (list *DoublyLinkedList) IsEmpty() bool {
	return list.head == nil
}

func (list *DoublyLinkedList) AddToHead(data interface{}) {
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
func (list *DoublyLinkedList) Add(data interface{}) {
	node := NewEmptyDoublyLinkedNode(data)
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

func (list *DoublyLinkedList) Find(data interface{}) *DoublyLinkedNode {
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

func (list *DoublyLinkedList) Clear() {
	list.head = nil
	list.tail = nil
	list.len = 0
}

func (list *DoublyLinkedList) RemoveTail() {
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

func (list *DoublyLinkedList) RemoveNode(node *DoublyLinkedNode) {
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

func (list *DoublyLinkedList) Remove(data interface{}) {
	if list.IsEmpty() {
		return
	}

	node := list.Find(data)
	list.RemoveNode(node)
}

// 单链表反转 三个指针前进
func (list *DoublyLinkedList) ReverseBySingle() *DoublyLinkedList {
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

func (list *DoublyLinkedList) PrintList() {
	list.PrintListWithDetail(false)
}

func (list *DoublyLinkedList) PrintListWithDetail(needDetail bool) {
	if list.IsEmpty() {
		log.Println("list is empty")
		return
	}

	node := list.head
	log.Println("len:", list.len, ",head:", list.head, ",tail:", list.tail)

	for {
		if needDetail {
			log.Println(
				cuibase.Green, "pre=", node.pre, cuibase.End,
				"data=", node.data,
				cuibase.Green, "next=", node.next, cuibase.End)
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
