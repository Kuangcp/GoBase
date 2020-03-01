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

	router.GET("/mybook/category/typeList", ListCategoryType)
	router.GET("/mybook/category/list", ListCategory)
	router.GET("/mybook/account/list", ListAccount)
	router.POST("/mybook/record/create", CreateRecord)

	e := router.Run(":10006")
	logger.Error(e)
}
