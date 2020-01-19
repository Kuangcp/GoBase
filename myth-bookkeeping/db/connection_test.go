package db

import (
	"io/ioutil"
	"log"
	"os"
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

func TestInitTableStructure(t *testing.T) {
	connection := getConnectionWithConfig(&ConnectionConfig{Path: "/tmp/tests.db", DriverName: "sqlite3"})
	defer connection.Close()

	file, e := os.Open("../resources/db.sql")
	if e != nil {
		log.Fatal(e)
	}
	bytes, e := ioutil.ReadAll(file)
	sql := string(bytes)
	result, e := connection.DB.Exec(sql)
	if e != nil {
		log.Fatal(e)
	}
	log.Println(result)
}
