package stream

import "fmt"

func ToString(item any) any {
	return fmt.Sprint(item)
}

func Self[R any](a any) R {
	return a.(R)
}
