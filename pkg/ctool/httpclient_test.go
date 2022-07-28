package ctool

import (
	"log"
	"testing"
)

func TestPostJson(t *testing.T) {
	type v struct {
		A string `json:"a"`
		B int    `json:"b"`
	}
	rsp, err := PostJson("http://127.0.0.1:9911/webapi/ec/EC_GetRecord",
		v{A: "33", B: 43})
	if err != nil {
		return
	}

	log.Println(string(rsp))

}
