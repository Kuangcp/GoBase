package cuibase

import "fmt"

var End = "\033[0m"

type Color string

const (
	Red         Color = "\033[0;31m"
	Green       Color = "\033[0;32m"
	Yellow      Color = "\033[0;33m"
	Blue        Color = "\033[0;34m"
	Purple      Color = "\033[0;35m"
	Cyan        Color = "\033[0;36m"
	White       Color = "\033[0;37m"
	LightRed    Color = "\033[0;91m"
	LightGreen  Color = "\033[0;92m"
	LightYellow Color = "\033[0;93m"
	LightBlue   Color = "\033[0;94m"
	LightPurple Color = "\033[0;95m"
	LightCyan   Color = "\033[0;96m"
	LightWhite  Color = "\033[0;97m"
)

// Print content with color
func (t Color) Print(content string) string {
	return string(t) + content + End
}

// Printf content with color and format
func (t Color) Printf(format string, content string) string {
	return string(t) + fmt.Sprintf(format, content) + End
}