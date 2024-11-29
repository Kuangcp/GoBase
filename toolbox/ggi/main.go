package main

import (
	"fmt"
)

func main() {
	info.Parse()
	if help {
		info.PrintHelp()
		return
	}

	if addRepo != "" {

		return
	}
	if delRepo != "" {
		return
	}

	if listRepo {

		return
	}

	fmt.Println(push, pull, allRepo)
	fmt.Println(addRepo, delRepo)
}
