package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	router := gin.Default()
	router.GET("/ping", HealthCheck)
	_ = router.Run()
}

func HealthCheck(c *gin.Context) {
	time.Sleep(time.Second * 1)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
