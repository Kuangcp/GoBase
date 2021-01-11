package main

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {

	param := QueryParam{
		app:       "1",
		query:     "2",
		secretKey: "3",
	}

	url := param.buildFinalURL()

	fmt.Printf("%v", url)
}
