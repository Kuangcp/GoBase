package main

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {

	param := queryParam{
		appId:     "1",
		query:     "2有",
		secretKey: "3",
	}

	url := param.buildFinalURL()

	fmt.Printf("%v", url)
}
