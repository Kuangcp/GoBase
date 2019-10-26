package linkedlist

import (
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
	li := NewLinkedList(nil, nil)
	for i := 0; i < 4; i++ {
		li.Add(i)
	}
	//li.Show()

	//for i := 0; i < 5; i++ {
	//	li.Add(NewCar(i, 4, "benz"))
	//}
	//li.Show()

	li.Remove(0)
	li.Show()
	li.Remove(1)
	li.Show()
	li.Remove(2)
	li.Show()
	li.Remove(3)
	li.Show()
}
