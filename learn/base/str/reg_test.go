package str

import (
	"log"
	"regexp"
	"testing"
)

func TestRegPossessive(t *testing.T) {
	// not support possessive mode
	compile, err := regexp.Compile("ab{1,3}+c")
	if err != nil {
		log.Panic(err)
	}

	mat := compile.MatchString("abbc")
	log.Println(mat)
}
