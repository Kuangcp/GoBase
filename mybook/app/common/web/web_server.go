package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/mybook/app/common"
	"github.com/kuangcp/gobase/mybook/app/common/config"
	_ "github.com/kuangcp/gobase/mybook/app/common/statik"
	"github.com/kuangcp/gobase/mybook/app/controller"
	"github.com/kuangcp/gobase/mybook/app/service"
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
	if port == config.DefaultPort {
		finalPort = strconv.Itoa(appConfig.Port)
	} else {
		finalPort = strconv.Itoa(port)
	}

	srv := &http.Server{
		Addr:    ":" + finalPort,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Warn("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Warn("Server exiting")
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
