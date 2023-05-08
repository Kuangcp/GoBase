package main

import (
	"testing"

	"github.com/go-redis/redis"
)

func Test_syncAllKey(t *testing.T) {
	Action(SyncAllKey, &redis.Options{
		Addr:     "172.16.19.227:16379",
		Password: "",
		DB:       1,
	}, &redis.Options{
		Addr:     "172.16.19.227:16379",
		Password: "",
		DB:       2,
	}, false)
}

func Test_syncKeyRecord(t *testing.T) {
	Action(SyncKeyRecord, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       2,
	}, &redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       5,
	}, false)
}
