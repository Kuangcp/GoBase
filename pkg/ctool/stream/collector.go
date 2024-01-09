package stream

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"reflect"
)

var (
	toStringType = ctool.NewSet(reflect.Uint, reflect.Int, reflect.Uint8, reflect.Int8, reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32, reflect.Uint64, reflect.Int64)
)

func ToSum[R ctool.Integer](s Stream) R {
	var r R
	for item := range s.source {
		r += item.(R)
	}
	return r
}

func ToSet[R comparable](s Stream) *ctool.Set[R] {
	return ToSetFunc[R](s, nil)
}

func ToSetFunc[R comparable](s Stream, fn func(s any) R) *ctool.Set[R] {
	result := ctool.NewSet[R]()
	for item := range s.source {
		if fn == nil {
			result.Add(item.(R))
		} else {
			result.Add(fn(item))
		}
	}
	return result
}

func ToList[R any](s Stream) []R {
	return ToListFunc[R](s, nil)
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

func ToMap[K comparable, V any](s Stream, key func(any) K, val func(any) V) map[K]V {
	result := make(map[K]V)
	for item := range s.source {
		result[key(item)] = val(item)
	}
	return result
}

func ToJoin(s Stream) string {
	return ToJoins(s, "")
}

func ToJoins(s Stream, split string) string {
	result := ""
	nonString := false
	first := true
	for item := range s.source {
		iType := reflect.TypeOf(item)
		if iType.Kind() == reflect.String {
			if !first {
				result += split
			}
			first = false
			result += item.(string)
		} else if toStringType.Contains(iType.Kind()) {
			if !first {
				result += split
			}
			first = false
			result += fmt.Sprint(item)
		} else {
			nonString = true
		}
	}
	if nonString {
		fmt.Println("warn: has no string type item")
	}
	return result
}
