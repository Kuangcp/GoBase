package ctool

import (
	"math/rand"
)

var list = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "I", "i",
	"J", "j", "K", "k", "L", "l", "M", "m", "N", "n", "O", "o", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t",
	"U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z"}
var alphaAndNum = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "I", "i",
	"J", "j", "K", "k", "L", "l", "M", "m", "N", "n", "O", "o", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t",
	"U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var alphaDict = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
var alphaDictLen = 52
var alphanumericDict = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
var alphanumericDictLen = 62

func RandomAlpha(length int) string {
	bys := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		idx := rand.Intn(alphaDictLen)
		bys = append(bys, alphaDict[idx:idx+1]...)
	}
	return string(bys)
}

func RandomAlNum(length int) string {
	bys := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		idx := rand.Intn(alphanumericDictLen)
		bys = append(bys, alphanumericDict[idx:idx+1]...)
	}

	return string(bys)
}

// RandomAlNumValid 生成非数字开头的字符串，数据库等地方使用
func RandomAlNumValid(length int) string {
	bys := make([]byte, 0, length)
	ls := 0
	for ls < length {
		idx := rand.Intn(alphanumericDictLen)
		single := alphanumericDict[idx : idx+1]
		if ls == 0 && single >= "0" && single <= "9" {
			continue
		}
		bys = append(bys, single...)
		ls++
	}

	return string(bys)
}
