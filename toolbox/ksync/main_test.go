package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestReadDir(t *testing.T) {
	open, err := os.Open("/home/kcp/test/ss/c/")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info(open)

	dir, err := ioutil.ReadDir("/home/kcp/test/ss/c/")
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Info(dir)
}

func TestNotify(t *testing.T) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("/home/zk/ksync")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	//<-make(chan struct{})

	time.Sleep(time.Hour)
}
