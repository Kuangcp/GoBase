package app

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kuangcp/gobase/cuibase"
	"github.com/wonderivan/logger"
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

type (
	LineChartVO struct {
		Lines    []LineVO `json:"lines"`
		Days     []string `json:"days"`
		KeyNames []string `json:"keyNames"`
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
	}
)

var commonLabel = LabelVO{Show: false, Position: "insideRight"}

func CalendarMap(c *gin.Context) {
	conn := GetConnection()
	data, err := conn.ZRange(TotalCount, 0, -1).Result()
	cuibase.CheckIfError(err)
	totalData, err := conn.ZRangeWithScores(TotalCount, 0, -1).Result()
	cuibase.CheckIfError(err)
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
	sort.Strings(yearList)

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

	GinSuccessWith(c, CalendarResultVO{Maps: mapList, Styles: styleList, Max: max})
}

func buildYear(data []string, scoreMap map[string]int) ([][2]string, int) {
	var result [][2]string
	max := 0
	var lastTime *time.Time = nil
	for _, day := range data {
		var dayTime, err = time.Parse("2006:01:02", day)
		cuibase.CheckIfError(err)

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
		if score > max {
			max = score
		}

		result = append(result, [2]string{dayTime.Format("2006-01-02"), strconv.Itoa(score)})
	}
	return result, max
}

func fillEmptyDay(startDay time.Time, endDay time.Time) [][2]string {
	var result [][2]string
	var indexDay = startDay
	if startDay.Equal(endDay) {
		return nil
	}
	for !indexDay.Equal(endDay) {
		result = append(result, [2]string{indexDay.Format("2006-01-02"), "0"})
		indexDay = indexDay.AddDate(0, 0, 1)
	}
	return result
}

func MultipleHeatMap(c *gin.Context) {
	param, err := parseParam(c)
	if err != nil {
		GinFailedWithMsg(c, err.Error())
		return
	}

	now := time.Now()
	weekday := now.Weekday()
	var weeksMap []*HeatMapVO

	var mutex = &sync.Mutex{}
	var latch sync.WaitGroup
	latch.Add(param.Weeks)
	max := 0
	for i := 0; i < param.Weeks; i++ {
		offset := int(weekday) + (7 * i)
		mapVO := buildDataFromFrame(7, offset)
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
	GinSuccessWith(c, weeksMap)
}

//HeatMap 热力图
func HeatMap(c *gin.Context) {
	param, err := parseParam(c)
	if err != nil {
		GinFailedWithMsg(c, err.Error())
		return
	}
	mapVO := buildDataFromFrame(param.Length, param.Offset)
	GinSuccessWith(c, mapVO)
}

func buildDataFromFrame(length int, offset int) *HeatMapVO {
	dayList := buildDayList(length, offset)

	// [weekday, hour, count], [weekday, hour, count]
	var result [168][3]int

	var mutex = &sync.Mutex{}
	// weekday -> hour -> count
	totalMap := make(map[int]map[int]int)
	var latch sync.WaitGroup
	latch.Add(len(dayList))

	for _, day := range dayList {
		var curDay = day
		go func() {
			defer latch.Done()

			readDetailToMap(curDay, mutex, totalMap)
		}()
	}
	latch.Wait()

	max := 0
	for weekday, v := range totalMap {
		chartIndex := 6 - weekday
		for hour, count := range v {
			//logger.Info(weekday, hour)
			if count > max {
				max = count
			}
			result[(chartIndex*24)+hour] = [...]int{chartIndex, hour, count}
		}
	}

	return &HeatMapVO{
		Max:   max,
		Data:  result,
		Start: dayList[0],
		End:   dayList[len(dayList)-1],
	}
}

func readDetailToMap(
	curDay string,
	mutex *sync.Mutex,
	totalMap map[int]map[int]int) {

	var lastCursor uint64 = 0
	first := true

	totalCount := 0
	for lastCursor != 0 || first {
		result, cursor, err := GetConnection().
			ZScan(GetDetailKeyByString(curDay), lastCursor, "", 1000).Result()
		cuibase.CheckIfError(err)
		lastCursor = cursor
		first = false
		for i := range result {
			if i%2 == 1 {
				continue
			}
			//logger.Info(result[i], result[i+1])

			parseInt, err := strconv.ParseInt(result[i], 0, 64)
			cuibase.CheckIfError(err)

			cur := time.Unix(parseInt/1000_000, 0)
			weekDay := int(cur.Weekday())

			mutex.Lock()

			dayMap := totalMap[weekDay]
			//curStr := cur.Format("2006:01:02")
			//if curStr != curDay {
			//	logger.Error("error detail data", curStr, curDay)
			//}
			if dayMap == nil {
				dayMap = make(map[int]int)
				totalMap[weekDay] = dayMap
			}
			dayMap[cur.Hour()] += 1

			mutex.Unlock()
		}
		totalCount += len(result)
	}
	//logger.Info(day, totalCount/2)
}

//LineMap 折线图 柱状图
func LineMap(c *gin.Context) {
	param, err := parseParam(c)
	if err != nil {
		GinFailedWithMsg(c, err.Error())
		return
	}
	dayList := buildDayList(param.Length, param.Offset)
	hotKey := hotKey(dayList, param.Top)
	nameMap := keyNameMap(hotKey)

	// keyNames
	var keyNames []string
	for _, v := range nameMap {
		keyNames = append(keyNames, v)
	}
	sort.Strings(keyNames)
	if len(keyNames) == 0 {
		GinFailed(c)
		return
	}

	// days
	var days []string
	showDayList := buildDayWithWeekdayList(param.Length, param.Offset)
	for _, day := range showDayList {
		score, err := GetConnection().ZScore(TotalCount, day.Day).Result()
		if err != nil {
			days = append(days, day.Day+"#"+day.WeekDay+"#0")
		} else {
			days = append(days, day.Day+"#"+day.WeekDay+"#"+strconv.Itoa(int(score)))
		}
	}
	if len(days) == 0 {
		GinFailed(c)
		return
	}

	// lines
	sortHotKeys := getMapKeys(hotKey)
	sort.Strings(sortHotKeys)
	var lines []LineVO
	commonLabel.Show = param.ShowLabel
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
			Label:     commonLabel,
			Color:     colorSet[keyCode%len(colorSet)],
		})
	}
	//logger.Info(lines)
	GinSuccessWith(c, LineChartVO{Lines: lines, Days: days, KeyNames: keyNames})
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
	}, nil
}

func keyNameMap(keyCode map[string]bool) map[string]string {
	result := make(map[string]string)
	for k := range keyCode {
		name, err := GetConnection().HGet(KeyMap, k).Result()
		if err != nil {
			result[k] = "unknown"
		}
		result[k] = name
	}
	return result
}

func hotKey(dayList []string, top int64) map[string]bool {
	start := time.Now().UnixNano()
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
	end := time.Now().UnixNano()
	logger.Info("hotKey: ", end-start)
	return keyCodeMap
}

func buildDayList(length int, offset int) []string {
	now := time.Now()

	var result []string
	start := now.AddDate(0, 0, -offset)
	for i := 0; i < length; i++ {
		tempTime := start.AddDate(0, 0, i)
		day := tempTime.Format("2006:01:02")
		if tempTime.After(now) {
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
		day := tempTime.Format("2006:01:02")
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
