package raft

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	raft1 := CreateRaft(44002, []string{"192.168.9.155:44003", "192.168.9.155:44004"})
	fmt.Println(raft1.Stat())

	raft2 := CreateRaft(44003, []string{"192.168.9.155:44002", "192.168.9.155:44004"})
	fmt.Println(raft2.Stat())

	raft3 := CreateRaft(44004, []string{"192.168.9.155:44002", "192.168.9.155:44003"})
	fmt.Println(raft3.Stat())

	go raft1.Start()
	go raft2.Start()
	raft3.Start()
}
