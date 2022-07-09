package ctool

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadLines(t *testing.T) {
	lines := ReadLines[int]("test.tsv", func(s string) bool {
		return true
	}, func(s string) int {
		sp := strings.Split(s, "\t")
		return len(sp)
	})
	fmt.Println(lines)

	linev := ReadTsvLines("test.tsv")
	for i := range linev {
		fmt.Println(len(linev[i]))
	}
}
