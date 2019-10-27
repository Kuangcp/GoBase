package linkedlist

import (
	"log"
	"strconv"
	"testing"
)

type Car struct {
	no    int
	wheel int
	brand string
}

func NewCar(no int, wheel int, brand string) *Car {
	return &Car{no: no, wheel: wheel, brand: brand}
}

// string interface like toString() in Java
func (car *Car) String() string {
	return "(no=" + strconv.Itoa(car.no) + ", brand=" + car.brand + ", wheel=" + strconv.Itoa(car.wheel) + ")"
}

func TestLinkedList_Add(t *testing.T) {
	list := NewLinkedList(nil, nil)
	for i := 0; i < 4; i++ {
		list.Add(i)
	}
	node := list.Find(2)
	re := node.data.(int) + 2
	log.Println(node.data, re)

	list.PrintList()
	list.Remove(3)
	list.PrintList()
	list.Clear()

	for i := 0; i < 5; i++ {
		list.Add(NewCar(i, 4, "benz"))
	}
	list.PrintList()

	car := list.Find(NewCar(1, 4, "benz"))
	log.Println(car)
}

func TestLinkedList_Reverse(t *testing.T) {
	list := NewLinkedList(nil, nil)
	for i := 0; i < 2; i++ {
		list.Add(i)
	}
	list.ReverseBySingle().PrintList()
	list.Clear()

	for i := 0; i < 5; i++ {
		list.Add(i)
	}
	list.ReverseBySingle().PrintList()
}