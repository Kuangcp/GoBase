package chart

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type Task struct {
	Name   string
	Start  time.Time
	End    time.Time
	Data   []opts.LineData
	Run    int64
	Weight int
}

func (t *Task) buildData(min, max, frame int64) {
	items := make([]opts.LineData, 0)

	start := t.Start.UnixMilli()
	end := t.End.UnixMilli()
	for min < max {
		if min > start && min < end {
			items = append(items, opts.LineData{Value: t.Weight})
		} else {
			items = append(items, opts.LineData{Value: nil})
		}
		min += frame
	}
	t.Data = items
}

func TestAnalyzeLog(t *testing.T) {
	logFile := "/home/zk/Downloads/firefox/1727597442521.log"
	finishLog := ctool.ReadStrLines(logFile, func(s string) bool {
		return strings.Contains(s, "records/s")
	})
	for _, l := range finishLog {
		//a := strings.Index(l, "StandAloneJobContainerCommunicator")
		b := strings.Index(l, "2024-")
		c := strings.Index(l, "Speed ")
		//fmt.Print(l[b:b+23], l[a+34:])
		fmt.Println(l[b:b+23], l[c+6:c+29])
	}
}

// 解析 datax log 得到每个任务运行时间段
func TestExtractLog(t *testing.T) {
	logFile := "/home/zk/Downloads/firefox/1727598245788.log"
	csvFile := "run-data29.csv"
	finishLog := ctool.ReadStrLines(logFile, func(s string) bool {
		return strings.Contains(s, "is successed, used")
	})

	//格式：id,start,end,weight
	writer := ctool.NewWriterIgnoreError(csvFile, true)
	defer writer.Close()
	for _, l := range finishLog {
		//fmt.Println(l)

		a := strings.Index(l, "taskId")
		b := strings.Index(l, "is successed")
		c := strings.Index(l, "2024-")
		name := l[a+7 : b-2]
		finishT := l[c : c+23]
		waste := l[b+19 : len(l)-4]

		fmt.Println(name, finishT, waste)

		// 注意time.Parse默认0时区
		finish, err := time.ParseInLocation("2006-01-02 15:04:05.000", finishT, time.Local)
		if err != nil {
			panic(err)
		}

		runMs, err := strconv.Atoi(waste)
		if err != nil {
			panic(err)
		}

		end := finish.Format("2006-01-02 15:04:05.000")
		start := finish.Add(-time.Millisecond * time.Duration(runMs)).Format("2006-01-02 15:04:05.000")

		writer.WriteLine(name + ", " + start + ", " + end)
	}
}

func TestConcatWeight(t *testing.T) {
	inCsv := "run-data29.csv"
	outCsv := "run-data29-w.csv"

	writer := ctool.NewWriterIgnoreError(outCsv, true)
	defer writer.Close()
	lines := ctool.ReadCsvLines(inCsv)
	// 27 28 使用period_int做拆分
	//var weight = []int64{4466162, 4555800, 3535773, 4846170,
	//	4996228, 3866885, 6370524, 7560920, 5747293, 7965687, 8302327, 6876701, 10128650, 7727466, 10768946, 11283096,
	//	9221867, 12960388, 13560612, 10678018, 14756670, 14868640, 11716800, 16813336, 16786591, 1}

	// 29 使用月份拆分
	var weight = []int64{9958032, 9958032, 9958032, 9958032, 9958032, 9958032,
		9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020, 9958020,
		8704046, 6897405, 6638680, 6174019, 2786928, 1,
	}

	for _, line := range lines {
		if len(line) == 0 || line[0] == "" {
			continue
		}
		name := line[0]
		//name = strings.Replace(name, "taskId", "", 1)
		//name = strings.Replace(name, "[", "", 1)
		//name = strings.Replace(name, "]", "", 1)
		taskId, err := strconv.Atoi(name)
		if err != nil {
			panic(err)
		}

		logger.Info("", taskId, weight[taskId])

		writer.WriteLine(name + ", " + line[1] + ", " + line[2] + ", " + fmt.Sprint(weight[taskId]))
	}
}

