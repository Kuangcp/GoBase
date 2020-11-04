package web

import (
	"log"
	"mybook/app/common"
	"mybook/app/common/config"
	_ "mybook/app/common/statik"
	"mybook/app/controller"
	"mybook/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ginhelper"
	"github.com/rakyll/statik/fs"
	"github.com/wonderivan/logger"
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
	if port == 0 {
		finalPort = strconv.Itoa(appConfig.Port)
	} else {
		finalPort = strconv.Itoa(port)
	}

	srv := &http.Server{
		Addr:    ":" + finalPort,
		Handler: router,
	}

	logger.Info("Start http://localhost:" + finalPort)
	ginhelper.GracefulExit(srv)
}

func registerRouter(router *gin.Engine) {
	// 分类
	router.GET(buildApi("/category/typeList"), common.ListCategoryType)
	router.GET(buildApi("/category/list"), common.ListCategory)

	// 账户
	router.GET(buildApi("/account/list"), controller.ListAccount)
	router.GET(buildApi("/account/balance"), controller.CalculateAccountBalance)

	// 账单
	router.POST(buildApi("/record/create"), controller.CreateRecord)
	router.GET(buildApi("/record/list"), controller.ListRecord)

	router.GET(buildApi("/record/category"), controller.CategoryRecord)

	router.GET(buildApi("/record/categoryDetail"), controller.CategoryDetailRecord)
	router.GET(buildApi("/record/categoryWeekDetail"), controller.WeekCategoryDetailRecord)
	router.GET(buildApi("/record/categoryMonthDetail"), controller.MonthCategoryDetailRecord)
}

func buildApi(path string) string {
	return config.DefaultUrlPath + path
}
