package stopwatch

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestStopWatch_PrettyPrint(t *testing.T) {
	stopWatch := NewWithName("action1")
	stopWatch.Start("task1")
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
	stopWatch.Start("1")
	http.Get("http://jd.com")
	stopWatch.Stop()
	stopWatch.Start("2")
	fmt.Println(stopWatch.PrettyPrint())
}
