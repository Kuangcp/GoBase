package ctool

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrentProgressTracker 并发安全的进度统计器
// 通过记录任务的开始和结束时间来计算剩余任务的执行耗时
type ConcurrentProgressTracker struct {
	total     int64     // 总任务数
	current   int64     // 已完成任务数（使用 atomic 操作）
	startTime time.Time // 整体开始时间
	width     int       // 进度条宽度

	// 任务耗时统计（使用滑动窗口）
	mu            sync.RWMutex    // 保护耗时统计
	taskDurations []time.Duration // 最近N个任务的耗时
	maxSamples    int             // 最多保存的样本数
	totalDuration time.Duration   // 所有已完成任务的总耗时（用于计算平均值）

	// 正在执行的任务
	activeTasks map[int64]time.Time // taskID -> startTime
	nextTaskID  int64               // 下一个任务ID
}

// NewConcurrentProgressTracker 创建新的并发安全进度统计器
func NewConcurrentProgressTracker() *ConcurrentProgressTracker {
	return &ConcurrentProgressTracker{
		total:         0,
		current:       0,
		startTime:     time.Now(),
		width:         50,
		taskDurations: make([]time.Duration, 0, 100),
		maxSamples:    100, // 默认保存最近100个任务的耗时
		totalDuration: 0,
		activeTasks:   make(map[int64]time.Time),
		nextTaskID:    1,
	}
}

// SetTotal 设置总任务数
func (cpt *ConcurrentProgressTracker) SetTotal(total int64) {
	if total < 0 {
		total = 0
	}
	atomic.StoreInt64(&cpt.total, total)
}

// SetWidth 设置进度条宽度
func (cpt *ConcurrentProgressTracker) SetWidth(width int) {
	if width > 0 {
		cpt.mu.Lock()
		cpt.width = width
		cpt.mu.Unlock()
	}
}

// SetMaxSamples 设置最多保存的样本数
func (cpt *ConcurrentProgressTracker) SetMaxSamples(maxSamples int) {
	if maxSamples > 0 {
		cpt.mu.Lock()
		cpt.maxSamples = maxSamples
		// 如果当前样本数超过新的限制，裁剪到新的大小
		if len(cpt.taskDurations) > maxSamples {
			cpt.taskDurations = cpt.taskDurations[len(cpt.taskDurations)-maxSamples:]
		}
		cpt.mu.Unlock()
	}
}

func (cpt *ConcurrentProgressTracker) RunTask(fn func(), callback func()) {
	taskID := cpt.StartTask()
	fn()
	cpt.FinishTask(taskID)
	callback()
}

// StartTask 标记任务开始，返回任务ID
func (cpt *ConcurrentProgressTracker) StartTask() int64 {
	cpt.mu.Lock()
	defer cpt.mu.Unlock()

	taskID := cpt.nextTaskID
	cpt.nextTaskID++
	cpt.activeTasks[taskID] = time.Now()
	return taskID
}

// FinishTask 标记任务完成
func (cpt *ConcurrentProgressTracker) FinishTask(taskID int64) {
	cpt.mu.Lock()
	defer cpt.mu.Unlock()

	startTime, exists := cpt.activeTasks[taskID]
	if !exists {
		// 任务不存在，可能是重复调用或无效ID
		return
	}

	// 计算任务耗时
	duration := time.Since(startTime)

	// 更新统计
	cpt.totalDuration += duration

	// 添加到滑动窗口
	cpt.taskDurations = append(cpt.taskDurations, duration)
	if len(cpt.taskDurations) > cpt.maxSamples {
		// 移除最旧的样本
		oldest := cpt.taskDurations[0]
		cpt.totalDuration -= oldest
		cpt.taskDurations = cpt.taskDurations[1:]
	}

	// 移除活跃任务
	delete(cpt.activeTasks, taskID)

	// 更新当前进度（原子操作）
	atomic.AddInt64(&cpt.current, 1)
}

