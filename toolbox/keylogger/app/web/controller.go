package web

import (
	"fmt"
	"github.com/kuangcp/gobase/pkg/ctool"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/gobase/pkg/stopwatch"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/store"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/pkg/ghelp"
	"github.com/kuangcp/logger"
)

var colorSet = [...]string{
	"rgb(46,199,201)",
	"rgb(182,162,222)",
	"rgb(90,177,239)",
	"rgb(255,185,128)",
	"rgb(216,122,128)",
	"rgb(141,152,179)",
	"rgb(229,207,13)",
	"rgb(151,181,82)",
	"rgb(149,112,109)",
	"rgb(220,105,170)",
	"rgb(7,162,164)",
	"rgb(154,127,209)",
	"rgb(88,141,213)",
	"rgb(245,153,78)",
	"rgb(192,80,80)",
	"rgb(89,103,140)",
	"rgb(201,171,0)",
	"rgb(126,176,10)",
	"rgb(111,85,83)",
	"rgb(193,64,137)",
}

const (
	TimeUnitDay  LineTimeUnit = "day"
	TimeUnitHour LineTimeUnit = "hour"
)

type (
	LineTimeUnit string

	LineChartVO struct {
		Lines    []*LineVO `json:"lines"`
		Days     []string  `json:"days"`
		KeyNames []string  `json:"keyNames"`
	}

	LineVO struct {
		Type  string `json:"type"`
		Name  string `json:"name"`
		Stack string `json:"stack"`
		Data  []int  `json:"data"`
		Color string `json:"color"`

		AreaStyle string  `json:"areaStyle"`
		Label     LabelVO `json:"label"`
	}
	LabelVO struct {
		Show     bool   `json:"show"`
		Position string `json:"position"`
	}

	HeatMapVO struct {
		Data  [168][3]int `json:"data"`
		Max   int         `json:"max"`
		Total int         `json:"total"`
		Start string      `json:"start"`
		End   string      `json:"end"`
	}

	CalendarHeatMapVO struct {
		Data             [][2]string `json:"data"`
		Type             string      `json:"type"`
		CoordinateSystem string      `json:"coordinateSystem"`
		CalendarIndex    int         `json:"calendarIndex"`
	}
	CalendarStyleVO struct {
		Range string `json:"range"`
	}

	CalendarResultVO struct {
		Maps   []CalendarHeatMapVO `json:"maps"`
		Styles []CalendarStyleVO   `json:"styles"`
		Max    int                 `json:"max"`
	}

	DayBO struct {
		Day     string
		WeekDay string
	}
	QueryParam struct {
		Length    int
		Offset    int
		Weeks     int
		Top       int64
		ChartType string
		ShowLabel bool
		HideZero  bool
		TimeUnit  LineTimeUnit
	}
)

var commonLabel = LabelVO{Show: false, Position: "insideRight"}

func valueOf(value string) (LineTimeUnit, error) {
	switch value {
	case string(TimeUnitDay):
		return TimeUnitDay, nil
	case string(TimeUnitHour):
		return TimeUnitHour, nil
	default:
		return "", fmt.Errorf("not support time unit")
	}
}

// SyncDetails 从Redis迁移入LevelDB（除当天外的数据）
func SyncDetails(c *gin.Context) {
	SyncAllDetails()
	ghelp.GinSuccess(c)
}

func ScheduleSyncAllDetails() {
	for range time.NewTicker(time.Hour * 24 * 5).C {
		logger.Info("sync details")
		SyncAllDetails()
	}
}

func SyncAllDetails() {
	conn := store.GetConnection()

	today := store.GetDetailKey(time.Now())
	var cursor uint64 = 0
	first := true
	for cursor != 0 || first {
		first = false
		keys, u, err := conn.Scan(cursor, store.Prefix+"*", 1000).Result()
		if err != nil {
			logger.Error(err)
		}
		for _, k := range keys {
			// 当天数据不全，不转移
			if today == k {
				logger.Info("ignore today", k)
				continue
			}
			if strings.HasSuffix(k, "detail") {
				logger.Info("Sync: ", k)
				syncDetail(k)
			}
		}

		cursor = u
	}
}

func syncDetail(key string) {
	list := store.QueryDetailByKey(key)
	if list == nil {
		logger.Warn("no data")
		return
	}

	text := ""
	for _, d := range list {
		text += d.ToString() + "\n"
	}
	db := store.GetDb()
	err := db.Put([]byte(key), []byte(text), nil)
	if err != nil {
		logger.Error(err)
		return
	}
	conn := store.GetConnection()
	conn.Del(key)
}

func ExportDetail(c *gin.Context) {
	store.ExportDetailToCsv(time.Now())
	ghelp.GinSuccess(c)
}

