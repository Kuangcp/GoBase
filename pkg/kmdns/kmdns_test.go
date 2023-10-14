package kmdns

import (
	"fmt"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	obj := New(time.Second*7, "test.server", "abc.server")
	obj.Server()
}

func TestClient(t *testing.T) {
	obj := New(time.Second*7, "test.server", "abc.server")
	server := obj.ClientRequest()
	fmt.Println(server)
}
