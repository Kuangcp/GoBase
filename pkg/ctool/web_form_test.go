package ctool

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

// 使用反射获取注解内容, 类型转换后将值注入结构体内
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `form:"l"`
		MaxResults int      `form:"max"`
		Exact      bool     `form:"x"`
	}
	data.MaxResults = 10 // set default
	if err := Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// ...rest of handler...
	fmt.Fprintf(resp, "Search: %+v\n", data)
}

func TestUnpack(t *testing.T) {
	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