func CalendarMap(c *gin.Context) {
	conn := store.GetConnection()
	data, err := conn.ZRange(store.TotalCount, 0, -1).Result()
	ctk.CheckIfError(err)
	totalData, err := conn.ZRangeWithScores(store.TotalCount, 0, -1).Result()
	ctk.CheckIfError(err)
	sort.Strings(data)

	scoreMap := make(map[string]int)
	for _, ele := range totalData {
		scoreMap[ele.Member.(string)] = int(ele.Score)
	}
	max := 0

	yearListMap := make(map[string][]string)
	for _, day := range data {
		fields := strings.Split(day, ":")
		yearStr := fields[0]
		if val, ok := yearListMap[yearStr]; ok {
			val = append(val, day)
			yearListMap[yearStr] = val
		} else {
			var val []string
			val = append(val, day)
			yearListMap[yearStr] = val
		}
	}

	// sort year
	var yearList []string
	for k := range yearListMap {
		yearList = append(yearList, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(yearList)))

	chartIndex := -1
	var mapList []CalendarHeatMapVO
	var styleList []CalendarStyleVO
	for i := range yearList {
		year := yearList[i]

		result, tempMax := buildYear(yearListMap[year], scoreMap)

		chartIndex += 1
		mapList = append(mapList, CalendarHeatMapVO{
			Type:             "heatmap",
			CoordinateSystem: "calendar",
			CalendarIndex:    chartIndex,
			Data:             result,
		})
		styleList = append(styleList, CalendarStyleVO{Range: year})
		if tempMax > max {
			max = tempMax
		}
	}

	ghelp.GinSuccessWith(c, CalendarResultVO{Maps: mapList, Styles: styleList, Max: max})
}

func buildYear(data []string, scoreMap map[string]int) ([][2]string, int) {
	var result [][2]string
	maxScore := 0
	var lastTime *time.Time = nil
	for _, day := range data {
		var dayTime, err = time.Parse(store.DateFormat, day)
		ctk.CheckIfError(err)

		if lastTime == nil {
			// fill year start to dayTime
			emptyDay := fillEmptyDay(dayTime.AddDate(0, 0, -dayTime.YearDay()+1), dayTime)
			result = append(result, emptyDay...)
			lastTime = &dayTime
		} else {
			emptyDay := fillEmptyDay(lastTime.AddDate(0, 0, 1), dayTime)
			result = append(result, emptyDay...)
			lastTime = &dayTime
		}
		score := scoreMap[day]
		if score > maxScore {
			maxScore = score
		}

		result = append(result, [2]string{dayTime.Format(ctk.YYYY_MM_DD), strconv.Itoa(score)})
	}
	return result, maxScore
}

func fillEmptyDay(startDay time.Time, endDay time.Time) [][2]string {
	var result [][2]string
	var indexDay = startDay
	if startDay.Equal(endDay) {
		return nil
	}
	for !indexDay.Equal(endDay) {
		result = append(result, [2]string{indexDay.Format(ctk.YYYY_MM_DD), "0"})
		indexDay = indexDay.AddDate(0, 0, 1)
	}
	return result
}

func MultipleHeatMap(c *gin.Context) {
	param, err := parseParam(c)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	now := time.Now()
	weekday := now.Weekday()
	var weeksMap []*HeatMapVO

	var mutex = &sync.Mutex{}
	max := 0
	for i := 0; i < param.Weeks; i++ {
		offset := int(weekday) + (7 * i)
		mapVO := buildDataByDatePeriod(7, offset)

		mutex.Lock()
		if mapVO.Max > max {
			max = mapVO.Max
		}
		weeksMap = append(weeksMap, mapVO)
		mutex.Unlock()
	}

	for _, vo := range weeksMap {
		vo.Max = max
	}
	ghelp.GinSuccessWith(c, weeksMap)
}

// HeatMap 热力图
func HeatMap(c *gin.Context) {
	param, err := parseParam(c)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}
	mapVO := buildDataByDatePeriod(param.Length, param.Offset)
	ghelp.GinSuccessWith(c, mapVO)
}

