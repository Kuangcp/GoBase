package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/cuibase"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
)

func Server(debugStatic bool, port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		ghelp.GinSuccessWith(c, "ok")
	})

	// 是否读取 statik 打包后的静态文件
	if debugStatic {
		router.Static("/static", "./static")
	} else {
		//// static file mapping
		//fileSystem, err := fs.New()
		//if err != nil {
		//	log.Fatal(err)
		//}
		//router.StaticFS("/static", fileSystem)
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

	ghelp.GracefulExit(srv)
}

func registerRouter(router *gin.Engine) {
	router.GET(buildPath("/listFile"), ListFile)
	router.POST(buildPath("/postFile"), CreateOrUpdateFile)
	router.GET(buildPath("/getFile"), FileContent)
	router.GET(buildPath("/switch"), SwitchFileState)
}

func buildPath(path string) string {
	return "/api/v1.0" + path
}
