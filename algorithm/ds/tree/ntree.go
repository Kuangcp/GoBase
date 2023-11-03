package main

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"strings"
)

type (
	NTree[T any] struct {
		Parent  *NTree[T]
		Childes []*NTree[T]
		Data    T
	}

	// n tree store in relation table
	ATable[T any] struct {
		Id       string `json:"id"`
		ParentId string `json:"parent_id"`
		Data     T      `json:"data"`
	}

	// ATree application n tree
	ATree[T any] struct {
		Id      string      `json:"id"`
		Childes []*ATree[T] `json:"childes"`
		Data    T           `json:"data"`
	}
)

func (a ATable[T]) String() string {
	return fmt.Sprintf("%6s[%6s] > %v", a.Id, a.ParentId, a.Data)
}

func (a ATree[T]) PrintJson() {
	marshal, _ := json.Marshal(a)
	//fmt.Println(string(marshal))

	var Options = &pretty.Options{Width: 80, Prefix: "", Indent: "    ", SortKeys: false}
	fmt.Println(string(pretty.PrettyOptions(marshal, Options)))
}

// dirAll 如果关键字命中中间层节点带出所有子节点
func (t *ATree[T]) Search(kwd string, dirAll bool) bool {
	if strings.Contains(fmt.Sprint(t.Data), kwd) && dirAll {
		return true
	}
	if len(t.Childes) == 0 {
		return false
	}
	var rm []int
	anyMatched := false
	for i, s := range t.Childes {
		matched := s.Search(kwd, dirAll)
		if !matched {
			rm = append(rm, i)
		} else {
			anyMatched = true
		}
	}
	//fmt.Println("delete: ", t.Id, rm)
	for i := len(rm) - 1; i >= 0; i-- {
		idx := rm[i]
		t.Childes = append(t.Childes[:idx], t.Childes[idx+1:]...)
	}

	return anyMatched
}
