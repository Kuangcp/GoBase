package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestId(t *testing.T) {
	cache := make(map[string]int)
	rand.Seed(time.Now().UnixNano())
	for j := 0; j < 100000; j++ {
		for i := 0; i < 5; i++ {
			finalID := fmt.Sprint(i) + "_" + uuid.New().String()[24:]
			cache[finalID] = 0
			if j == 4 && i == 0 {

				fmt.Println(finalID)

			}
		}
	}
	fmt.Println(len(cache))


}