func buildDataByDatePeriod(length int, offset int) *HeatMapVO {
	dayList := buildDayList(length, offset)

	// data: [weekday, hour, count], [weekday, hour, count]
	var result [168][3]int

	//TODO no mutex or sync.Map use read write lock
	var mutex = &sync.Mutex{}
	// weekday -> hour -> count
	totalMap := make(map[int]map[int]int)
	var latch sync.WaitGroup
	latch.Add(len(dayList))

	watch := stopwatch.NewWithName("")
	watch.Start(fmt.Sprint(len(dayList), "day"))
	for _, day := range dayList {
		var curDay = day
		go func() {
			defer latch.Done()

			readDetailToMap(curDay, mutex, totalMap)
		}()
	}
	latch.Wait()
	logger.Debug(watch.PrettyPrint())

	total := 0
	max := 0
	for weekday, v := range totalMap {
		chartIndex := 6 - weekday
		for hour, count := range v {
			//logger.Info(weekday, hour)
			if count > max {
				max = count
			}
			total += count
			result[(chartIndex*24)+hour] = [...]int{chartIndex, hour, count}
		}
	}

	strings.Replace(dayList[0], ":", "-", -1)
	return &HeatMapVO{
		Max:   max,
		Total: total,
		Data:  result,
		Start: strings.Replace(dayList[0], ":", "-", -1),
		End:   strings.Replace(dayList[len(dayList)-1], ":", "-", -1),
	}
}

func readDetailToMap(curDay string, mutex *sync.Mutex, totalMap map[int]map[int]int) {

	list := store.QueryDetailByDay(curDay)
	for _, d := range list {
		curStrokeTime := time.Unix(d.HitTime/1000_000, 0)
		weekDay := int(curStrokeTime.Weekday())

		actionRoundLock(mutex, func() {
			dayMap := totalMap[weekDay]
			//curStr := cur.Format(DateFormat)
			//if curStr != curDay {
			//	logger.Error("error detail data", curStr, curDay)
			//}
			if dayMap == nil {
				dayMap = make(map[int]int)
				totalMap[weekDay] = dayMap
			}
			dayMap[curStrokeTime.Hour()] += 1
		})
	}
}

func actionRoundLock(mutex *sync.Mutex, action func()) {
	mutex.Lock()
	defer mutex.Unlock()
	action()
}
func HourLineChart(c *gin.Context, param *QueryParam) {

}

// LineMap 折线图 柱状图
func LineMap(c *gin.Context) {
	param, err := parseParam(c)
	if err != nil {
		ghelp.GinFailedWithMsg(c, err.Error())
		return
	}

	if param.TimeUnit == TimeUnitHour {
		HourLineChart(c, param)
		return
	}

	dayList := buildDayList(param.Length, param.Offset)
	hotKey := hotKey(dayList, param.Top)
	if len(hotKey) == 0 {
		ghelp.GinFailed(c)
		return
	}

	nameMap := keyNameMap(hotKey)

	// keyNames
	var keyNames []string
	for _, v := range nameMap {
		keyNames = append(keyNames, v)
	}
	sort.Strings(keyNames)
	if len(keyNames) == 0 {
		ghelp.GinFailed(c)
		return
	}

	// days
	var days []string
	showDayList := buildDayWithWeekdayList(param.Length, param.Offset)
	conn := store.GetConnection()
	for _, day := range showDayList {
		score, err := conn.ZScore(store.TotalCount, day.Day).Result()
		if err != nil {
			score = 0
		}
		kpm, kpmErr := conn.Get(store.GetTodayMaxKPMKeyByString(day.Day)).Result()
		if kpmErr != nil {
			kpm = "0"
		}

		dayStr := fmt.Sprintf("%s|%s|%d|%s", strings.Replace(day.Day, ":", "-", 2), day.WeekDay, int(score), kpm)
		days = append(days, dayStr)
	}
	if len(days) == 0 {
		ghelp.GinFailed(c)
		return
	}

	// lines
	sortHotKeys := getMapKeys(hotKey)
	sort.Strings(sortHotKeys)
	var lines []*LineVO
	commonLabel.Show = param.ShowLabel
	for _, key := range sortHotKeys {
		var hitPreDay []int
		for _, day := range dayList {
			result, err := conn.ZScore(store.GetRankKeyByString(day), key).Result()
			if err != nil {
				result = 0
			}
			hitPreDay = append(hitPreDay, int(result))
		}

		keyCode, err := strconv.Atoi(key)
		ctk.CheckIfError(err)
		lines = append(lines, &LineVO{
			Type:      param.ChartType,
			Name:      nameMap[key],
			Data:      hitPreDay,
			Stack:     "all",
			AreaStyle: "{normal: {}}",
			Label:     commonLabel,
			Color:     colorSet[keyCode%len(colorSet)],
		})
	}
	//logger.Info(lines)
	if param.HideZero {
		zeroDay := ctool.NewSet[int]()
		var newDays []string
		for i := range days {
			dayStr := days[i]
			if isIdleDay(dayStr, 8000) {
				//if isZeroDay(dayStr) {
				zeroDay.Add(i)
			} else {
				newDays = append(newDays, dayStr)
			}
		}
		if zeroDay.IsEmpty() {
			ghelp.GinSuccessWith(c, LineChartVO{Lines: lines, Days: days, KeyNames: keyNames})
			return
		}

		for i := range lines {
			line := lines[i]
			var newData []int
			for idx := range line.Data {
				if zeroDay.Contains(idx) {
					continue
				}
				newData = append(newData, line.Data[idx])
			}
			line.Data = newData
		}

		ghelp.GinSuccessWith(c, LineChartVO{Lines: lines, Days: newDays, KeyNames: keyNames})
	} else {
		ghelp.GinSuccessWith(c, LineChartVO{Lines: lines, Days: days, KeyNames: keyNames})
	}
}

