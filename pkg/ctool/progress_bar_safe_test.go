package ctool

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
	"time"
)

// TestConcurrentProgressTracker_Basic 基本用法测试
func TestConcurrentProgressTracker_Basic(t *testing.T) {
	cpt := NewConcurrentProgressTracker()
	var totalNum int64 = 100
	cpt.SetTotal(totalNum)
	cpt.SetWidth(30)

	fmt.Println("=== 基本用法测试 ===")
	fmt.Println("初始状态:", cpt.String())

	// 使用 StartTask/FinishTask 来跟踪任务
	for i := 0; i < int(totalNum); i++ {
		//taskID := cpt.StartTask()
		//time.Sleep(time.Duration(rand.IntN(100)+60) * time.Millisecond)
		//cpt.FinishTask(taskID)
		//if i%3 == 0 {
		//	fmt.Printf("进度 %d: %s\n", i+1, cpt.String())
		//}

		fi := i
		cpt.RunTask(func() {
			time.Sleep(time.Duration(rand.IntN(100)+60) * time.Millisecond)
		}, func() {
			if i%3 == 2 {
				fmt.Printf(" %s 进度 %d\n", cpt.String(), fi+1)
			}
		})
	}

	fmt.Println("完成状态:", cpt.String())
	fmt.Printf("是否完成: %v\n", cpt.IsComplete())
	fmt.Printf("当前进度: %d/%d\n", cpt.GetCurrent(), cpt.GetTotal())
	fmt.Printf("百分比: %.2f%%\n", cpt.GetPercentage())
}

// TestConcurrentProgressTracker_Concurrent 并发场景测试
func TestConcurrentProgressTracker_Concurrent(t *testing.T) {
	cpt := NewConcurrentProgressTracker()
	totalTasks := int64(4000)
	cpt.SetTotal(totalTasks)
	cpt.SetWidth(40)
	cpt.SetMaxSamples(int(totalTasks))

	fmt.Println("\n=== 并发场景测试 ===")
	fmt.Println("初始状态:", cpt.String())
	watch := NewStopWatch()
	watch.Start("start")

	var wg sync.WaitGroup
	workerCount := 10 // 10个并发worker

	// 启动多个goroutine并发执行任务
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			// 每个worker执行多个任务
			tasksPerWorker := int(totalTasks) / workerCount
			for j := 0; j < tasksPerWorker; j++ {
				taskID := cpt.StartTask()

				// 模拟不同耗时的任务
				duration := time.Duration(rand.IntN(100)+30) * time.Millisecond
				time.Sleep(duration)

				//CalculatePiTylor(5000, 6000)

				cpt.FinishTask(taskID)
			}
		}(i)
	}

	// 定期打印进度
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				current := cpt.GetCurrent()
				if current < totalTasks {
					fmt.Printf("进度: %s (活跃任务: %d)\n",
						cpt.String(), cpt.GetActiveTaskCount())
				}
			case <-done:
				return
			}
		}
	}()

	// 等待所有任务完成
	wg.Wait()
	close(done)
	watch.Stop()

	fmt.Println("\n最终状态:", cpt.String())
	fmt.Printf("是否完成: %v\n%v", cpt.IsComplete(), watch.PrettyPrint())

	// 显示任务统计信息
	stats := cpt.GetTaskStatsDetail()
	fmt.Printf("\n任务统计信息:\n")
	fmt.Printf("样本数量: %d\n", stats.SampleCount)
	fmt.Printf("平均耗时: %v\n", stats.AvgDuration)
	fmt.Printf("最小耗时: %v\n", stats.MinDuration)
	fmt.Printf("P30 耗时: %v\n", stats.P30Duration)
	fmt.Printf("P50 耗时: %v\n", stats.P50Duration)
	fmt.Printf("P75 耗时: %v\n", stats.P75Duration)
	fmt.Printf("P90 耗时: %v\n", stats.P90Duration)
	fmt.Printf("P95 耗时: %v\n", stats.P95Duration)
	fmt.Printf("P99 耗时: %v\n", stats.P99Duration)
	fmt.Printf("最大耗时: %v\n", stats.MaxDuration)
}

// TestConcurrentProgressTracker_Update 使用 Update 方法（向后兼容）
func TestConcurrentProgressTracker_Update(t *testing.T) {
	cpt := NewConcurrentProgressTracker()
	cpt.SetTotal(20)
	cpt.SetWidth(30)

	fmt.Println("\n=== Update 方法测试（不记录任务耗时）===")
	fmt.Println("初始状态:", cpt.String())

	// 使用 Update 方法批量更新
	for i := 0; i < 4; i++ {
		cpt.Update(5) // 每次增加5
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("更新后: %s\n", cpt.String())
	}

	fmt.Println("完成状态:", cpt.String())
}

