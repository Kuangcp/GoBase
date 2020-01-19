package service

import (
	"log"
	"reflect"
	"testing"

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
