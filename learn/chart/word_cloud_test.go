package chart

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/kuangcp/gobase/pkg/ctool"
	"github.com/kuangcp/logger"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloud(t *testing.T) {
	createWordCloud(generateWordCloudData(wordCloudData))
}
func latestLogFile() string {
	fi := ""
	err := filepath.WalkDir("log/",
		func(path string, d fs.DirEntry, err error) error {
			if strings.HasSuffix(path, "log") {
				fi = path
			}
			return nil
		})
	if err != nil {
		return ""
	}

	return fi
}

// algorithm/ds/tree/trie_tokenizer_test.go TestDir
func TestReadFile(t *testing.T) {
	cache := make(map[string]interface{})
	ignore := ctool.NewSet("使用", "一个", "可以", "如果", "这个")

	// TestDir 生成的日志文件绝对路径
	//logfile := "log/"
	logfile := latestLogFile()
	logger.Info("parse file", logfile)
	ctool.ReadLines(logfile, func(s string) bool {
		return len(s) > 0
	}, func(s string) bool {
		fields := strings.Fields(s)
		if ignore.Contains(fields[1]) {
			return false
		}
		cache[fields[1]] = fields[0]
		return false
	})
	logger.Info(cache)
	createWordCloud(generateWordCloudData(cache))
}

var wordCloudData = map[string]interface{}{
	"Bitcoin":      10000,
	"Ethereum":     8000,
	"Cardano":      5000,
	"Polygon":      4000,
	"Polkadot":     3000,
	"Chainlink":    2500,
	"Solana":       2000,
	"Ripple":       1500,
	"Decentraland": 1000,
	"Tron":         800,
	"Sandbox":      500,
	"Litecoin":     200,
}

func generateWordCloudData(data map[string]interface{}) (items []opts.WordCloudData) {
	items = make([]opts.WordCloudData, 0)
	for k, v := range data {
		items = append(items, opts.WordCloudData{Name: k, Value: v})
	}
	return
}

func createWordCloud(items []opts.WordCloudData) {
	wc := charts.NewWordCloud()
	wc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Note WordCloud",
			Subtitle: "all markdown file",
		}))
	wc.AddSeries("wordcloud", items).
		SetSeriesOptions(
			charts.WithWorldCloudChartOpts(
				opts.WordCloudChart{
					SizeRange: []float32{10, 42},
					// The shape of the "cloud" to draw. Can be any polar equation represented as a
					// callback function, or a keyword present.
					//
					//Available presents are circle (default),
					// cardioid (apple or heart shape curve, the most known polar equation), diamond (alias of square),
					// triangle-forward, triangle, (alias of triangle-upright, pentagon,
					//Shape: "triangle",
					Shape:         "circle",
					RotationRange: []float32{-30, 30},
				}),
		)
	f, _ := os.Create("word_cloud.html")
	_ = wc.Render(f)
}
