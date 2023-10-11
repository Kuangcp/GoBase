package main

import (
	"fmt"
)

type (
	BaseBird struct {
		age int
	}

	DerivedBird struct {
		BaseBird
	}
)

func (this *BaseBird) Cal() {
	this.Add()
}
func (this *BaseBird) Add() {
	this.age = this.age + 1
	fmt.Printf("base add. age=%d\n", this.age)
}

func (this *DerivedBird) Add() {
	this.age = this.age + 2
	fmt.Printf("add. age=%d\n", this.age)
}

func (this *BaseBird) AddWithInterface() {
	this.age = this.age + 1
	fmt.Printf("base interface add. age=%d\n", this.age)
}

func (this *DerivedBird) AddWithInterface() {
	this.age = this.age + 2
	fmt.Printf("interface add. age=%d\n", this.age)
}

type Bird interface {
	AddWithInterface()
}

func Cal(bird Bird) {
	bird.AddWithInterface()
}

// https://zhuanlan.zhihu.com/p/157786743
func main() {
	// 使用组合的方式 实现不了
	var b1 = BaseBird{age: 1}
	var b2 = DerivedBird{BaseBird{1}}
	b1.Cal()
	b2.Cal()

	// 使用接口的方式实现
	Cal(&BaseBird{age: 1})
	Cal(&DerivedBird{BaseBird{1}})
}
