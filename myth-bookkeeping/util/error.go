package util

import "log"

func AssertNoError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
