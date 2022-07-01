package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/tidwall/pretty"
)

var (
	debug    bool
	help     bool
	txt      bool
	uglyJSON string
	indent   string
)

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&uglyJSON, "s", "", "json string")
	flag.StringVar(&indent, "i", "    ", "indent string")
	flag.BoolVar(&debug, "d", false, "debug log")
	flag.BoolVar(&txt, "t", false, "no color")
}

func helpInfo() {
	fmt.Printf("usage:\n\n")
	flag.PrintDefaults()

	fmt.Println("\neg:\n   1. echo '{\"id\":1}' | pretty-json")
	fmt.Println("   2. pretty-json -s '{\"id\":1}'")
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
		if err != nil && debug {
			fmt.Printf("read error. result:%s error: %v\n", result, err)
		}
		uglyJSON = result
	}

	var Options = &pretty.Options{Width: 80, Prefix: "", Indent: indent, SortKeys: false}
	if txt {
		fmt.Println(string(pretty.PrettyOptions([]byte(uglyJSON), Options)))
	} else {
		fmt.Printf("%s\n", pretty.Color(pretty.PrettyOptions([]byte(uglyJSON), Options), pretty.TerminalStyle))
	}
}
