package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/config"
	"github.com/wonderivan/logger"
)

func Server(_ []string) {
	appConfig := config.GetAppConfig()
	if !appConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.GET("/ping", HealthCheck)
	router.Static("/static", "./conf/static")
	router.StaticFile("/favicon.ico", "./conf/static/favicon.ico")

	router.GET("/mybook/category/typeList", ListCategoryType)
	router.GET("/mybook/category/list", ListCategory)
	router.GET("/mybook/account/list", ListAccount)
	router.POST("/mybook/record/create", CreateRecord)
	router.GET("/mybook/record/list", ListRecord)

	logger.Info("Open http://localhost:10006/static/")
	e := router.Run(":10006")
	logger.Error(e)
}
