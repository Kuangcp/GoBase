package util

import (
	"github.com/kuangcp/logger"
)

func AssertNoError(e error) {
	if e != nil {
		logger.Emer(e)
		panic(e.Error())
	}
}

func RecordError(e error) {
	if e != nil {
		logger.Error(e)
	}
}
