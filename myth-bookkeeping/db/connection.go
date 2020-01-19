package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type (
	Connection struct {
		DB *sql.DB
	}
	ConnectionConfig struct {
		Path       string
		DriverName string
	}
)

var config = &ConnectionConfig{Path: "/tmp/test.db", DriverName: "sqlite3"}

func GetConnectionConfig() *ConnectionConfig {
	if config != nil {
		return config
	}
	config = &ConnectionConfig{Path: "/tmp/bookkeeping.db", DriverName: "sqlite3"}
	return config
}

func GetConnection() *Connection {
	return getConnectionWithConfig(GetConnectionConfig())
}

func getConnectionWithConfig(config *ConnectionConfig) *Connection {
	db, err := sql.Open(config.DriverName, config.Path)
	if err != nil {
		log.Fatal(err)
	}

	return &Connection{DB: db}
}

func (this *Connection) Close() {
	e := this.DB.Close()
	if e != nil {
		log.Fatal(e)
	}
}
