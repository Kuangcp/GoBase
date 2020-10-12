package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/kuangcp/gobase/keylogger/app/statik"
	"github.com/rakyll/statik/fs"
	"github.com/wonderivan/logger"
)

func Server(debugStatic bool, port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		GinSuccessWith(c, "ok")
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

	logger.Info("http://localhost" + srv.Addr)

	gracefulExit(srv)
}

func gracefulExit(srv *http.Server) {
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %s\n", err)
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
	router.GET(buildPath("/lineMap"), LineMap)
	router.GET(buildPath("/heatMap"), HeatMap)
	router.GET(buildPath("/weeksHeatMap"), MultipleHeatMap)
	router.GET(buildPath("/calendarMap"), CalendarMap)
}

func buildPath(path string) string {
	return "/api/v1.0" + path
}