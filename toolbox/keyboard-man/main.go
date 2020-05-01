package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

const prefix = "keyboard:"
const lastEvent = prefix + "last-event"

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
		"-p": func(params []string) {
			cuibase.AssertParamCount(2, "input device")
			device, _ := Open("/dev/input/" + params[2])
			if device == nil {
				return
			}
			defer device.Release()

			fmt.Println(device)
			fmt.Println("key map:")
			fmt.Println(device.Capabilities)
		},
		"-s": ListenDevice,
	}, HelpInfo)
}

func ListenDevice(params []string) {
	var event = ""
	if len(params) > 2 {
		event = params[2]
	}

	connection := initConnection()
	defer connection.Close()

	if event == "" {
		last := connection.Get(lastEvent)
		if last == nil {
			return
		}
		event = last.Val()
	} else {
		connection.GetSet(lastEvent, event)
	}
	if event == "" {
		return
	}

	device, _ := Open("/dev/input/" + event)
	if device == nil {
		return
	}
	defer device.Release()

	for true {
		inputEvents, err := device.Read()
		if err != nil || inputEvents == nil || len(inputEvents) == 0 {
			continue
		}

		handleEvents(inputEvents, connection)
	}
}

func handleEvents(inputEvents []InputEvent, conn *redis.Client) {
	today := time.Now().Format("2006:01:02")
	for _, inputEvent := range inputEvents {
		if inputEvent.Code == 0 {
			continue
		}

		event := NewKeyEvent(&inputEvent)
		if event.State != KeyDown {
			continue
		}

		//fmt.Printf("%v           %v\n", event, inputEvent)
		conn.ZIncr(prefix+today+":rank", redis.Z{Score: 1, Member: event.Scancode})
		conn.ZIncr(prefix+"total", redis.Z{Score: 1, Member: today})

		// actual store us not ns
		conn.ZAdd(prefix+today+":detail", redis.Z{Score: float64(event.Scancode), Member: inputEvent.Time.Nano()})
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
				Comment: "Listen keyboard with last device or specific device",
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

func initConnection() *redis.Client {
	target := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       5,
	})
	return target
}
