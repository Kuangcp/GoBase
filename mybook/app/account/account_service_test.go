package account

import (
	"testing"
)

func TestAddCount(t *testing.T) {
	type args struct {
		account *Account
	}
	account := &Account{Name: "test", InitAmount: 0, TypeId: 1}
	tests := []struct {
		name string
		args args
	}{
		{name: "", args: args{account: account}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddAccount(tt.args.account)
		})
	}
}
