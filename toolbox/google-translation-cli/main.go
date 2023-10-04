package main

import (
	"fmt"
	"github.com/mind1949/googletrans"
	"golang.org/x/text/language"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 0 {
		return
	}
	params := googletrans.TranslateParams{
		Src:  "auto",
		Dest: language.SimplifiedChinese.String(),
		Text: fmt.Sprintf("%v", args),
	}

	translated, err := googletrans.Translate(params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("text: %q \npronunciation: %q", translated.Text, translated.Pronunciation)
}
