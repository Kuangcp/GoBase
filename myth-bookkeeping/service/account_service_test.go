package service

import (
	"github.com/kuangcp/gobase/myth-bookkeeping/constant"
	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
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

func TestInitAccount(t *testing.T) {
	AddAccount(&domain.Account{TypeId: constant.CASH_TYPE, Name: "现金", InitAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.CREDIT_TYPE, Name: "招商信用卡", InitAmount: 0, MaxAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.CREDIT_TYPE, Name: "农行信用卡", InitAmount: 0, MaxAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.CREDIT_TYPE, Name: "花呗", InitAmount: 0, MaxAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.ONLINE_TYPE, Name: "支付宝", InitAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.ONLINE_TYPE, Name: "微信", InitAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.DEPOSIT_TYPE, Name: "储蓄卡", InitAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.FINANCE_TYPE, Name: "招商基金", InitAmount: 0})
	AddAccount(&domain.Account{TypeId: constant.FINANCE_TYPE, Name: "支付宝基金", InitAmount: 0})
}