// Update 简单的增量更新（向后兼容，但不记录任务耗时）
func (cpt *ConcurrentProgressTracker) Update(delta int64) {
	if delta <= 0 {
		return
	}
	atomic.AddInt64(&cpt.current, delta)
}

// SetCurrent 直接设置当前进度值
func (cpt *ConcurrentProgressTracker) SetCurrent(current int64) {
	if current < 0 {
		current = 0
	}
	atomic.StoreInt64(&cpt.current, current)
}

// GetAverageTaskDuration 获取平均任务耗时
func (cpt *ConcurrentProgressTracker) GetAverageTaskDuration() time.Duration {
	cpt.mu.RLock()
	defer cpt.mu.RUnlock()

	if len(cpt.taskDurations) == 0 {
		return 0
	}
	return cpt.totalDuration / time.Duration(len(cpt.taskDurations))
}

// GetCurrent 获取当前进度值
func (cpt *ConcurrentProgressTracker) GetCurrent() int64 {
	return atomic.LoadInt64(&cpt.current)
}

// GetTotal 获取总任务数
func (cpt *ConcurrentProgressTracker) GetTotal() int64 {
	return atomic.LoadInt64(&cpt.total)
}

// GetPercentage 获取当前进度百分比
func (cpt *ConcurrentProgressTracker) GetPercentage() float64 {
	total := atomic.LoadInt64(&cpt.total)
	if total <= 0 {
		return 0
	}
	current := atomic.LoadInt64(&cpt.current)
	percentage := float64(current) / float64(total) * 100
	if percentage > 100 {
		return 100
	}
	return percentage
}

// IsComplete 检查是否已完成
func (cpt *ConcurrentProgressTracker) IsComplete() bool {
	total := atomic.LoadInt64(&cpt.total)
	if total <= 0 {
		return false
	}
	current := atomic.LoadInt64(&cpt.current)
	return current >= total
}

// String 拼接进度条字符串
// 格式: [#########>        ] 45% (123/456) 12.34s elapsed, 15.67s remaining
func (cpt *ConcurrentProgressTracker) String() string {
	total := atomic.LoadInt64(&cpt.total)
	current := atomic.LoadInt64(&cpt.current)

	if total <= 0 {
		cpt.mu.RLock()
		width := cpt.width
		cpt.mu.RUnlock()
		return fmt.Sprintf("[%s] 0%% (0/0) calculating...", strings.Repeat(" ", width))
	}

	// 计算百分比
	percentage := float64(current) / float64(total) * 100
	if percentage > 100 {
		percentage = 100
	}

	// 计算已执行时间
	elapsed := time.Since(cpt.startTime)
	elapsedSeconds := elapsed.Seconds()

	// 计算剩余时间（基于已完成任务的平均耗时）
	var remainingSeconds float64
	cpt.mu.RLock()
	var avgDuration time.Duration
	if len(cpt.taskDurations) > 0 {
		avgDuration = cpt.totalDuration / time.Duration(len(cpt.taskDurations))
	}
	remainingItems := total - current
	width := cpt.width
	cpt.mu.RUnlock()

	if avgDuration > 0 && remainingItems > 0 {
		remainingSeconds = avgDuration.Seconds() * float64(remainingItems)
	} else if current > 0 {
		// 如果没有任务耗时统计，使用总时间估算
		avgTimePerItem := elapsedSeconds / float64(current)
		remainingSeconds = avgTimePerItem * float64(remainingItems)
	} else {
		remainingSeconds = 0
	}

	// 构建进度条

	filled := int(float64(width) * percentage / 100)
	if filled > width {
		filled = width
	}
	empty := width - filled

	bar := strings.Repeat("█", filled)
	if filled < width {
		bar += "░"
		empty--
		if empty > 0 {
			bar += strings.Repeat(" ", empty)
		}
	}

	// 格式化时间
	elapsedStr := formatDuration(elapsedSeconds)
	remainingStr := formatDuration(remainingSeconds)

	return fmt.Sprintf("[%s] %.1f%% (%d/%d) %6s elapsed, %6s remaining",
		bar, percentage, current, total, elapsedStr, remainingStr)
}

