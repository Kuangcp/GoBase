package linkedlist

import (
	"github.com/kuangcp/gobase/cuibase"
	"log"
)

type LinkedList struct {
	head *DoublyLinkedNode
	tail *DoublyLinkedNode
	len  int
}

func NewLinkedList(head *DoublyLinkedNode, tail *DoublyLinkedNode) *LinkedList {
	return &LinkedList{head: head, tail: tail}
}

type DoublyLinkedNode struct {
	pre  *DoublyLinkedNode
	data interface{}
	next *DoublyLinkedNode
}

func NewDoublyLinkedNode(pre *DoublyLinkedNode, data interface{}, next *DoublyLinkedNode) *DoublyLinkedNode {
	return &DoublyLinkedNode{pre: pre, data: data, next: next}
}

func (list *LinkedList) IsEmpty() bool {
	return list.head == nil
}

func (list *LinkedList) Add(data interface{}) {
	node := NewDoublyLinkedNode(nil, data, nil)
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

func (list *LinkedList) Find(data interface{}) *DoublyLinkedNode {
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

func (list *LinkedList) Clear() {
	list.head = nil
	list.tail = nil
	list.len = 0
}

func (list *LinkedList) Remove(data interface{}) {
	if list.IsEmpty() {
		return
	}

	node := list.Find(data)
	if node != nil {
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
}

// 单链表反转 三个指针前进
func (list *LinkedList) ReverseBySingle() *LinkedList {
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

func (list *LinkedList) PrintList() {
	list.PrintListWithDetail(false)
}

func (list *LinkedList) PrintListWithDetail(needDetail bool) {
	if list.IsEmpty() {
		log.Println("list is empty")
		return
	}

	node := list.head
	println()
	log.Println("head=", list.head, "tail=", list.tail, "len=", list.len)

	for {
		if needDetail {
			print("cur=", node, " ")
			log.Println(cuibase.Green, "pre=", node.pre, cuibase.End,
				"data=", node.data, cuibase.Green, "next=", node.next, cuibase.End)
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
