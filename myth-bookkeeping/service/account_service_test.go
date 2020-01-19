package service

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/kuangcp/gobase/myth-bookkeeping/domain"
)

func TestQueryAllAccounts(t *testing.T) {
	tests := []struct {
		name string
		want []domain.Account
	}{
		{name: "name", want: QueryAllAccounts()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QueryAllAccounts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAllAccounts() = %v, want %v", got, tt.want)
			} else {
				log.Println(got)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type args struct {
		account *domain.Account
	}
	account := &domain.Account{Name: "name", InitAmount: 0, Type: 1, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix(), DeletedAt: 0}
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
