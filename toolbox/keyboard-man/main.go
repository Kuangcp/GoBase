package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

var lastEvent = "last-event"

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
			if len(params) > 2 {
				ListenDevice(params[2])
			} else {
				ListenDevice("")
			}
		},
	}, HelpInfo)
}

func ListenDevice(event string) {
	connection := initConnection()
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
	defer device.Release()

	//fmt.Println(device)
	//fmt.Println(device.Capabilities)

	today := time.Now().Format("2006-01-02")

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
			//fmt.Printf("%v \n", event)
			connection.ZIncr(today, redis.Z{Score: 1, Member: event.Scancode})
			connection.ZIncr("total", redis.Z{Score: 1, Member: today})
			connection.ZAdd("detail-"+today, redis.Z{Score: float64(inputEvent.Time.Nano()), Member: event.Scancode})
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
