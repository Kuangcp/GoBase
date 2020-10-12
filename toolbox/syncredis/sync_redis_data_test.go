package main

import (
	"testing"

	"github.com/go-redis/redis"
)

func Test_syncAllKey(t *testing.T) {
	Action(false, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       1,
	}, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       2,
	}, SyncAllKey)
}

func Test_syncKeyRecord(t *testing.T) {
	Action(false, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       2,
	}, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       5,
	}, SyncKeyRecord)
}
