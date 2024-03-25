package common

import (
	"fmt"
	"sync"
	"testing"
)

func TestSyncRun(t *testing.T) {
	lock := &sync.Mutex{}
	SyncRun(lock, func() error {
		fmt.Println("1")
		x := 0
		fmt.Println(1 / x)
		return nil
	})
	SyncRun(lock, func() error {
		fmt.Println("2")
		return nil
	})
}

func TestSyncRuns(t *testing.T) {
	lock := &sync.Mutex{}
	SyncRuns(lock, func() {
		fmt.Println("1")
		x := 0
		fmt.Println(1 / x)
	})
	SyncRuns(lock, func() {
		fmt.Println("2")
	})
}
