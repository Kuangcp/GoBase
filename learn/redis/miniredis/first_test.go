package miniredis

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

func TestSomething(t *testing.T) {
	// 相当于当前进程内开启了一个redis服务端, 并且同样通过redis客户端发起tcp连接去做交互
	s := miniredis.RunT(t)

	// Optionally set some keys your code expects:
	s.Set("foo", "bar")
	s.HSet("some", "other", "key")

	// Run your code and see if it behaves.
	// An example using the redigo library from "github.com/gomodule/redigo/redis":
	fmt.Println(s.Addr())
	c, err := redis.Dial("tcp", s.Addr())
	if err != nil {
		panic(err)
	}
	_, err = c.Do("SET", "foo", "bar")

	// Optionally check values in redis...
	if got, err := s.Get("foo"); err != nil || got != "bar" {
		t.Error("'foo' has the wrong value")
	}
	// ... or use a helper for that:
	s.CheckGet(t, "foo", "bar")

	// TTL and expiration:
	s.Set("foo", "bar")
	s.SetTTL("foo", 10*time.Second)
	s.FastForward(11 * time.Second)
	if s.Exists("foo") {
		t.Fatal("'foo' should not have existed anymore")
	}
	time.Sleep(time.Minute)
}
