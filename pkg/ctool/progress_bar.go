package ctool

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar 进度条结构体
type ProgressBar struct {
	total      int64     // 总个数
	current    int64     // 当前进度
	startTime  time.Time // 开始时间
	lastUpdate time.Time // 上次更新时间
	width      int       // 进度条宽度（字符数）
}
type Task struct {
	startTime time.Time // 开始时间
	endTime   time.Time // 上次更新时间
}

// NewProgressBar 创建新的进度条对象
func NewProgressBar() *ProgressBar {
	now := time.Now()
	return &ProgressBar{
		total:      0,
		current:    0,
		startTime:  now,
		lastUpdate: now,
		width:      50, // 默认宽度50个字符
	}
}

// SetTotal 设置总个数
func (pb *ProgressBar) SetTotal(total int64) {
	pb.total = total
	if pb.total < 0 {
		pb.total = 0
	}
}

// SetWidth 设置进度条宽度（可选方法）
func (pb *ProgressBar) SetWidth(width int) {
	if width > 0 {
		pb.width = width
	}
}

// Update 更新当前进度，增加delta值
func (pb *ProgressBar) Update(delta int64) {
	pb.current += delta
	if pb.current < 0 {
		pb.current = 0
	}
	if pb.total > 0 && pb.current > pb.total {
		pb.current = pb.total
	}
	pb.lastUpdate = time.Now()
}

// SetCurrent 直接设置当前进度值（可选方法）
func (pb *ProgressBar) SetCurrent(current int64) {
	pb.current = current
	if pb.current < 0 {
		pb.current = 0
	}
	if pb.total > 0 && pb.current > pb.total {
		pb.current = pb.total
	}
	pb.lastUpdate = time.Now()
}

// String 拼接进度条字符串，类似pacman风格
// 格式: [#########>        ] 45% (123/456) 12.34s elapsed, 15.67s remaining
func (pb *ProgressBar) String() string {
	if pb.total <= 0 {
		return fmt.Sprintf("[%s] 0%% (0/0) calculating...", strings.Repeat(" ", pb.width))
	}

	// 计算百分比
	percentage := float64(pb.current) / float64(pb.total) * 100
	if percentage > 100 {
		percentage = 100
	}

	// 计算已执行时间
	elapsed := time.Since(pb.startTime)
	elapsedSeconds := elapsed.Seconds()

	// 计算剩余时间
	var remainingSeconds float64
	if pb.current > 0 {
		avgTimePerItem := elapsedSeconds / float64(pb.current)
		remainingItems := pb.total - pb.current
		remainingSeconds = avgTimePerItem * float64(remainingItems)
	} else {
		remainingSeconds = 0
	}

	// 构建进度条
	filled := int(float64(pb.width) * percentage / 100)
	if filled > pb.width {
		filled = pb.width
	}
	empty := pb.width - filled

	bar := strings.Repeat("#", filled)
	if filled < pb.width {
		bar += ">"
		empty--
		if empty > 0 {
			bar += strings.Repeat(" ", empty)
		}
	}

	// 格式化时间
	elapsedStr := formatDuration(elapsedSeconds)
	remainingStr := formatDuration(remainingSeconds)

	return fmt.Sprintf("[%s] %.1f%% (%d/%d) %s elapsed, %s remaining",
		bar, percentage, pb.current, pb.total, elapsedStr, remainingStr)
}

// formatDuration 格式化秒数为可读的时间字符串
func formatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.2fs", seconds)
	} else if seconds < 3600 {
		minutes := int(seconds / 60)
		secs := seconds - float64(minutes*60)
		return fmt.Sprintf("%dm %.2fs", minutes, secs)
	} else {
		hours := int(seconds / 3600)
		remaining := seconds - float64(hours*3600)
		minutes := int(remaining / 60)
		secs := remaining - float64(minutes*60)
		return fmt.Sprintf("%dh %dm %.2fs", hours, minutes, secs)
	}
}

// GetPercentage 获取当前进度百分比
func (pb *ProgressBar) GetPercentage() float64 {
	if pb.total <= 0 {
		return 0
	}
	percentage := float64(pb.current) / float64(pb.total) * 100
	if percentage > 100 {
		return 100
	}
	return percentage
}

// GetCurrent 获取当前进度值
func (pb *ProgressBar) GetCurrent() int64 {
	return pb.current
}

// GetTotal 获取总个数
func (pb *ProgressBar) GetTotal() int64 {
	return pb.total
}

// IsComplete 检查是否已完成
func (pb *ProgressBar) IsComplete() bool {
	return pb.total > 0 && pb.current >= pb.total
}

// Reset 重置进度条
func (pb *ProgressBar) Reset() {
	pb.current = 0
	pb.startTime = time.Now()
	pb.lastUpdate = pb.startTime
}
