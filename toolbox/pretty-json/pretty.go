package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tidwall/pretty"
	"os"
)

var help bool
var uglyJSON string
var indent string

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&uglyJSON, "s", "", "json string")
	flag.StringVar(&indent, "i", "\t", "indent string, default tab")
}

func helpInfo() {
	fmt.Printf("usage:\n\n")
	flag.PrintDefaults()

	fmt.Println("\neg:\n 1. echo '{\"id\":1}' | pretty-json")
	fmt.Println(" 2. pretty-json -s '{\"id\":1}'")
}

func main() {
	flag.Parse()

	if help {
		helpInfo()
		return
	}

	if uglyJSON == "" {
		reader := bufio.NewReader(os.Stdin)
		result, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		uglyJSON = result
	}

	var Options = &pretty.Options{Width: 80, Prefix: "", Indent: indent, SortKeys: false}
	fmt.Printf("%s\n", pretty.Color(pretty.PrettyOptions([]byte(uglyJSON), Options), pretty.TerminalStyle))
}
