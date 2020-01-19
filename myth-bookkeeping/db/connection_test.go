package db

import (
	"reflect"
	"testing"
)

func TestGetConnectionConfig(t *testing.T) {
	tests := []struct {
		name string
		want *ConnectionConfig
	}{
		{name: "", want: &ConnectionConfig{Path: "/tmp/test.db", DriverName: "sqlite3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConnectionConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConnectionConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitDB(t *testing.T) {
	connection := GetConnection()
	defer connection.Close()

	connection.DB.Exec("create table account ")
}
