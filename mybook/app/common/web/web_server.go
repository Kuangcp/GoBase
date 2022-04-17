package web

import (
	"embed"
	"fmt"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"log"
	"mybook/app/account"
	"mybook/app/category"
	"mybook/app/common"
	"mybook/app/common/config"
	"mybook/app/loan"
	"mybook/app/record"
	"mybook/app/report"
	"mybook/app/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/logger"
)

func Server(staticFS embed.FS) {
	//if config.AppConf.DBFilePath == config.DefaultDBPath {
	//	common.AutoMigrateAll()
	//}

	router := gin.New()
	router.GET("/ping", common.HealthCheck)

	registerModule(router)
	registerStaticFile(staticFS, router)
	registerServerApi(router)

	// start web server by specific port
	var finalPort = fmt.Sprintf(":%v", config.AppConf.Port)
	srv := &http.Server{
		Addr:    finalPort,
		Handler: router,
	}
	logger.Info("Start http://localhost" + finalPort)

	ghelp.GracefulExit(srv)
}

// 注册前端内容
func registerStaticFile(staticFS embed.FS, router *gin.Engine) {
	// 是否读取二进制内嵌静态文件
	if config.AppConf.DebugStatic {
		router.Static("/s", "./mybook-static/dist")
		router.StaticFile("/favicon.ico", "./conf/static/favicon.ico")
	} else {
		resource := &ghelp.StaticResource{
			StaticFS: staticFS,
			Path:     "mybook-static/dist",
		}
		router.StaticFS("/s", http.FS(resource))

		router.GET("/favicon.ico", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "s/favicon.ico")
		})
	}

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "s/")
	})
}

func registerModule(router *gin.Engine) {
	router.Use(supportCORS)
	router.Use(gin.Logger())

	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ghelp.GinFailedWithMsg(c, err)
			return
		}
		ghelp.GinFailed(c)
	}))
}

func supportCORS(c *gin.Context) {
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

// 服务端API注册
func registerServerApi(router *gin.Engine) {
	// 分类
	router.GET(buildApi("/category/listCategoryType"), common.ListCategoryType)
	router.GET(buildApi("/category/listCategory"), category.ListCategory)
	router.GET(buildApi("/category/listCategoryTree"), category.ListCategoryTree)

	// 账户
	router.GET(buildApi("/account/listAccount"), account.ListAccount)
	router.POST(buildApi("/account/createAccount"), account.CreateNewAccount)

	// 用户
	router.GET(buildApi("/user/listUser"), user.ListUser)

	// 账单
	router.GET(buildApi("/record/calBalance"), record.QueryAccountBalance)
	router.POST(buildApi("/record/createRecord"), record.CreateRecord)
	router.POST(buildApi("/loan/create"), loan.CreateLoan)

	router.GET(buildApi("/record/listRecord"), record.ListRecord)
	router.GET(buildApi("/record/category"), record.CategoryRecord)

	router.GET(buildApi("/record/categoryDetail"), record.CategoryDetailRecord)
	router.GET(buildApi("/record/categoryWeekDetail"), record.WeekCategoryDetailRecord)
	router.GET(buildApi("/record/categoryMonthDetail"), record.MonthCategoryDetailRecord)

	router.GET(buildApi("/report/categoryPeriod"), report.CategoryPeriodReport) // 各分类周期报表
	router.GET(buildApi("/report/balanceReport"), report.BalanceReport)         // 余额报表
}

func buildApi(path string) string {
	return config.DefaultUrlPath + path
}
