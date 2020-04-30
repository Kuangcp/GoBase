package main

import (
	"fmt"
	"strings"

	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

// TODO store in redis
func Listen(event string) {
	if event == "" {
		return
	}

	device, _ := Open("/dev/input/" + event)
	defer device.Release()

	//fmt.Println(device)
	//fmt.Println(device.Capabilities)

	for true {
		inputEvents, err := device.Read()
		if err != nil || inputEvents == nil || len(inputEvents) == 0 {
			continue
		}

		for _, inputEvent := range inputEvents {
			event := NewKeyEvent(&inputEvent)
			if event.State != KeyDown {
				continue
			}
			println(event.String())
		}
	}
}

func HelpInfo(_ []string) {
	info := cuibase.HelpInfo{
		Description: "Format markdown file, generate catalog",
		VerbLen:     -3,
		ParamLen:    -5,
		Params: []cuibase.ParamInfo{
			{
				Verb:    "-h",
				Param:   "",
				Comment: "Help info",
			}, {
				Verb:    "-s",
				Param:   "<device>",
				Comment: "Listen keyboard with last input device or specific device",
			}, {
				Verb:    "-l",
				Param:   "",
				Comment: "List keyboard device",
			}, {
				Verb:    "-ld",
				Param:   "",
				Comment: "List all device",
			},
		}}
	cuibase.Help(info)
}

func main() {
	cuibase.RunAction(map[string]func(params []string){
		"-h": HelpInfo,
		"-ld": func(_ []string) {
			devices, _ := ListInputDevices()
			for _, dev := range devices {
				fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
			}
		},
		"-l": func(_ []string) {
			devices, _ := ListInputDevices()
			for _, dev := range devices {
				if strings.Contains(dev.Name, "Keyboard") {
					fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
				}
			}
		},
		"-s": func(params []string) {
			if len(params) > 1 {
				Listen(params[2])
			} else {
				Listen("")
			}
		},
	}, HelpInfo)
}
