package ctk

import "fmt"

type Value interface {
	String() string
	Set(string) error
}

type ArrayFlags []string

// Value ...
func (i *ArrayFlags) String() string {
	return fmt.Sprint(*i)
}

// Set 方法是flag.Value接口, 设置flag Value的方法.
// 通过多个flag指定的值， 所以我们追加到最终的数组上.
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
