package stream

import (
	"fmt"
	"reflect"
)

func ToList[R any](s Stream) []R {
	return ToListFunc[R](s, nil)
}

func ToMap[K comparable, V any](s Stream, key func(any) K, val func(any) V) map[K]V {
	result := make(map[K]V)
	for item := range s.source {
		result[key(item)] = val(item)
	}
	return result
}

func ToListFunc[R any](s Stream, fn func(s any) R) []R {
	var result []R
	for item := range s.source {
		if fn == nil {
			result = append(result, item.(R))
		} else {
			result = append(result, fn(item))
		}
	}
	return result
}

func ToJoin(s Stream) string {
	result := ""
	nonString := false
	for item := range s.source {
		iType := reflect.TypeOf(item)
		if iType.Kind() == reflect.String {
			result += item.(string)
		} else {
			nonString = true
		}
	}
	if nonString {
		fmt.Println("warn: has no string type item")
	}
	return result
}
