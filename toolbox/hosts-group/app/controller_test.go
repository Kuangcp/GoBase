package app

import (
	"fmt"
	"strings"
	"testing"
)

func Test_buildOneFileBlock(t *testing.T) {
	result := buildFileBlock("tes", "fdsfsd\nfdsjsi\n\nfdsjk")
	fmt.Println(result)
}

func TestTrim(t *testing.T) {
	fmt.Println(strings.TrimSpace("new s"))
}
