package prime

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

func findTwoPrimeFactor(target int) int {
	pairCount := 0
	for i := 2; i < int(math.Sqrt(float64(target)))+1; i++ {
		temp := target % i
		if temp == 0 {
			fmt.Printf("%11d * %11d = %11d \n", i, target/i, target)
			pairCount++
		}
	}
	if pairCount == 0 {
		fmt.Println("prime: ", target)
	}
	return pairCount
}

func TestOne(t *testing.T) {
	findTwoPrimeFactor(7140229933)
}

func TestWithChannel(t *testing.T) {
	var latch sync.WaitGroup
	latch.Add(1000)
	for i := 0; i < 1000; i++ {
		j := i
		go func() {
			target := 6541367000 + j
			findTwoPrimeFactor(target)
			latch.Done()
		}()
	}
	latch.Wait()
}

func TestInterval(t *testing.T) {
	total := 0
	for i := 0; i < 1000; i++ {
		target := 6541367000 + i
		count := findTwoPrimeFactor(target)
		total += count
	}
	fmt.Println(total)
}
