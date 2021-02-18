package app

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

func Test_buildOneFileBlock(t *testing.T) {
	fmt.Print(buildFileBlock("one", "fdsfsd\nfdsjsi\n\nfdsjk"))
	fmt.Print(buildFileBlock("two", "1 line\n 2 line \n key"))
}

func TestTrim(t *testing.T) {
	fmt.Println(strings.TrimSpace("new s"))
	fmt.Println(utf8.RuneCountInString("级级级级级级级级"))
}
