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

func RandomAlpha(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(list))
		result += list[idx] + " "
	}

	return result
}

func RandomAlNum(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		idx := rand.Intn(len(alphaAndNum))
		result += alphaAndNum[idx]
	}

	return result
}

// RandomAlNumValid 生成非数字开头的字符串，数据库等地方使用
func RandomAlNumValid(length int) string {
	result := ""
	for len(result) < length {
		idx := rand.Intn(len(alphaAndNum))
		single := alphaAndNum[idx]
		if len(result) == 0 && single >= "0" && single <= "9" {
			continue
		}
		result += single
	}

	return result
}