func isZeroDay(day string) bool {
	return strings.HasSuffix(day, "0|0")
}
func isIdleDay(day string, low int) bool {
	arr := strings.Split(day, "|")
	total := arr[len(arr)-2]
	atoi, err := strconv.Atoi(total)
	if err != nil {
		return false
	}
	return atoi < low
}

func getMapKeys(m map[string]bool) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func parseParam(c *gin.Context) (*QueryParam, error) {
	length := c.DefaultQuery("length", "7")
	offset := c.DefaultQuery("offset", "0")
	top := c.DefaultQuery("top", "2")
	chartType := c.DefaultQuery("type", "bar")
	showLabel := c.DefaultQuery("showLabel", "false")
	weeks := c.DefaultQuery("weeks", "1")
	timeUnit := c.DefaultQuery("timeUnit", "day")
	hideZero := c.DefaultQuery("hideZero", "false")

	unit, err := valueOf(timeUnit)
	if err != nil {
		return nil, err
	}
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		return nil, err
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, err
	}
	topInt, err := strconv.ParseInt(top, 10, 64)
	if err != nil {
		return nil, err
	}
	showLabelBool, err := strconv.ParseBool(showLabel)
	if err != nil {
		return nil, err
	}
	hideZeroBool, err := strconv.ParseBool(hideZero)
	if err != nil {
		return nil, err
	}

	weeksInt, err := strconv.Atoi(weeks)
	if err != nil {
		return nil, err
	}

	topInt -= 1
	if topInt < 0 {
		topInt = 0
	}
	if lengthInt <= 0 {
		lengthInt = 1
	}
	return &QueryParam{
		Length:    lengthInt,
		Offset:    offsetInt,
		Top:       topInt,
		Weeks:     weeksInt,
		ChartType: chartType,
		ShowLabel: showLabelBool,
		TimeUnit:  unit,
		HideZero:  hideZeroBool,
	}, nil
}

func keyNameMap(keyCode map[string]bool) map[string]string {
	result := make(map[string]string)
	for k := range keyCode {
		name, err := store.GetConnection().HGet(store.KeyMap, k).Result()
		if err != nil {
			result[k] = "unknown"
		}
		result[k] = name
	}
	return result
}

func hotKey(dayList []string, top int64) map[string]bool {
	//start := time.Now().UnixNano()
	keyCodeMap := make(map[string]bool)
	for i := range dayList {
		result, err := store.GetConnection().ZRevRange(store.GetRankKeyByString(dayList[i]), 0, top).Result()
		if err != nil {
			logger.Warn("get hot key error", err)
			continue
		}

		for _, s := range result {
			keyCodeMap[s] = true
		}
	}
	//end := time.Now().UnixNano()
	//logger.Info("hotKey: ", end-start)
	return keyCodeMap
}

// 不超过今天
func buildDayList(length int, offset int) []string {
	now := time.Now()

	var result []string
	start := now.AddDate(0, 0, -offset)
	for i := 0; i < length; i++ {
		cursor := start.AddDate(0, 0, i)
		day := cursor.Format(store.DateFormat)
		if cursor.After(now) {
			return result
		}
		result = append(result, day)
	}
	return result
}

func buildDayWithWeekdayList(length int, offset int) []DayBO {
	now := time.Now()

	var result []DayBO
	start := now.AddDate(0, 0, -offset)
	for i := 0; i < length; i++ {
		tempTime := start.AddDate(0, 0, i)
		day := tempTime.Format(store.DateFormat)
		if tempTime.After(now) {
			return result
		}
		result = append(result, DayBO{Day: day, WeekDay: buildWeekDay(tempTime.Weekday())})
	}
	return result
}

// 周
func buildWeekDay(weekday time.Weekday) string {
	switch weekday {
	case time.Monday:
		return "一"
	case time.Tuesday:
		return "二"
	case time.Wednesday:
		return "三"
	case time.Thursday:
		return "四"
	case time.Friday:
		return "五"
	case time.Saturday:
		return "六"
	case time.Sunday:
		return "七"
	}
	return ""
}
