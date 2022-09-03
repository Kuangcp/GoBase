package ctool

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestStopWatch_PrettyPrint(t *testing.T) {
	stopWatch := NewStopWatchWithName("action1")
	stopWatch.Start("task1 new request")
	time.Sleep(time.Second * 3)
	stopWatch.Stop()
	stopWatch.Start("task2")
	http.Get("https://baidu.com")
	stopWatch.Stop()
	fmt.Println(stopWatch.PrettyPrint())

	stopWatch = NewStopWatchWithName("action")
	stopWatch.Stop()
	fmt.Println(stopWatch.PrettyPrint())

	stopWatch = NewStopWatchWithName("action2")
	stopWatch.Start("A")
	http.Get("http://jd.com")
	stopWatch.Stop()
	stopWatch.Start("B")

	fmt.Println(stopWatch.PrettyPrint())
}

func TestRepeatStart(t *testing.T) {
	watch := NewStopWatchWithName("task")
	watch.Start("a")
	http.Get("http://jd.com")
	watch.Start("b")
	http.Get("http://jd.com")
	watch.Stop()
	watch.Start("x")
	watch.Stop()
	println(watch.PrettyPrint())
	println(watch.GetTotalDuration().Milliseconds())
	println(watch.GetTaskCount())
}

func TestMissingLastStop(t *testing.T) {
	watch := NewStopWatchWithName("task")
	watch.Start("a")
	http.Get("http://jd.com")
	println(watch.PrettyPrint())
}

func TestFmtDuration(t *testing.T) {
	fmt.Println(fmtDuration(time.Nanosecond * 8938860000))
	fmt.Println(time.Millisecond * 9345)
}
