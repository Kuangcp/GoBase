package core

import (
	"github.com/kuangcp/logger"
	"testing"
)

func TestSplit(t *testing.T) {

	var s = []string{"ss", "a", "a", "a", "a", "a", "b"}
	array := splitArray(s, 5)
	logger.Info(array)
}
