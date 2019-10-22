package cuibase

import (
	"log"
	"os"
)

var red = "\033[0;31m"
var green = "\033[0;32m"
var yellow = "\033[0;33m"
var blue = "\033[0;34m"
var purple = "\033[0;35m"
var cyan = "\033[0;36m"
var white = "\033[0;37m"
var end = "\033[0m"

// AssertParamCount os.Args 参数构成: 0 文件 1 参数 2 参数
func AssertParamCount(count int, msg string) {
	param := os.Args
	flag := enoughCount(param, count)
	if !flag {
		log.Printf("param count less than %v \n", count)
		log.Fatal(msg)
	}
}

func enoughCount(param []string, count int) bool {
	return len(param) > count
}
