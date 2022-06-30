package ctk

import (
	"fmt"
	"testing"
)

func TestLoop(t *testing.T) {
	v := make(map[string]string)
	v["xx"] = "rr"

	for s := range v {
		fmt.Println(s)
	}
}