var (
	// 忽略权重区别
	ignoreWeight = false

	// 时间划分的段数
	point         int64 = 300
	width, height       = 6500, 600

	// 比例缩小权重
	weightFold = 100_000

	// 超大图样式宽度的阈值
	hugeWidthFlag = 1500
	noSelect      = true
)

// 生成 运行序列 仿甘特图
func TestRunSerial(t *testing.T) {
	tasks, minT, maxT := parseCsv("run-data29-w.csv")

	logger.Info("parse: ", len(tasks), time.UnixMilli(minT), time.UnixMilli(maxT))
	if ignoreWeight {
		for _, task := range tasks {
			task.Weight = 1
		}
	}

	renderLineChart(maxT, minT, tasks)
}

// parseCsv 格式：id,start,end,weight
func parseCsv(path string) ([]*Task, int64, int64) {
	lines := ctool.ReadCsvLines(path)
	var tasks []*Task
	var minT, maxT int64 = time.Now().UnixMilli(), 0
	for _, row := range lines {
		if len(row) == 0 || len(row) == 1 {
			continue
		}
		name := row[0]
		task := &Task{Name: name}

		// 注意time.Parse默认0时区
		parse1, err := time.ParseInLocation("2006-01-02 15:04:05.000", strings.TrimSpace(row[1]), time.Local)
		if err != nil {
			panic(err)
		}
		parse2, err := time.ParseInLocation("2006-01-02 15:04:05.000", strings.TrimSpace(row[2]), time.Local)
		if err != nil {
			panic(err)
		}

		weight := 1
		if len(row) > 3 {
			atoi, err := strconv.Atoi(strings.TrimSpace(row[3]))
			if err == nil {
				weight = atoi / weightFold
			}
		}
		task.Start = parse1
		task.End = parse2
		task.Weight = weight
		tasks = append(tasks, task)
		if minT > task.Start.UnixMilli() {
			minT = task.Start.UnixMilli()
		}
		if maxT < task.End.UnixMilli() {
			maxT = task.End.UnixMilli()
		}
	}
	return tasks, minT, maxT
}

func renderLineChart(maxT int64, minT int64, tasks []*Task) {
	frame := (maxT - minT) / point
	var names []string
	for _, task := range tasks {
		names = append(names, task.Name)
	}
	sort.Slice(names, func(i, j int) bool {
		a, _ := strconv.Atoi(names[i])
		b, _ := strconv.Atoi(names[j])

		return a < b
	})

	sta := make(map[string]bool)
	for _, name := range names {
		sta[name] = noSelect
	}

	ch := charts.NewLine()
	show := true
	ch.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Gantt",
			Subtitle: "Datax task run serial",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  fmt.Sprint(width) + "px",
			Height: fmt.Sprint(height) + "px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Data:     names,
			Selected: sta,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{

				Type: "cross",
				Label: &opts.Label{
					BackgroundColor: "#6a7985",
				},
			},
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{Show: &show},
			},
		}),
	)

	for _, task := range tasks {
		task.buildData(minT, maxT, frame)

		ch.AddSeries(task.Name, task.Data).SetSeriesOptions(
			charts.WithLineChartOpts(
				opts.LineChart{
					Stack: "Total",
				}),
			charts.WithAreaStyleOpts(
				opts.AreaStyle{
					Opacity: 0.2,
				}),
			charts.WithLabelOpts(
				opts.Label{
					Show: opts.Bool(true),
				}),
		)
	}

	var xs []string
	cursor := minT
	for cursor < maxT {
		flag := time.UnixMilli(cursor).Format("15:04:05")
		xs = append(xs, flag)
		cursor += frame
	}

	ch.SetXAxis(xs)

	f, _ := os.Create("gantt.html")
	html := ch.RenderContent()
	tmp := string(html)

	if width > hugeWidthFlag {
		tmp = strings.Replace(tmp, "justify-content: center;", "", 1)
		tmp = strings.Replace(tmp, "margin-top:30px;", "margin-top:300px;", 1)
	}
	f.Write([]byte(tmp))

	//f, _ := os.Create("gantt.html")
	//_ = ch.Render(f)
}
