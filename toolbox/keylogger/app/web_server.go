package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/kuangcp/gobase/keylogger/app/statik"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/pkg/ginhelper"
	"github.com/rakyll/statik/fs"
	"github.com/wonderivan/logger"
)

func Server(debugStatic bool, port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		ginhelper.GinSuccessWith(c, "ok")
	})

	// 是否读取 statik 打包后的静态文件
	if debugStatic {
		router.Static("/static", "./static")
	} else {
		// static file mapping
		fileSystem, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}
		router.StaticFS("/static", fileSystem)
	}
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/favicon.ico")
	})

	// backend logic router
	registerRouter(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	url := "http://localhost" + srv.Addr
	logger.Info(url)
	if !debugStatic {
		_ = cuibase.OpenBrowser(url)
	}

	ginhelper.GracefulExit(srv)
}

func registerRouter(router *gin.Engine) {
	router.GET(buildPath("/lineMap"), LineMap)
	router.GET(buildPath("/heatMap"), HeatMap)
	router.GET(buildPath("/weeksHeatMap"), MultipleHeatMap)
	router.GET(buildPath("/calendarMap"), CalendarMap)
}

func buildPath(path string) string {
	return "/api/v1.0" + path
}
