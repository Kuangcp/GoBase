package data_source

import (
	"reflect"
	"testing"
)

func TestGetConnectionConfig(t *testing.T) {
	tests := []struct {
		name string
		want *ConnectionConfig
	}{
		{name: "", want: &ConnectionConfig{Path: "/tmp/test.data_source", DriverName: "sqlite3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDBConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDBConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
