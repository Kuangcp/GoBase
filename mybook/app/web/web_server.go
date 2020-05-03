package web

import (
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/config"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/wonderivan/logger"
	"strconv"
)

func Server(_ []string) {
	appConfig := config.GetAppConfig()
	if !appConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if appConfig.Path == config.DefaultPath {
		service.AutoMigrateAll()
	}

	router := gin.Default()

	router.GET("/ping", HealthCheck)
	router.Static("/static", "./conf/static")
	router.StaticFile("/favicon.ico", "./conf/static/favicon.ico")

	logicRouter(router)

	logger.Info("Open http://localhost:10006/static/")
	e := router.Run(":" + strconv.Itoa(appConfig.Port))
	logger.Error(e)
}

func logicRouter(router *gin.Engine) {
	api := "/mybook"
	router.GET(api+"/category/typeList", ListCategoryType)
	router.GET(api+"/category/list", ListCategory)

	router.GET(api+"/account/list", ListAccount)

	router.POST(api+"/record/create", CreateRecord)
	router.GET(api+"/record/list", ListRecord)

	router.GET(api+"/record/category", CategoryRecord)

	router.GET(api+"/record/categoryDetail", CategoryDetailRecord)
	router.GET(api+"/record/categoryWeekDetail", WeekCategoryDetailRecord)
	router.GET(api+"/record/categoryMonthDetail", MonthCategoryDetailRecord)
}
