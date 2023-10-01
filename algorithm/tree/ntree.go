package main

type (
	NTree[T any] struct {
		Parent  *NTree[T]
		Childes []*NTree[T]
		Data    T
	}
)
