package main

import (
	"github.com/go-redis/redis"
	"github.com/kuangcp/logger"
	"math/rand"
	"testing"
	"time"
)

const (
	Lock  = "global-lock"
	Queue = "queue"
	Doing = "doing"
)

func TestTryLock(t *testing.T) {
	option := redis.Options{Addr: "192.168.9.155:6667", Password: "", DB: 0}
	InitConnection(option, true)

	// 全局锁, 注册消费任务, 如果任务完成就空出额度
	// 异常处理: 任务持续挂起超长时间, 任务消费过程异常中断线程, 任务完成了但是没空出额度,

	go func() {
		for t := range time.NewTicker(time.Second * 1).C {
			AddTask(t.Format("2006-01-02 15:04:05.000"))
		}
	}()

	time.Sleep(time.Second * 6)

	go func() {
		for range time.NewTicker(time.Second * 2).C {
			HandleTask()
		}
	}()

	time.Sleep(time.Minute * 10)
}

func AddTask(id string) {
	conn := GetConnection()
	result, _ := conn.SetNX(Lock, "", time.Second*1).Result()
	if !result {
		logger.Info("try lock failed")
		return
	}

	defer func() {
		conn.Del(Lock)
	}()

	logger.Info("add", id)
	conn.LPush(Queue, id)
}

func popId() string {
	conn := GetConnection()
	result, _ := conn.SetNX(Lock, "", time.Second*1).Result()
	if !result {
		logger.Info("try lock failed")
		return ""
	}

	defer func() {
		conn.Del(Lock)
	}()

	id, err := conn.RPopLPush(Queue, Doing).Result()
	if err != nil {
		logger.Error(err)
		return ""
	}
	if id == "" {
		return ""
	}

	return id
}

func HandleTask() {
	id := popId()
	if id == "" {
		return
	}
	waste := 1000 + rand.Intn(6000)
	logger.Info("start", id, waste)
	time.Sleep(time.Millisecond * time.Duration(waste))
	conn := GetConnection()
	conn.LRem(Doing, 1, id)
	logger.Info("delete doing", id)
}
