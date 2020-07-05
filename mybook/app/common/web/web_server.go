package web

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/common"
	"github.com/kuangcp/gobase/mybook/app/common/config"
	record "github.com/kuangcp/gobase/mybook/app/controller"
	"github.com/kuangcp/gobase/mybook/app/service"
	"github.com/rakyll/statik/fs"
	"github.com/wonderivan/logger"

	_ "github.com/kuangcp/gobase/mybook/app/common/statik"
)

func Server(debug bool) {
	appConfig := config.GetAppConfig()
	if !appConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if appConfig.Path == config.DefaultPath {
		service.AutoMigrateAll()
	}

	router := gin.Default()
	router.GET("/ping", common.HealthCheck)

	// static file not package in
	if debug {
		router.Static("/static", "./conf/static")
		router.StaticFile("/favicon.ico", "./conf/static/favicon.ico")
	} else {
		// static file mapping
		fileSystem, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}
		router.StaticFS("/static", fileSystem)
		router.GET("/favicon.ico", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "static/favicon.ico")
		})
	}

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/")
	})

	// backend logic router
	backendRouter(router)

	port := strconv.Itoa(appConfig.Port)
	logger.Info("Open http://localhost:" + port)
	e := router.Run(":" + port)
	logger.Error(e)
}

func backendRouter(router *gin.Engine) {
	api := "/api"
	router.GET(api+"/category/typeList", common.ListCategoryType)
	router.GET(api+"/category/list", common.ListCategory)

	router.GET(api+"/account/list", record.ListAccount)
	router.GET(api+"/account/balance", record.AccountBalance)

	router.POST(api+"/record/create", record.CreateRecord)
	router.GET(api+"/record/list", record.ListRecord)

	router.GET(api+"/record/category", record.CategoryRecord)

	router.GET(api+"/record/categoryDetail", record.CategoryDetailRecord)
	router.GET(api+"/record/categoryWeekDetail", record.WeekCategoryDetailRecord)
	router.GET(api+"/record/categoryMonthDetail", record.MonthCategoryDetailRecord)
}
