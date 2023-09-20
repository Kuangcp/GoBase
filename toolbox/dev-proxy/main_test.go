package main

import (
	"fmt"
	"github.com/kuangcp/gobase/toolbox/dev-proxy/core"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestMatchPrefix(t *testing.T) {
	compile := regexp.MustCompile("(.*)//192.168.16.90:32009/tg-easy-fetch/(.*)")
	result := compile.ReplaceAllString("http://192.168.16.90:32009/tg-easy-fetch/index.html", "$1//localhost:8081/api/$2")
	fmt.Println("replace: ", result)

	result = compile.ReplaceAllString("http://192.168.16.90:32009/tg-easy-fetch2/index.html", "$1//localhost:8081/api/$2")
	fmt.Println("replace: ", result)

	submatch := compile.FindStringSubmatch("http://192.168.16.90:32009/tg-easy-fetch/index.html")
	fmt.Println(submatch)

	fmt.Println(strings.HasPrefix("xxx", ""))
}

func TestHostName(t *testing.T) {
	name, err := os.Hostname()
	println(name)
	println(err)
}

func TestTransGeneric(t *testing.T) {
	type A[T comparable] struct {
		id int
		no T
	}
	//a := A[string]{id: 3, no: "sss"}
	//b := (A[int])(a)
}

func TestPrefix(t *testing.T) {
	var s []byte
	s = []byte("-----------------------------38012652472498993413112866420\nContent-Disposition: form-data; name=\"tableName\"\n\n")
	println(len(s))

	fmt.Println([]byte("\n"))
	//result := filterFileType(s)
	result := core.FilterFormType(s)
	fmt.Println(string(result))

	start := time.Now().UnixMilli()
	for i := 0; i < 10000; i++ {
		judgeStr(s)
		//judgeByte(s)
		//filterFormType(s)
	}
	end := time.Now().UnixMilli()
	fmt.Println("cost: ", end-start)
}

func judgeStr(s []byte) bool {
	str := string(s)
	if strings.HasPrefix(str, "------") {
		return true
	}
	return false
}
