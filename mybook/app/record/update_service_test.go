package record

import (
	"github.com/kuangcp/gobase/pkg/ghelp"
	"log"
	"testing"
)

func TestParseFloat(t *testing.T) {
	assertValue(parsePrice("123").Data, 12300)
	assertValue(parsePrice("10").Data, 1000)

	assertValue(parsePrice("3.3").Data, 330)
	assertValue(parsePrice("3.03").Data, 303)
	assertValue(parsePrice("3.00").Data, 300)
	assertValue(parsePrice("3.30").Data, 330)
	assertValue(parsePrice("3.33").Data, 333)

	assertFailedResult(parsePrice("0"))
	assertFailedResult(parsePrice("-0"))

	assertFailedResult(parsePrice("3.000"))
	assertFailedResult(parsePrice("3.333"))
	assertFailedResult(parsePrice("3.3.99"))
	assertFailedResult(parsePrice("3..3"))
	assertFailedResult(parsePrice("-3.9"))
	assertFailedResult(parsePrice("-3.933333"))
	assertFailedResult(parsePrice("-3.933333"))

	assertFailedResult(parsePrice("test"))
	assertFailedResult(parsePrice("3.33d"))
	assertFailedResult(parsePrice("x3.33"))
}

func assertFailedResult(vo ghelp.ResultVO) {
	if vo.IsSuccess() {
		log.Fatalf("error: %v", vo)
	}
}
func assertValue(actual interface{}, except interface{}) {
	if actual != except {
		log.Fatalf("except: %v actual: %v\n", except, actual)
	}
}
