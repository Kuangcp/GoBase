package service

import (
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
	"testing"
)

func TestQueryAllAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
	    go QueryAll()
	}
}

func TestInsert(t *testing.T) {
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
			Insert(tt.args.account)
		})
	}
}
