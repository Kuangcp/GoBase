package rx

import (
	"fmt"
	"github.com/reactivex/rxgo/v2"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	observable := rxgo.Just("Hello, World!")()
	ch := observable.Observe()
	item := <-ch
	fmt.Println(item.V)
}
