package util

import (
	"encoding/json"
	"strconv"
)

func BuildCHCharFormat(expectLen int, str string) string {
	return "%" + strconv.Itoa(expectLen-len(str)/3*2) + "s"
}

func Json(data interface{}) string {
	bytes, e := json.Marshal(data)
	if e != nil {
		return " ERROR"
	}
	return string(bytes)
}
