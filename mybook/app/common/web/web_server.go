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

func Server(debugStatic bool, port int) {
	appConfig := config.GetAppConfig()
	if !appConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if appConfig.Path == config.DefaultDBPath {
		service.AutoMigrateAll()
	}

	router := gin.Default()
	router.GET("/ping", common.HealthCheck)

	// 是否读取 statik 打包后的静态文件
	if debugStatic {
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
	registerRouter(router)

	// start web server by specific port
	var finalPort string
	if port == config.DefaultPort {
		finalPort = strconv.Itoa(appConfig.Port)
	} else {
		finalPort = strconv.Itoa(port)
	}
	logger.Info("Open http://localhost:" + finalPort)
	e := router.Run(":" + finalPort)
	logger.Error(e)
}

func registerRouter(router *gin.Engine) {
	router.GET(buildApi("/category/typeList"), common.ListCategoryType)
	router.GET(buildApi("/category/list"), common.ListCategory)

	router.GET(buildApi("/account/list"), record.ListAccount)
	router.GET(buildApi("/account/balance"), record.CalculateAccountBalance)

	router.POST(buildApi("/record/create"), record.CreateRecord)
	router.GET(buildApi("/record/list"), record.ListRecord)

	router.GET(buildApi("/record/category"), record.CategoryRecord)

	router.GET(buildApi("/record/categoryDetail"), record.CategoryDetailRecord)
	router.GET(buildApi("/record/categoryWeekDetail"), record.WeekCategoryDetailRecord)
	router.GET(buildApi("/record/categoryMonthDetail"), record.MonthCategoryDetailRecord)
}

func buildApi(path string) string {
	return config.DefaultUrlPath + path
}
