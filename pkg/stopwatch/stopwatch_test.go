package stopwatch

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestStopWatch_PrettyPrint(t *testing.T) {
	stopWatch := NewWithName("action1")
	stopWatch.Start("task1 new request")
	time.Sleep(time.Second * 1)
	stopWatch.Stop()
	stopWatch.Start("task2")
	http.Get("https://baidu.com")
	stopWatch.Stop()
	fmt.Println(stopWatch.PrettyPrint())

	stopWatch = NewWithName("action")
	stopWatch.Stop()
	fmt.Println(stopWatch.PrettyPrint())

	stopWatch = NewWithName("action2")
	stopWatch.Start("A")
	http.Get("http://jd.com")
	stopWatch.Stop()
	stopWatch.Start("B")
	fmt.Println(stopWatch.PrettyPrint())
}

func TestRepeatStart(t *testing.T) {
	watch := NewWithName("task")
	watch.Start("a")
	http.Get("http://jd.com")
	watch.Start("b")
	http.Get("http://jd.com")
	watch.Stop()
	println(watch.PrettyPrint())
}

func TestMissingLastStop(t *testing.T) {
	watch := NewWithName("task")
	watch.Start("a")
	http.Get("http://jd.com")
	println(watch.PrettyPrint())
}
