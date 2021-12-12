package record

import (
	"github.com/kuangcp/gobase/pkg/ghelp"
	"log"
	"testing"
)

func TestParseFloat(t *testing.T) {
	var result ghelp.ResultVO
	result = parseAmount("123")
	assertValue(result.Data, 12300)

	result = parseAmount("3.3")
	assertValue(result.Data, 330)

	result = parseAmount("3.33")
	assertValue(result.Data, 333)

	result = parseAmount("test")
	assertValue(result.IsFailed(), true)

	result = parseAmount("3.33d")
	assertValue(result.IsFailed(), true)
}

func assertValue(actual interface{}, except interface{}) {
	if actual != except {
		log.Fatalf("except: %v actual: %v\n", except, actual)
	}
}
