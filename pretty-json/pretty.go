package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tidwall/pretty"
	"os"
)

var help bool
var uglyJson string
var indent string

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&uglyJson, "s", "", "json string")
	flag.StringVar(&indent, "i", "\t", "indent string, default tab")
}

func main() {
	flag.Parse()

	if help {
		fmt.Printf("usage:\n\n")
		flag.PrintDefaults()

		fmt.Println("\neg:\n echo '{\"id\":1}' | pretty_json")
		return
	}

	if uglyJson == "" {
		reader := bufio.NewReader(os.Stdin)
		result, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read error:", err)
			return
		}
		uglyJson = result
	}

	var Options = &pretty.Options{Width: 80, Prefix: "", Indent: indent, SortKeys: false}
	fmt.Printf("%s\n", pretty.Color(pretty.PrettyOptions([]byte(uglyJson), Options), pretty.TerminalStyle))
}
