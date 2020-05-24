package util

import (
	"github.com/wonderivan/logger"
	"log"
)

func AssertNoError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func RecordError(e error) {
	if e != nil {
		logger.Error(e)
	}
}
