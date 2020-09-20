package main

import (
	"github.com/go-redis/redis"
	"testing"
)

func Test_syncAllKey(t *testing.T) {
	Action(&redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       0,
	}, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       1,
	}, SyncAllKey)
}

func Test_syncKeyRecord(t *testing.T) {
	Action(&redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       2,
	}, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       5,
	}, SyncKeyRecord)
}
