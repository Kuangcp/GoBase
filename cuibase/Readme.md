# cuibase
> 命令行工具 基础脚手架

## Quick Start

```go
func HelpInfo(_ []string) {
	info := cuibase.HelpInfo{
		Description: "",
		VerbLen:     -3,
		ParamLen:    -5,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "Help info",
			},
		}}
	cuibase.Help(info)
}

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h": HelpInfo,
	}, HelpInfo)
}
```