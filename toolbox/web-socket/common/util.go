package common

import (
	"github.com/kuangcp/logger"
	"sync"
)

func SyncRun(lock *sync.Mutex, action func() error) error {
	defer func() {
		a := recover()
		if a != nil {
			logger.Error(a)
		}
		lock.Unlock()
	}()
	lock.Lock()
	err := action()
	return err
}

func SyncRuns(lock *sync.Mutex, action func()) {
	_ = SyncRun(lock, func() error {
		action()
		return nil
	})
}
