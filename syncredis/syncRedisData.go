package main

import (
	"log"

	"github.com/go-redis/redis"
)

func main() {
	origin := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "",
		DB:       1,
	})

	target := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6666",
		Password: "",
		DB:       2,
	})

	_, err := origin.Ping().Result()
	if err != nil {
		log.Fatal("origin can not connetction ", err)
	}

	_, err = target.Ping().Result()
	if err != nil {
		log.Fatal("target can not connetction ", err)
	}

	log.Println("start sync")
	result, _ := origin.Keys("*").Result()
	log.Println("total key: ", result)

	for i := range result {
		key := result[i]
		keyType, _ := origin.Type(key).Result()
		log.Println(key, keyType)
		switch keyType {
		case STRING:
			val, _ := origin.Get(key).Result()
			log.Println("value: ", val)
		}
	}
}
