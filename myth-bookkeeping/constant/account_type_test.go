package constant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetById(t *testing.T) {
	online := GetById(ONLINE_TYPE)
	if online.TypeId != ONLINE_TYPE {
		assert.Fail(t, "type error")
	}
}
