package constant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetById(t *testing.T) {
	online := GetAccountTypeByIndex(ACCOUNT_ONLINE)
	if online.Index != ACCOUNT_ONLINE {
		assert.Fail(t, "type error")
	}
}