// Reset 重置进度统计器
func (cpt *ConcurrentProgressTracker) Reset() {
	atomic.StoreInt64(&cpt.current, 0)
	cpt.mu.Lock()
	cpt.startTime = time.Now()
	cpt.taskDurations = cpt.taskDurations[:0]
	cpt.totalDuration = 0
	cpt.activeTasks = make(map[int64]time.Time)
	cpt.nextTaskID = 1
	cpt.mu.Unlock()
}

// GetActiveTaskCount 获取当前正在执行的任务数
func (cpt *ConcurrentProgressTracker) GetActiveTaskCount() int {
	cpt.mu.RLock()
	defer cpt.mu.RUnlock()
	return len(cpt.activeTasks)
}

// TaskStats 任务统计信息
type TaskStats struct {
	AvgDuration time.Duration // 平均耗时
	MinDuration time.Duration // 最小耗时
	MaxDuration time.Duration // 最大耗时
	P30Duration time.Duration // P30 百分位耗时
	P50Duration time.Duration // P50 百分位耗时（中位数）
	P75Duration time.Duration // P75 百分位耗时
	P90Duration time.Duration // P90 百分位耗时
	P95Duration time.Duration // P95 百分位耗时
	P99Duration time.Duration // P99 百分位耗时
	SampleCount int           // 样本数量
}

// GetTaskStats 获取任务统计信息（旧版本，保持向后兼容）
func (cpt *ConcurrentProgressTracker) GetTaskStats() (avgDuration, minDuration, maxDuration time.Duration, sampleCount int) {
	stats := cpt.GetTaskStatsDetail()
	return stats.AvgDuration, stats.MinDuration, stats.MaxDuration, stats.SampleCount
}

// GetTaskStatsDetail 获取详细的任务统计信息（包含百分位数）
func (cpt *ConcurrentProgressTracker) GetTaskStatsDetail() TaskStats {
	cpt.mu.RLock()
	defer cpt.mu.RUnlock()

	sampleCount := len(cpt.taskDurations)
	if sampleCount == 0 {
		return TaskStats{}
	}

	// 复制一份用于排序（不影响原数据）
	sortedDurations := make([]time.Duration, sampleCount)
	copy(sortedDurations, cpt.taskDurations)
	sort.Slice(sortedDurations, func(i, j int) bool {
		return sortedDurations[i] < sortedDurations[j]
	})

	stats := TaskStats{
		SampleCount: sampleCount,
		AvgDuration: cpt.totalDuration / time.Duration(sampleCount),
		MinDuration: sortedDurations[0],
		MaxDuration: sortedDurations[sampleCount-1],
		P30Duration: calculatePercentile(sortedDurations, 0.30),
		P50Duration: calculatePercentile(sortedDurations, 0.50),
		P75Duration: calculatePercentile(sortedDurations, 0.75),
		P90Duration: calculatePercentile(sortedDurations, 0.90),
		P95Duration: calculatePercentile(sortedDurations, 0.95),
		P99Duration: calculatePercentile(sortedDurations, 0.99),
	}

	return stats
}

// calculatePercentile 计算百分位数
// sortedData 必须是已排序的切片
// percentile 取值范围 [0, 1]
func calculatePercentile(sortedData []time.Duration, percentile float64) time.Duration {
	if len(sortedData) == 0 {
		return 0
	}

	// 使用线性插值法计算百分位数
	index := percentile * float64(len(sortedData)-1)
	lower := int(index)
	upper := lower + 1

	if upper >= len(sortedData) {
		return sortedData[len(sortedData)-1]
	}

	// 线性插值
	fraction := index - float64(lower)
	return time.Duration(float64(sortedData[lower]) + fraction*float64(sortedData[upper]-sortedData[lower]))
}
