package ctool

import (
	"fmt"
	"log"
	"sort"
)

type (
	NumberDis struct {
		Tips   string
		MaxVal float64
		MinVal float64
		Avg    float64
		P30    float64
		P50    float64
		P75    float64
		P90    float64
		P95    float64
	}
)

func NumberDistribution(data []float64) *NumberDis {
	if len(data) == 0 {
		log.Println("数据为空")
		return nil
	}

	// 过滤掉 0 值
	filtered := make([]float64, 0, len(data))
	for _, v := range data {
		if v != 0 {
			filtered = append(filtered, v)
		}
	}

	if len(filtered) == 0 {
		//log.Println("过滤后数据为空（所有值都是 0）")
		return nil
	}

	// 创建副本并排序
	sorted := make([]float64, len(filtered))
	copy(sorted, filtered)
	sort.Float64s(sorted)

	// 计算最大值和最小值
	maxVal := sorted[len(sorted)-1]
	minVal := sorted[0]

	// 计算平均值
	var sum float64
	for _, v := range filtered {
		sum += v
	}
	avg := sum / float64(len(filtered))

	// 计算百分位数的辅助函数
	percentile := func(p float64) float64 {
		if len(sorted) == 0 {
			return 0
		}
		if len(sorted) == 1 {
			return sorted[0]
		}
		index := p * float64(len(sorted)-1)
		lower := int(index)
		upper := lower + 1
		if upper >= len(sorted) {
			return sorted[len(sorted)-1]
		}
		weight := index - float64(lower)
		return sorted[lower]*(1-weight) + sorted[upper]*weight
	}

	// 计算各个百分位数
	p30 := percentile(0.30)
	p50 := percentile(0.50)
	p75 := percentile(0.75)
	p90 := percentile(0.90)
	p95 := percentile(0.95)

	// 打印结果
	tips := ""
	filLen := len(filtered)
	allLen := len(data)
	if allLen > filLen {
		tips += fmt.Sprintf("数据统计 (原始数据 %d 条，过滤后 %d 条，已忽略 %d 个 0 值):\n",
			allLen, filLen, allLen-filLen)
	} else {
		tips += fmt.Sprintf("数据统计 %d 条:\n", allLen)
	}

	tips += fmt.Sprintf("  最大值: %.2f\n", maxVal)
	tips += fmt.Sprintf("  最小值: %.2f\n", minVal)
	tips += fmt.Sprintf("  平均值: %.2f\n", avg)
	tips += fmt.Sprintf("  30%%水位线: %.2f\n", p30)
	tips += fmt.Sprintf("  50%%水位线: %.2f\n", p50)
	tips += fmt.Sprintf("  75%%水位线: %.2f\n", p75)
	tips += fmt.Sprintf("  90%%水位线: %.2f\n", p90)
	tips += fmt.Sprintf("  95%%水位线: %.2f\n", p95)
	return &NumberDis{
		Tips:   tips,
		MaxVal: maxVal,
		MinVal: minVal,
		Avg:    avg,
		P30:    p30,
		P50:    p50,
		P75:    p75,
		P90:    p90,
		P95:    p95,
	}
}
