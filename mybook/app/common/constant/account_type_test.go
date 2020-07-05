package constant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetById(t *testing.T) {
	online := GetAccountTypeByIndex(AccountOnline)
	if online.Index != AccountOnline {
		assert.Fail(t, "type error")
	}
}
