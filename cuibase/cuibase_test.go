package main

import (
	"testing"
)

func TestAssertCount(t *testing.T) {
	list := make([]string, 1)
	flag := enoughCount(list, 2)
	if flag {
		t.Fail()
	}
	t.Log(flag)

	list = make([]string, 2)
	flag = enoughCount(list, 1)
	if !flag {
		t.Fail()
	}
}