// TestConcurrentProgressTracker_Mixed 混合使用 StartTask/FinishTask 和 Update
func TestConcurrentProgressTracker_Mixed(t *testing.T) {
	cpt := NewConcurrentProgressTracker()
	cpt.SetTotal(30)

	fmt.Println("\n=== 混合使用测试 ===")

	// 一部分使用 StartTask/FinishTask
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			taskID := cpt.StartTask()
			time.Sleep(time.Duration(rand.IntN(50)+20) * time.Millisecond)
			cpt.FinishTask(taskID)
		}()
	}

	// 另一部分使用 Update
	for i := 0; i < 20; i++ {
		cpt.Update(1)
		time.Sleep(10 * time.Millisecond)
	}

	wg.Wait()

	fmt.Println("最终状态:", cpt.String())
	fmt.Printf("当前进度: %d/%d\n", cpt.GetCurrent(), cpt.GetTotal())
}

// TestConcurrentProgressTracker_Reset 重置测试
func TestConcurrentProgressTracker_Reset(t *testing.T) {
	cpt := NewConcurrentProgressTracker()
	cpt.SetTotal(10)

	fmt.Println("\n=== 重置测试 ===")

	// 执行一些任务
	for i := 0; i < 5; i++ {
		taskID := cpt.StartTask()
		time.Sleep(20 * time.Millisecond)
		cpt.FinishTask(taskID)
	}

	fmt.Println("重置前:", cpt.String())

	// 重置
	cpt.Reset()
	cpt.SetTotal(10)

	fmt.Println("重置后:", cpt.String())

	// 重新执行任务
	for i := 0; i < 10; i++ {
		taskID := cpt.StartTask()
		time.Sleep(20 * time.Millisecond)
		cpt.FinishTask(taskID)
	}

	fmt.Println("完成后:", cpt.String())
}

// TestConcurrentProgressTracker_EdgeCases 边界情况测试
func TestConcurrentProgressTracker_EdgeCases(t *testing.T) {
	cpt := NewConcurrentProgressTracker()

	fmt.Println("\n=== 边界情况测试 ===")

	// 测试总数为0
	cpt.SetTotal(0)
	fmt.Println("总数为0:", cpt.String())

	// 测试负数
	cpt.SetTotal(100)
	cpt.Update(-10) // 应该被忽略
	fmt.Println("负数delta:", cpt.String())

	// 测试超出总数
	cpt.SetCurrent(150)
	fmt.Println("超出总数:", cpt.String())
	fmt.Printf("当前值: %d, 总数: %d\n", cpt.GetCurrent(), cpt.GetTotal())

	// 测试无效任务ID
	cpt.FinishTask(999) // 不存在的任务ID
	fmt.Println("无效任务ID后:", cpt.String())
}

// ExampleConcurrentProgressTracker 示例：展示典型使用场景
func ExampleConcurrentProgressTracker() {
	// 创建一个进度跟踪器
	tracker := NewConcurrentProgressTracker()
	tracker.SetTotal(100)
	tracker.SetWidth(40)

	// 模拟并发处理100个任务
	var wg sync.WaitGroup
	workerCount := 5

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				// 开始任务
				taskID := tracker.StartTask()

				// 模拟任务执行
				time.Sleep(time.Duration(rand.IntN(50)+20) * time.Millisecond)

				// 完成任务
				tracker.FinishTask(taskID)
			}
		}()
	}

	// 定期打印进度
	ticker := time.NewTicker(300 * time.Millisecond)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				if !tracker.IsComplete() {
					fmt.Println(tracker.String())
				}
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	ticker.Stop()
	close(done)

	fmt.Println(tracker.String())
}

