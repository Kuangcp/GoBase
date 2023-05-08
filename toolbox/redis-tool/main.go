package main

import (
	"flag"
	"github.com/go-redis/redis"
)

var (
	debugFlag = false
	bigKey    = false
	queryKey  string

	fromAddr string
	fromPwd  string
	fromDB   int

	toAddr string
	toPwd  string
	toDB   int
)

func init() {
	flag.StringVar(&fromAddr, "addr", "127.0.0.1:6379", "origin redis address")
	flag.StringVar(&fromPwd, "pwd", "", "origin redis password")
	flag.IntVar(&fromDB, "db", 2, "origin redis db")

	flag.StringVar(&toAddr, "t.addr", "127.0.0.1:6379", "target redis address")
	flag.StringVar(&toPwd, "t.pwd", "", "target redis password")
	flag.IntVar(&toDB, "t.db", 3, "target redis db")

	flag.BoolVar(&bigKey, "bk", false, "scan big key")
	flag.StringVar(&queryKey, "key", "", "query key")
}

func main() {
	flag.Parse()
	originOpt := &redis.Options{
		Addr:     fromAddr,
		Password: fromPwd,
		DB:       fromDB,
	}

	if bigKey {
		scanBigKey(originOpt)
		return
	}

	if queryKey != "" {
		queryKeyDetail(originOpt)
		return
	}

	Action(SyncAllKey,
		originOpt,
		&redis.Options{
			Addr:     toAddr,
			Password: toPwd,
			DB:       toDB,
		}, false)
}
