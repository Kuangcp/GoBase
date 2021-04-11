package dal

import (
	"mybook/app/common/config"
	"reflect"
	"testing"
)

func TestGetConnectionConfig(t *testing.T) {
	tests := []struct {
		name string
		want *config.AppConfig
	}{
		{name: "", want: &config.AppConfig{DBFilePath: "/tmp/test.db", DriverName: "sqlite3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.InitAppConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDBConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
