package main

import (
	"encoding/json"
	"fmt"
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
		},
		"yAxis": map[string]interface{}{
			"type": "value",
		},
		"series": series,
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func buildLineChartOption(result *AnalysisResult) (string, error) {
	type seriesItem struct {
		Name      string            `json:"name"`
		Type      string            `json:"type"`
		Stack     string            `json:"stack"`
		AreaStyle map[string]string `json:"areaStyle"`
		Data      []int             `json:"data"`
		LineStyle map[string]string `json:"lineStyle,omitempty"`
		ItemStyle map[string]string `json:"itemStyle,omitempty"`
	}

	series := []seriesItem{
		{
			Name: "Added", Type: "line", Stack: "Total",
			AreaStyle: map[string]string{"color": "#4caf50"},
			LineStyle: map[string]string{"color": "#4caf50"},
			ItemStyle: map[string]string{"color": "#4caf50"},
			Data:      result.AddedLineSeries,
		},
		{
			Name: "Deleted", Type: "line", Stack: "Total",
			AreaStyle: map[string]string{"color": "#f44336"},
			LineStyle: map[string]string{"color": "#f44336"},
			ItemStyle: map[string]string{"color": "#f44336"},
			Data:      result.DeletedLineSeries,
		},
	}

	opt := map[string]interface{}{
		"tooltip": map[string]string{"trigger": "axis"},
		"legend": map[string]interface{}{
			"data":   []string{"Added", "Deleted"},
			"bottom": 0,
		},
		"grid": map[string]interface{}{
			"left": "3%", "right": "4%", "bottom": "15%", "top": "3%",
			"containLabel": true,
		},
		"xAxis": map[string]interface{}{
			"type": "category", "data": result.DateRange, "boundaryGap": false,
		},
		"yAxis": map[string]interface{}{"type": "value"},
		"series": series,
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
		},
		"yAxis": map[string]interface{}{"type": "value"},
		"series": series,
	}

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
		},
		"yAxis": map[string]interface{}{
			"type": "category",
			"data": days,
			"inverse": true,
			"splitArea": map[string]interface{}{"show": true},
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
		},
		"yAxis": map[string]interface{}{
			"type": "value",
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
		},
		"yAxis": map[string]interface{}{
			"type": "value",
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
			},
		},
		"yAxis": map[string]interface{}{
			"type": "value",
		},
		"series": series,
	}

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
			},
		},
		"yAxis": map[string]interface{}{
			"type":        "value",
			"minInterval": 1,
		},
		"series": []map[string]interface{}{series},
	}

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
			},
		},
		"yAxis": map[string]interface{}{
			"type":        "value",
			"minInterval": 1,
		},
		"series": []map[string]interface{}{series},
	}

	b, err := json.Marshal(opt)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func authorNames(series []AuthorDayData) []string {
	names := make([]string, len(series))
	for i, s := range series {
		names[i] = s.Name
	}
	return names
}