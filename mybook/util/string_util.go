package util

import "strconv"

func BuildCHCharFormat(expectLen int, str string) string {
	return "%" + strconv.Itoa(expectLen-len(str)/3*2) + "s"
}
