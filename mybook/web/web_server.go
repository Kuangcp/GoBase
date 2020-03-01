package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/conf"
	"github.com/wonderivan/logger"
)

func Server(_ []string) {
	conf.LoadConfig()
	router := gin.Default()

	router.GET("/ping", HealthCheck)
	router.Static("/static", "./web/static")

	router.POST("/record/create", CreateRecord)
	router.GET("/record/typeList", ListRecordType)
	router.GET("/record/categoryList", ListCategory)

	e := router.Run(":10006")
	logger.Error(e)
}
