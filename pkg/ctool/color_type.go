package ctool

import "fmt"

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
	End         Color = "\033[0m"
)

// Print content with color
func (t Color) Print(content string) string {
	return string(t) + content + string(End)
}

// Print content with color no end
func (t Color) PrintNoEnd(content string) string {
	return string(t) + content
}

// Print content with color
func (t Color) Println(content string) string {
	return string(t) + content + string(End) + "\n"
}

// Printf content with color and format
func (t Color) Printf(format string, a ...interface{}) string {
	return string(t) + fmt.Sprintf(format, a...) + string(End)
}

func (t Color) String() string {
	return string(t)
}
