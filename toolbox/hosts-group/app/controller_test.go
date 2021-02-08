package app

import (
	"fmt"
	"testing"
)

func Test_buildOneFileBlock(t *testing.T) {
	result := buildFileBlock("tes", "fdsfsd\nfdsjsi\n\nfdsjk")
	fmt.Println(result)
}
