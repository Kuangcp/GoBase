package ctool

import (
	"encoding/json"
	"fmt"
	"testing"
)

func c() ResultVO[string] {
	return SuccessWith("sss")
}

func d() ResultVO[int] {
	r := c()
	if r.IsFailed() {
		return FailedWithMsg[int](r.Msg)
	}
	return SuccessWith(88)
}

func TestCheck(t *testing.T) {
	type B struct {
		A string `json:"a"`
		C string `json:"c"`
	}
	with := SuccessWith(B{A: "ifdjsifjds"})
	marshal, _ := json.Marshal(with)
	fmt.Println(string(marshal))

	fmt.Println(with.JSONStr(), string([]byte{}))
}
