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

func TestConcatWeight(t *testing.T) {
	writer := ctool.NewWriterIgnoreError("run-data.csv", true)
	defer writer.Close()
	lines := ctool.ReadCsvLines("/home/zk/Work/supersonic/sssss.csv")
	var weight = []int64{4466162,
		4555800,
		3535773,
		4846170,
		4996228,
		3866885,
		6370524,
		7560920,
		5747293,
		7965687,
		8302327,
		6876701,
		10128650,
		7727466,
		10768946,
		11283096,
		9221867,
		12960388,
		13560612,
		10678018,
		14756670,
		14868640,
		11716800,
		16813336,
		16786591, 1}
	for _, line := range lines {
		name := line[0]
		name = strings.Replace(name, "taskId", "", 1)
		name = strings.Replace(name, "[", "", 1)
		name = strings.Replace(name, "]", "", 1)
		taskId, err := strconv.Atoi(name)
		if err != nil {
			panic(err)
		}

		logger.Info("", taskId, weight[taskId])

		writer.WriteLine(name + ", " + line[1] + ", " + line[2] + ", " + fmt.Sprint(weight[taskId]))
	}
}

// 运行序列
func TestRunSerial(t *testing.T) {
	tasks, minT, maxT := parseCsv("run-data.csv")
	renderLineChart(maxT, minT, tasks)
}

func parseCsv(path string) ([]*Task, int64, int64) {
	// id,start,end,weight
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
		atoi, err := strconv.Atoi(strings.TrimSpace(row[3]))
		if err != nil {
			atoi = 1
		} else {
			atoi = atoi / 100_000
		}
		task.Start = parse1
		task.End = parse2
		task.Weight = atoi
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
	var point int64 = 100
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

	ch := charts.NewLine()
	ch.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Gantt",
			Subtitle: "all markdown file",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1600px",
			Height: "600px",
		}),
		charts.WithLegendOpts(opts.Legend{
			Data: names,
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
	_ = ch.Render(f)
}
