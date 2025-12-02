package ctool

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func TestProgressBar(t *testing.T) {
	// 创建进度条对象
	pb := NewProgressBar()

	// 设置总个数
	pb.SetTotal(100)

	// 测试初始状态
	fmt.Println("初始状态:")
	fmt.Println(pb.String())
	fmt.Println()

	// 模拟进度更新
	for i := 0; i < 100; i++ {
		// 更新进度，每次增加10
		pb.Update(1)
		if i%10 == 0 {
			//fmt.Printf("进度更新 %d: %s\n", i+1, pb.String())
			fmt.Printf("%s\n", pb.String())
		}

		al := rand.IntN(100) + 20
		time.Sleep(time.Duration(al) * time.Millisecond) // 模拟处理时间
	}

	fmt.Println("---------------------")
	// 测试完成状态
	fmt.Println("\n完成状态:")
	fmt.Println(pb.String())
	fmt.Printf("是否完成: %v\n", pb.IsComplete())
}

func TestProgressBarWithDelta(t *testing.T) {
	pb := NewProgressBar()
	pb.SetTotal(50)
	pb.SetWidth(30)

	fmt.Println("测试delta更新:")
	fmt.Println("初始:", pb.String())

	// 使用不同的delta值更新
	pb.Update(5)
	fmt.Println("+5:  ", pb.String())

	pb.Update(10)
	fmt.Println("+10: ", pb.String())

	pb.Update(15)
	fmt.Println("+15: ", pb.String())

	pb.Update(20)
	fmt.Println("+20: ", pb.String())
}

func TestProgressBarEdgeCases(t *testing.T) {
	// 测试边界情况
	pb := NewProgressBar()

	// 测试总数为0
	pb.SetTotal(0)
	fmt.Println("总数为0:", pb.String())

	// 测试负数
	pb.SetTotal(100)
	pb.Update(-50) // 应该不会变成负数
	fmt.Println("负数delta:", pb.String())

	// 测试超出总数
	pb.SetTotal(100)
	pb.SetCurrent(150) // 应该被限制为100
	fmt.Println("超出总数:", pb.String())
	fmt.Printf("当前值: %d, 总数: %d\n", pb.GetCurrent(), pb.GetTotal())
}
