package dal

import (
	"reflect"
	"testing"

	"github.com/kuangcp/gobase/mybook/conf"
)

func TestGetConnectionConfig(t *testing.T) {
	tests := []struct {
		name string
		want *conf.ConnectionConfig
	}{
		{name: "", want: &conf.ConnectionConfig{Path: "/tmp/test.db", DriverName: "sqlite3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDBConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDBConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}