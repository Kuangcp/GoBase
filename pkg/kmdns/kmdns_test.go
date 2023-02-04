package kmdns

import (
	"fmt"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	obj := New("test.server", time.Second*7)
	obj.Server()
}

func TestClient(t *testing.T) {
	obj := New("test.server", time.Second*7)
	server := obj.ClientRequest()
	fmt.Println(server)
}
