package service

import (
	"github.com/kuangcp/gobase/mybook/domain"
	"testing"
)

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		go ListAccounts()
	}
}

func TestAddCount(t *testing.T) {
	type args struct {
		account *domain.Account
	}
	account := &domain.Account{Name: "test", InitAmount: 0, TypeId: 1}
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