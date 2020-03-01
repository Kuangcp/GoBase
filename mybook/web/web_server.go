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

	router.GET("/mybook/record/typeList", ListRecordType)
	router.GET("/mybook/record/accountList", ListAccount)
	router.GET("/mybook/record/categoryList", ListCategory)
	router.POST("/mybook/record/create", CreateRecord)

	e := router.Run(":10006")
	logger.Error(e)
}
