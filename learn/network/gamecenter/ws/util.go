package ws

import "sync"

func ActionRoundLock(lock *sync.Mutex, action func() error) error {
	lock.Lock()
	err := action()
	lock.Unlock()
	return err
}

func SilentActionRoundLock(lock *sync.Mutex, action func()) {
	lock.Lock()
	action()
	lock.Unlock()
}