// TestConcurrentProgressTracker_GetTaskStatsDetail 测试详细统计信息（包含百分位数）
func TestConcurrentProgressTracker_GetTaskStatsDetail(t *testing.T) {
	tracker := NewConcurrentProgressTracker()
	tracker.SetTotal(100)
	tracker.SetMaxSamples(100)

	// 模拟不同耗时的任务
	durations := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
		40 * time.Millisecond,
		50 * time.Millisecond,
		60 * time.Millisecond,
		70 * time.Millisecond,
		80 * time.Millisecond,
		90 * time.Millisecond,
		100 * time.Millisecond,
	}

	// 执行任务
	for _, d := range durations {
		taskID := tracker.StartTask()
		time.Sleep(d)
		tracker.FinishTask(taskID)
	}

	// 获取详细统计信息
	stats := tracker.GetTaskStatsDetail()

	fmt.Printf("\n任务统计信息:\n")
	fmt.Printf("样本数量: %d\n", stats.SampleCount)
	fmt.Printf("平均耗时: %v\n", stats.AvgDuration)
	fmt.Printf("最小耗时: %v\n", stats.MinDuration)
	fmt.Printf("P30 耗时: %v\n", stats.P30Duration)
	fmt.Printf("P50 耗时: %v\n", stats.P50Duration)
	fmt.Printf("P75 耗时: %v\n", stats.P75Duration)
	fmt.Printf("P90 耗时: %v\n", stats.P90Duration)
	fmt.Printf("P95 耗时: %v\n", stats.P95Duration)
	fmt.Printf("P99 耗时: %v\n", stats.P99Duration)
	fmt.Printf("最大耗时: %v\n", stats.MaxDuration)

	// 验证基本正确性
	if stats.SampleCount != 10 {
		t.Errorf("期望样本数量为 10，实际为 %d", stats.SampleCount)
	}

	// 验证顺序关系：最小 <= P30 <= P50 <= P75 <= P90 <= P95 <= P99 <= 最大
	if stats.MinDuration > stats.P30Duration {
		t.Errorf("最小耗时 %v 应该 <= P30耗时 %v", stats.MinDuration, stats.P30Duration)
	}
	if stats.P30Duration > stats.P50Duration {
		t.Errorf("P30耗时 %v 应该 <= P50耗时 %v", stats.P30Duration, stats.P50Duration)
	}
	if stats.P50Duration > stats.P75Duration {
		t.Errorf("P50耗时 %v 应该 <= P75耗时 %v", stats.P50Duration, stats.P75Duration)
	}
	if stats.P75Duration > stats.P90Duration {
		t.Errorf("P75耗时 %v 应该 <= P90耗时 %v", stats.P75Duration, stats.P90Duration)
	}
	if stats.P90Duration > stats.P95Duration {
		t.Errorf("P90耗时 %v 应该 <= P95耗时 %v", stats.P90Duration, stats.P95Duration)
	}
	if stats.P95Duration > stats.P99Duration {
		t.Errorf("P95耗时 %v 应该 <= P99耗时 %v", stats.P95Duration, stats.P99Duration)
	}
	if stats.P99Duration > stats.MaxDuration {
		t.Errorf("P99耗时 %v 应该 <= 最大耗时 %v", stats.P99Duration, stats.MaxDuration)
	}

	// 验证平均值在合理范围内
	if stats.AvgDuration < stats.MinDuration || stats.AvgDuration > stats.MaxDuration {
		t.Errorf("平均耗时 %v 应该在 [%v, %v] 范围内", stats.AvgDuration, stats.MinDuration, stats.MaxDuration)
	}
}

// TestConcurrentProgressTracker_BackwardCompatibility 测试向后兼容性
func TestConcurrentProgressTracker_BackwardCompatibility(t *testing.T) {
	tracker := NewConcurrentProgressTracker()
	tracker.SetTotal(10)

	// 执行一些任务
	for i := 0; i < 5; i++ {
		taskID := tracker.StartTask()
		time.Sleep(10 * time.Millisecond)
		tracker.FinishTask(taskID)
	}

	// 使用旧版本的 GetTaskStats
	avgOld, minOld, maxOld, countOld := tracker.GetTaskStats()

	// 使用新版本的 GetTaskStatsDetail
	statsNew := tracker.GetTaskStatsDetail()

	// 验证两者结果一致
	if avgOld != statsNew.AvgDuration {
		t.Errorf("平均耗时不一致: 旧版本 %v, 新版本 %v", avgOld, statsNew.AvgDuration)
	}
	if minOld != statsNew.MinDuration {
		t.Errorf("最小耗时不一致: 旧版本 %v, 新版本 %v", minOld, statsNew.MinDuration)
	}
	if maxOld != statsNew.MaxDuration {
		t.Errorf("最大耗时不一致: 旧版本 %v, 新版本 %v", maxOld, statsNew.MaxDuration)
	}
	if countOld != statsNew.SampleCount {
		t.Errorf("样本数量不一致: 旧版本 %d, 新版本 %d", countOld, statsNew.SampleCount)
	}
}
