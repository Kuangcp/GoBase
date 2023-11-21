package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"testing"
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

func TestSingleToCluster(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.16.33:6379",
		Password: "zkredis",
		DB:       2,
	})
	result, err := client.HGetAll(context.Background(), "tg-fetch:trans-apply-map").Result()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)

	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"192.168.16.203:6589", "192.168.16.204:6589", "192.168.16.205:6589"}, //"192.168.16.203:6590", "192.168.16.204:6590", "192.168.16.205:6590",

		Password: "jszt20.v",
	})

	log.Println(cluster.Get(context.Background(), "ping").Result())

	SyncAllKey(client, cluster)
}
