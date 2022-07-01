package web

import (
	"embed"
	"net/http"

	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/logger"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
)

func Server(f embed.FS, debugStatic, notOpenPage bool, port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		ghelp.GinSuccessWith(c, "ok")
	})

	// 是否读取构建后静态目录
	if debugStatic {
		router.Static("/s", "./static")
	} else {
		resource := &ghelp.StaticResource{
			StaticFS: f,
			Path:     "static",
		}
		router.StaticFS("/s", http.FS(resource))
	}
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "s/")
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "s/favicon.ico")
	})

	// backend logic router
	registerRouter(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	url := "http://localhost" + srv.Addr
	logger.Info(url)
	if !notOpenPage {
		_ = ctk.OpenBrowser(url)
	}

	ghelp.GracefulExit(srv)
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
