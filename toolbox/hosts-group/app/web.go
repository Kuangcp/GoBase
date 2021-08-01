package app

import (
	"embed"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
	"github.com/zserge/lorca"
)

func OpenWebView(url string) {
	if "windows" == runtime.GOOS {
		ui, err := lorca.New(url, "", 1024, 768)
		if err != nil {
			log.Fatal(err)
		}
		defer ui.Close()
		<-ui.Done()
	} else {
		command := exec.Command("webview-cmd-wrapper",
			"-url", url, "-width", "1024", "-height", "768", "-title", "Hosts-Group")
		err := command.Start()
		if err != nil {
			logger.Fatal(err.Error())
		}
	}
}

func WebServer(f embed.FS, port string) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		ghelp.GinSuccessWith(c, "ok")
	})

	// 是否读取 embed 打包后的静态文件
	if DebugStatic {
		router.Static("/static", "./static")
	} else {
		resource := &ghelp.StaticResource{
			StaticFS: f,
			Path:     "static",
		}
		router.StaticFS("/static", http.FS(resource))
	}
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "static/")
	})

	// backend logic router
	registerRouter(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logger.Info("http://localhost" + srv.Addr)
	ghelp.GracefulExitWithHook(srv, func() {
		systray.Quit()
	})
}

func registerRouter(router *gin.Engine) {
	router.GET(buildPath("/queryMode"), QueryMode)

	router.GET(buildPath("/listFile"), ListFile)
	router.POST(buildPath("/postFile"), CreateOrUpdateFile)
	router.GET(buildPath("/getFile"), FileContent)
	router.GET(buildPath("/currentHosts"), CurrentHosts)
	router.GET(buildPath("/switch"), SwitchFileState)
}

func buildPath(path string) string {
	return "/api/v1.0" + path
}
