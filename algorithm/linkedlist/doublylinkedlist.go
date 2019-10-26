package linkedlist

import (
	"github.com/kuangcp/gobase/cuibase"
	"log"
)

type LinkedList struct {
	head *DoublyLinkedNode
	tail *DoublyLinkedNode
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
	return list == nil || list.head == nil
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
}

func (list *LinkedList) Find(data interface{}) *DoublyLinkedNode {
	if list == nil || list.IsEmpty() {
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
func (list *LinkedList) Remove(data interface{}) {
	if list == nil || list.IsEmpty() {
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
			return
		}

		// isTail
		if next == nil {
			list.tail = pre
			pre.next = nil
			return
		}

		// remove current node
		pre.next = next
		next.pre = pre
	}
}

func (list *LinkedList) Show() {
	list.ShowWithDetail(false)
}
func (list *LinkedList) ShowWithDetail(needDetail bool) {
	if list == nil || list.IsEmpty() {
		log.Println("list is empty")
		return
	}

	node := list.head
	println()
	log.Println("head=", list.head, "tail=", list.tail)

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
