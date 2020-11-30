package web

import (
	"github.com/kuangcp/gobase/pkg/ghelp"
	"log"
	"mybook/app/common"
	"mybook/app/common/config"
	_ "mybook/app/common/statik"
	"mybook/app/controller"
	"mybook/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		router.Static("/static", "./mybook-static/dist")
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

	router.Use(Cors())
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
	ghelp.GracefulExit(srv)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}

func registerRouter(router *gin.Engine) {
	// 分类
	router.GET(buildApi("/category/typeList"), common.ListCategoryType)
	router.GET(buildApi("/category/list"), common.ListCategory)
	router.GET(buildApi("/category/listTree"), common.ListCategoryTree)

	// 账户
	router.GET(buildApi("/account/list"), controller.ListAccount)
	router.GET(buildApi("/account/balance"), controller.CalculateAccountBalance)

	// 账单
	router.POST(buildApi("/record/createRecord"), controller.CreateRecord)
	router.GET(buildApi("/record/list"), controller.ListRecord)

	router.GET(buildApi("/record/category"), controller.CategoryRecord)

	router.GET(buildApi("/record/categoryDetail"), controller.CategoryDetailRecord)
	router.GET(buildApi("/record/categoryWeekDetail"), controller.WeekCategoryDetailRecord)
	router.GET(buildApi("/record/categoryMonthDetail"), controller.MonthCategoryDetailRecord)

	router.GET(buildApi("/report/categoryMonth"), controller.CategoryMonthMap)

}

func buildApi(path string) string {
	return config.DefaultUrlPath + path
}
