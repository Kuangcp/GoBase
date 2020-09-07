package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/rakyll/statik/fs"
	"github.com/wonderivan/logger"
)

var (
	ChartsType = "bar"
)

type (
	LineVO struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Stack string `json:"stack"`
		Data  []int  `json:"data"`
	}

	QueryParam struct {
		Length    int
		Offset    int
		Top       int64
		ChartType string
	}
)

func Server(debugStatic bool, port string) {
	router := gin.Default()
	//router.GET("/ping", common.HealthCheck)

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

	// backend logic router
	registerRouter(router)

	srv := &http.Server{
		Addr:    ":" + port,
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
	router.GET(buildPath("/hotKeyWithNum"), HotKeyWithNum)
	router.GET(buildPath("/recentDay"), RecentDay)
	router.GET(buildPath("/hotKey"), HotKey)
}

func buildPath(path string) string {
	return "/api/v1.0" + path
}

func HotKey(c *gin.Context) {
	param := parseParam(c)

	dayList := buildDayList(param.Length, param.Offset)

	hotKey := hotKey(dayList, param.Top)

	nameMap := keyNameMap(hotKey)
	var result []string
	for _, v := range nameMap {
		result = append(result, v)
	}
	sort.Strings(result)
	GinSuccessWith(c, result)
}

func RecentDay(c *gin.Context) {
	param := parseParam(c)

	var result []string
	dayList := buildDayList(param.Length, param.Offset)
	for _, day := range dayList {
		score, err := GetConnection().ZScore(TotalCount, day).Result()
		if err != nil {
			result = append(result, day+"#0")
		} else {
			result = append(result, day+"#"+strconv.Itoa(int(score)))
		}
	}

	GinSuccessWith(c, result)
}

func getKeys(m map[string]bool) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func HotKeyWithNum(c *gin.Context) {
	param := parseParam(c)

	dayList := buildDayList(param.Length, param.Offset)
	//logger.Info(dayList)

	hotKey := hotKey(dayList, param.Top)
	//logger.Info(hotKey)

	nameMap := keyNameMap(hotKey)
	sortHotKeys := getKeys(hotKey)
	sort.Strings(sortHotKeys)
	var lines []LineVO
	for _, key := range sortHotKeys {
		var hitPreDay []int
		for _, day := range dayList {
			result, err := GetConnection().ZScore(Prefix+day+":rank", key).Result()
			if err != nil {
				result = 0
			}
			hitPreDay = append(hitPreDay, int(result))
		}
		vo := LineVO{Type: param.ChartType, Stack: "all", Name: nameMap[key], Data: hitPreDay}
		lines = append(lines, vo)
	}
	//logger.Info(lines)

	GinSuccessWith(c, lines)
}

func parseParam(c *gin.Context) QueryParam {
	length := c.Query("length")
	offset := c.Query("offset")
	top := c.Query("top")
	chartType := c.Query("type")

	if length == "" {
		length = "7"
	}
	if offset == "" {
		offset = "0"
	}
	if top == "" {
		top = "3"
	}

	lengthInt, err := strconv.Atoi(length)
	cuibase.CheckIfError(err)
	offsetInt, err := strconv.Atoi(offset)
	cuibase.CheckIfError(err)
	topInt, err := strconv.ParseInt(top, 10, 64)
	cuibase.CheckIfError(err)

	if chartType == "" {
		chartType = "bar"
	}
	return QueryParam{
		Length:    lengthInt,
		Offset:    offsetInt,
		Top:       topInt,
		ChartType: chartType,
	}
}

func keyNameMap(keyCode map[string]bool) map[string]string {
	result := make(map[string]string)
	for k, _ := range keyCode {
		name, err := GetConnection().HGet(KeyMap, k).Result()
		cuibase.CheckIfError(err)
		result[k] = name
	}
	return result
}

func hotKey(dayList []string, top int64) map[string]bool {
	keyCodeMap := make(map[string]bool)
	for i := range dayList {
		result, err := GetConnection().ZRevRange(Prefix+dayList[i]+":rank", 0, top).Result()
		cuibase.CheckIfError(err)
		for _, s := range result {
			keyCodeMap[s] = true
		}
	}
	return keyCodeMap
}

func buildDayList(length int, offset int) []string {
	now := time.Now()

	var result []string
	start := now.AddDate(0, 0, -offset)
	for i := 0; i < length; i++ {
		day := start.AddDate(0, 0, i).Format("2006:01:02")
		result = append(result, day)
	}
	return result
}
