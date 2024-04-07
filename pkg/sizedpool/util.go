package sizedpool

import (
	"log"
)

func Go(act func()) {
	go func() {
		defer func() {
			//捕获抛出的panic
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
		act()
	}()
}
