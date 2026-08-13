package main

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"
)

func buildCommitChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Stack     string   `json:"stack"`
		AreaStyle struct{} `json:"areaStyle"`
		Data      []int    `json:"data"`
	}

	var series []seriesItem
	for _, s := range result.AuthorSeries {
		item := seriesItem{
			Name:  s.Name,
			Type:  "line",
			Stack: "Total",
			Data:  s.Data,
		}
		item.AreaStyle = struct{}{}
		series = append(series, item)
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   authorNames(result.AuthorSeries),
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type":        "category",
			"data":        result.DateRange,
			"boundaryGap": false,
			"axisLabel":   map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": series,
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildLineChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Stack     string   `json:"stack"`
		AreaStyle struct{} `json:"areaStyle"`
		Data      []int    `json:"data"`
	}

	var series []seriesItem
	for i, a := range result.AuthorAddedSeries {
		netData := make([]int, len(a.Data))
		for j := range a.Data {
			netData[j] = a.Data[j]
			if j < len(result.AuthorDeletedSeries[i].Data) {
				netData[j] -= result.AuthorDeletedSeries[i].Data[j]
			}
		}
		item := seriesItem{
			Name:  a.Name,
			Type:  "line",
			Stack: "Total",
			Data:  netData,
		}
		item.AreaStyle = struct{}{}
		series = append(series, item)
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   authorNames(result.AuthorAddedSeries),
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.DateRange, "boundaryGap": false,
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": series,
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func hasAny(data []int) bool {
	for _, v := range data {
		if v > 0 {
			return true
		}
	}
	return false
}

func buildAuthorLineChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string   `json:"name"`
		Type      string   `json:"type"`
		Stack     string   `json:"stack"`
		AreaStyle struct{} `json:"areaStyle"`
		Data      []int    `json:"data"`
	}

	var series []seriesItem
	for i, s := range result.AuthorAddedSeries {
		total := make([]int, len(s.Data))
		for j := range s.Data {
			total[j] = s.Data[j]
			if j < len(result.AuthorDeletedSeries[i].Data) {
				total[j] += result.AuthorDeletedSeries[i].Data[j]
			}
		}
		item := seriesItem{
			Name:  s.Name,
			Type:  "line",
			Stack: "Total",
			Data:  total,
		}
		item.AreaStyle = struct{}{}
		series = append(series, item)
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   authorNames(result.AuthorAddedSeries),
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.DateRange, "boundaryGap": false,
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": series,
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildHourWeekChartOption(result *AnalysisResult) (string, error) {
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	var hours []string
	for i := 0; i < 24; i++ {
		hours = append(hours, fmt.Sprintf("%02d", i))
	}

	maxVal := 0
	var data [][3]int
	for d := 0; d < 7; d++ {
		for h := 0; h < 24; h++ {
			v := result.HourWeekData[d][h]
			if v > maxVal {
				maxVal = v
			}
			data = append(data, [3]int{h, d, v})
		}
	}

	opt := map[string]interface{}{
		"tooltip": map[string]interface{}{
			"position": "top",
		},
		"grid": map[string]interface{}{
			"left":        "2%",
			"right":       "4%",
			"bottom":      "15%",
			"top":         "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category",
			"data": hours,
			"splitArea": map[string]interface{}{"show": true},
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type": "category",
			"data": days,
			"inverse": true,
			"splitArea": map[string]interface{}{"show": true},
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"visualMap": map[string]interface{}{
			"min":        0,
			"max":        maxVal,
			"calculable": true,
			"orient":     "horizontal",
			"left":       "center",
			"bottom":     0,
			"inRange": map[string]interface{}{
				"color": []string{"#f5f5f5", "#c6e48b", "#7bc96f", "#239a3b", "#196127"},
			},
		},
		"series": []map[string]interface{}{
			{
				"name": "Commits",
				"type": "heatmap",
				"data": data,
				"label": map[string]interface{}{
					"show": false,
				},
				"emphasis": map[string]interface{}{
					"itemStyle": map[string]interface{}{
						"shadowBlur":  10,
						"shadowColor": "rgba(0,0,0,0.5)",
					},
				},
			},
		},
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildMonthOfYearChartOption(result *AnalysisResult) (string, error) {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "10%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": months,
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": []map[string]interface{}{
			{
				"type": "bar",
				"data": result.MonthOfYearData[:],
				"itemStyle": map[string]interface{}{
					"color": "#5470c6",
				},
			},
		},
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildWeekOfYearChartOption(result *AnalysisResult) (string, error) {
	colors := []string{"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de",
		"#3ba272", "#fc8452", "#9a60b4", "#ea7ccc", "#97b552", "#95706d"}

	type seriesItem struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		Stack     string `json:"stack"`
		Data      []int  `json:"data"`
		ItemStyle struct {
			Color string `json:"color"`
		} `json:"itemStyle"`
	}

	var series []seriesItem

	// Sort authors by total commits descending
	type authorTotal struct {
		name  string
		data  []int
		total int
	}
	var ats []authorTotal
	for _, s := range result.AuthorWeekOfYearSeries {
		total := 0
		for _, v := range s.Data {
			total += v
		}
		ats = append(ats, authorTotal{name: s.Name, data: s.Data, total: total})
	}
	sort.Slice(ats, func(i, j int) bool {
		return ats[i].total > ats[j].total
	})

	topN := 10
	if len(ats) < topN {
		topN = len(ats)
	}

	for i := 0; i < topN; i++ {
		item := seriesItem{
			Name:  ats[i].name,
			Type:  "bar",
			Stack: "total",
			Data:  ats[i].data,
		}
		item.ItemStyle.Color = colors[i%len(colors)]
		series = append(series, item)
	}

	// Merge remaining into "Others"
	if len(ats) > topN {
		others := make([]int, len(result.WeekOfYearLabels))
		for i := topN; i < len(ats); i++ {
			for j := range ats[i].data {
				others[j] += ats[i].data[j]
			}
		}
		item := seriesItem{
			Name:  "Others",
			Type:  "bar",
			Stack: "total",
			Data:  others,
		}
		item.ItemStyle.Color = "#b0b0b0"
		series = append(series, item)
	}

	var legendNames []string
	for _, s := range series {
		legendNames = append(legendNames, s.Name)
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   legendNames,
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.WeekOfYearLabels,
			"axisLabel": map[string]interface{}{"color": "#e0e0e0", "rotate": 45},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"dataZoom": []map[string]interface{}{
			{"type": "inside", "start": 0, "end": 100},
			{"type": "slider", "start": 0, "end": 100, "bottom": 30},
		},
		"series": series,
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildYearMonthChartOption(result *AnalysisResult) (string, error) {
	n := len(result.YearMonthLabels)
	revLabels := make([]string, n)
	revData := make([]int, n)
	for i := 0; i < n; i++ {
		revLabels[i] = result.YearMonthLabels[n-1-i]
		revData[i] = result.YearMonthData[n-1-i]
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   []string{"Commits"},
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": revLabels,
			"boundaryGap": false,
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": []map[string]interface{}{
			{
				"name": "Commits",
				"type": "line",
				"smooth": true,
				"data": revData,
				"areaStyle": map[string]interface{}{
					"color": "rgba(84,112,198,0.2)",
				},
				"lineStyle": map[string]interface{}{
					"color": "#5470c6",
				},
				"itemStyle": map[string]interface{}{
					"color": "#5470c6",
				},
			},
		},
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildYearChartOption(result *AnalysisResult) (string, error) {
	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "10%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.YearLabels,
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": []map[string]interface{}{
			{
				"name": "Commits",
				"type": "bar",
				"data": result.YearData,
				"itemStyle": map[string]interface{}{
					"color": "#5470c6",
				},
			},
		},
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildCumChartOption(result *AnalysisResult, isCommit bool) (string, error) {
	var source []AuthorDayData
	if isCommit {
		source = result.AuthorCumCommitSeries
	} else {
		source = result.AuthorCumAddedSeries
	}

	dayLabels := result.AllDayLabels
	if len(dayLabels) == 0 {
		dayLabels = result.YearMonthLabels
	}

	type seriesItem struct {
		Name      string                `json:"name"`
		Type      string                `json:"type"`
		Symbol    string                `json:"symbol"`
		Data      [][]interface{}       `json:"data"`
		LineStyle map[string]interface{} `json:"lineStyle,omitempty"`
	}

	var series []seriesItem
	for _, s := range source {
		data := make([][]interface{}, len(s.Data))
		for i, v := range s.Data {
			t, _ := time.Parse("2006-01-02", dayLabels[i])
			data[i] = []interface{}{t.UnixMilli(), v}
		}
		series = append(series, seriesItem{
			Name:   s.Name,
			Type:   "line",
			Symbol: "none",
			Data:   data,
		})
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   authorNames(source),
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "time",
			"axisLabel": map[string]interface{}{
				"formatter": `{yyyy}-{MM}`,
				"color":     "#e0e0e0",
			},
		},
		"yAxis": map[string]interface{}{
			"type":      "value",
			"axisLabel": map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": series,
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildFileChartOption(result *AnalysisResult) (string, error) {
	var data [][]interface{}
	for i, v := range result.FileChartData {
		t, err := time.Parse("2006-01-02", result.FileChartLabels[i])
		if err != nil {
			continue
		}
		data = append(data, []interface{}{t.UnixMilli(), v})
	}

	series := map[string]interface{}{
		"name":      "Files",
		"type":      "line",
		"symbol":    "none",
		"data":      data,
		"areaStyle": map[string]interface{}{"color": "rgba(84,112,198,0.2)"},
		"lineStyle": map[string]interface{}{"color": "#5470c6"},
		"itemStyle": map[string]interface{}{"color": "#5470c6"},
	}

	opt := map[string]interface{}{
		"tooltip": map[string]interface{}{
			"trigger": "axis",
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "10%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "time",
			"axisLabel": map[string]interface{}{
				"formatter": "{yyyy}-{MM}",
				"color":     "#e0e0e0",
			},
		},
		"yAxis": map[string]interface{}{
			"type":        "value",
			"minInterval": 1,
			"axisLabel":   map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": []map[string]interface{}{series},
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildLocChartOption(result *AnalysisResult) (string, error) {
	var data [][]interface{}
	for i, v := range result.LocChartData {
		t, err := time.Parse("2006-01-02", result.LocChartLabels[i])
		if err != nil {
			continue
		}
		data = append(data, []interface{}{t.UnixMilli(), v})
	}

	series := map[string]interface{}{
		"name":      "Lines of Code",
		"type":      "line",
		"symbol":    "none",
		"data":      data,
		"areaStyle": map[string]interface{}{"color": "rgba(76,175,80,0.2)"},
		"lineStyle": map[string]interface{}{"color": "#4caf50"},
		"itemStyle": map[string]interface{}{"color": "#4caf50"},
	}

	opt := map[string]interface{}{
		"tooltip": map[string]interface{}{
			"trigger": "axis",
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "10%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "time",
			"axisLabel": map[string]interface{}{
				"formatter": "{yyyy}-{MM}",
				"color":     "#e0e0e0",
			},
		},
		"yAxis": map[string]interface{}{
			"type":        "value",
			"minInterval": 1,
			"axisLabel":   map[string]interface{}{"color": "#e0e0e0"},
		},
		"series": []map[string]interface{}{series},
	}

	applyDarkTheme(opt)
	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildRadarChartOption(scores []float64, overallGrade string) (string, error) {
	roundedScores := make([]int, len(scores))
	for i, s := range scores {
		roundedScores[i] = int(math.Round(s))
	}
	opt := map[string]interface{}{
		"tooltip": map[string]interface{}{},
		"radar": map[string]interface{}{
			"indicator": []map[string]interface{}{
				{"name": "活跃度", "max": 100},
				{"name": "项目规模", "max": 100},
				{"name": "代码健康", "max": 100},
				{"name": "协作多样", "max": 100},
				{"name": "技术债", "max": 100},
				{"name": "研发节奏", "max": 100},
			},
			"shape":  "circle",
			"center": []interface{}{"50%", "50%"},
			"radius": "65%",
			"axisName": map[string]interface{}{
				"textStyle": map[string]interface{}{"color": "#e0e0e0"},
			},
			"splitArea": map[string]interface{}{
				"areaStyle": map[string]interface{}{
					"color": []string{"rgba(26,26,46,0.3)", "rgba(26,26,46,0.1)"},
				},
			},
			"splitLine": map[string]interface{}{
				"lineStyle": map[string]interface{}{"color": "rgba(42,42,74,0.6)"},
			},
		},
		"series": []map[string]interface{}{
			{
				"type": "radar",
				"name": "综合评级：" + overallGrade,
				"data": []map[string]interface{}{
					{
						"value": roundedScores,
						"areaStyle": map[string]interface{}{
							"color": "rgba(84,112,198,0.3)",
						},
						"lineStyle": map[string]interface{}{
							"color": "#5470c6", "width": 2,
						},
						"itemStyle": map[string]interface{}{
							"color": "#5470c6",
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func applyDarkTheme(opt map[string]interface{}) {
	opt["textStyle"] = map[string]interface{}{"color": "#e0e0e0"}
	if legend, ok := opt["legend"].(map[string]interface{}); ok {
		legend["textStyle"] = map[string]interface{}{"color": "#e0e0e0"}
	}
}

func authorNames(series []AuthorDayData) []string {
	names := make([]string, len(series))
	for i, s := range series {
		names[i] = s.Name
	}
	return names
}