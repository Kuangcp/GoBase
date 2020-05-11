package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

const Prefix = "keyboard:"
const LastInputEvent = Prefix + "last-event"

var info = cuibase.HelpInfo{
	Description: "Format markdown file, generate catalog",
	Version:     "1.0.0",
	VerbLen:     -3,
	ParamLen:    -9,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help info",
			Handler: nil,
		}, {
			Verb:    "-s",
			Param:   "<device>",
			Comment: "Listen keyboard with last device or specific device",
			Handler: ListenDevice,
		}, {
			Verb:    "-l",
			Param:   "",
			Comment: "List keyboard device",
			Handler: func(_ []string) {
				devices, _ := ListInputDevices()
				for _, dev := range devices {
					printKeyDevice(dev)
				}
			},
		}, {
			Verb:    "-ld",
			Param:   "",
			Comment: "List all device",
			Handler: func(_ []string) {
				devices, _ := ListInputDevices()
				for _, dev := range devices {
					fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
				}
			},
		}, {
			Verb:    "-p",
			Param:   "",
			Comment: "Print key map",
			Handler: PrintKeyMap,
		},
	}}

func main() {
	cuibase.RunActionFromInfo(info, nil)
}

func PrintKeyMap(params []string) {
	connection := initConnection()
	defer closeConnection(connection)

	event := ""
	if len(params) < 3 {
		last := connection.Get(LastInputEvent)
		if last == nil {
			return
		}
		event = last.Val()
	} else {
		event = params[2]
	}
	if event == "" {
		fmt.Printf("%vPlease select inputDevice %v\n", cuibase.Red, cuibase.End)
		return
	}

	device, _ := Open("/dev/input/" + event)
	if device == nil {
		return
	}

	fmt.Println(device)
	fmt.Printf("\n%vkey map:  %v", cuibase.LightGreen, cuibase.End)
	fmt.Println(device.Capabilities)
}

func printKeyDevice(dev *InputDevice) {
	for _, codes := range dev.Capabilities {
		for _, code := range codes {
			if code.Name == "KEY_F" {
				fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
				return
			}
		}
	}
}

func ListenDevice(params []string) {
	var event = ""
	if len(params) > 2 {
		event = params[2]
	}

	connection := initConnection()
	defer closeConnection(connection)

	if event == "" {
		last := connection.Get(LastInputEvent)
		if last == nil {
			return
		}
		event = last.Val()
	} else {
		connection.GetSet(LastInputEvent, event)
	}
	if event == "" {
		return
	}

	device, _ := Open("/dev/input/" + event)
	if device == nil {
		return
	}
	defer closeDevice(device)

	success := false
	for true {
		inputEvents, err := device.Read()
		if err != nil || inputEvents == nil || len(inputEvents) == 0 {
			continue
		}

		handleResult := handleEvents(inputEvents, connection)
		if !success && handleResult {
			success = handleResult
			fmt.Printf("%s  listen success. %s \n", cuibase.Green, cuibase.End)
		}
	}
}

func handleEvents(inputEvents []InputEvent, conn *redis.Client) bool {
	today := time.Now().Format("2006:01:02")
	flag := false
	for _, inputEvent := range inputEvents {
		if inputEvent.Code == 0 {
			continue
		}

		event := NewKeyEvent(&inputEvent)
		if event.State != KeyDown {
			continue
		}

		flag = true
		//fmt.Printf("%v           %v\n", event, inputEvent)
		conn.ZIncr(Prefix+today+":rank", redis.Z{Score: 1, Member: event.Scancode})
		conn.ZIncr(Prefix+"total", redis.Z{Score: 1, Member: today})

		// actual store us not ns
		conn.ZAdd(Prefix+today+":detail",
			redis.Z{Score: float64(event.Scancode), Member: inputEvent.Time.Nano() / 1000})
	}
	return flag
}

func initConnection() *redis.Client {
	target := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       5,
	})
	return target
}

func closeDevice(device *InputDevice) {
	err := device.Release()
	if err != nil {
		fmt.Println("release device error: ", err)
	}
}

func closeConnection(client *redis.Client) {
	err := client.Close()
	if err != nil {
		fmt.Println("close redis connection error: ", err)
	}
}
