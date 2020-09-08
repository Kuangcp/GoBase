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
	_ "github.com/kuangcp/gobase/keyboard-man/app/statik"
	"github.com/rakyll/statik/fs"
	"github.com/wonderivan/logger"
)

type (
	LineVO struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		Stack     string `json:"stack"`
		Data      []int  `json:"data"`
		AreaStyle string `json:"areaStyle"`
		Color     string `json:"color"`
	}

	QueryParam struct {
		Length    int
		Offset    int
		Top       int64
		ChartType string
	}
)

var colorSet = [...]string{
	"#c23531",
	"#2f4554",
	"#61a0a8",
	"#d48265",
	"#91c7ae",
	"#749f83",
	"#ca8622",
	"#bda29a",
	"#6e7074",
	"#546570",
	"#c4ccd3",
}

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

	// backend logic router
	registerRouter(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logger.Info("http://localhost" + srv.Addr)

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
	router.GET(buildPath("/hotKeyWithCount"), HotKeyWithNum)
	router.GET(buildPath("/recentDay"), RecentDay)
	router.GET(buildPath("/hotKeyName"), HotKey)
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
			result, err := GetConnection().ZScore(GetRankKeyByString(day), key).Result()
			if err != nil {
				result = 0
			}
			hitPreDay = append(hitPreDay, int(result))
		}

		keyCode, err := strconv.Atoi(key)
		cuibase.CheckIfError(err)
		lines = append(lines, LineVO{
			Type:      param.ChartType,
			Name:      nameMap[key],
			Data:      hitPreDay,
			Stack:     "all",
			AreaStyle: "{normal: {}}",
			Color:     colorSet[keyCode%len(colorSet)],
		})
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
		top = "2"
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

	topInt -= 1
	if topInt < 0 {
		topInt = 0
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
	for k := range keyCode {
		name, err := GetConnection().HGet(KeyMap, k).Result()
		cuibase.CheckIfError(err)
		result[k] = name
	}
	return result
}

func hotKey(dayList []string, top int64) map[string]bool {
	keyCodeMap := make(map[string]bool)
	for i := range dayList {
		result, err := GetConnection().ZRevRange(GetRankKeyByString(dayList[i]), 0, top).Result()
		if err != nil {
			logger.Warn("get hot key error", err)
			continue
		}

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
