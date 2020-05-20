package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

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
		}, {
			Verb:    "-ca",
			Param:   "",
			Comment: "Cache key map",
			Handler: CacheKeyMap,
		}, {
			Verb:    "-d",
			Param:   "<x> <duration>",
			Comment: "Print daily total by before x day ago and duration",
			Handler: func(params []string) {
				printAnyByDay(params, printTotalByDate)
			},
		}, {
			Verb:    "-dr",
			Param:   "<x> <duration>",
			Comment: "Print daily rank by before x day ago and duration",
			Handler: func(params []string) {
				printAnyByDay(params, printRankByDate)
			},
		},
	}}

func main() {
	cuibase.RunActionFromInfo(info, nil)
}

func printAnyByDay(params []string, action func(time time.Time, conn *redis.Client)) {
	connection := initConnection()
	defer closeConnection(connection)

	now := time.Now()
	if len(params) == 3 {
		day, err := strconv.Atoi(params[2])
		cuibase.CheckIfError(err)
		now = now.AddDate(0, 0, - day)
	} else if len(params) == 4 {
		indexDay, err := strconv.Atoi(params[2])
		cuibase.CheckIfError(err)
		durationDay, err := strconv.Atoi(params[3])
		cuibase.CheckIfError(err)
		for i := 0; i < durationDay; i++ {
			action(now.AddDate(0, 0, -indexDay+i), connection)
		}
		return
	}
	action(now, connection)
}

func printRankByDate(time time.Time, conn *redis.Client) {
	today := time.Format("2006:01:02")

	all := conn.HGetAll(KeyMap)
	var keyMap map[string]string
	if all != nil {
		keyMap = all.Val()
	}
	totalScore := conn.ZScore(TotalCount, today)

	fmt.Printf("%s%-9s%s | %s | Total: %v\n", cuibase.Green, time.Weekday(), cuibase.End, today, int64(totalScore.Val()))
	score := conn.ZRevRangeByScoreWithScores(GetRankKey(time), redis.ZRangeBy{Min: "0", Max: "10000"})
	if len(keyMap) != 0 {
		for _, v := range score.Val() {
			fmt.Printf("%4v %v\n", v.Score, keyMap[v.Member.(string)][4:])
		}
	} else {
		for _, v := range score.Val() {
			fmt.Printf("%4v %v\n", v.Score, v.Member)
		}
	}
}

func printTotalByDate(time time.Time, conn *redis.Client) {
	today := time.Format("2006:01:02")
	score := conn.ZScore(TotalCount, today)
	fmt.Printf("%s%-9s%s %s %v\n", cuibase.Green, time.Weekday(), cuibase.End, today, int64(score.Val()))
}

func CacheKeyMap(params []string) {
	device := getDevice(params)
	if device == nil {
		return
	}
	_, codes := findActualBoardMap(device)
	if codes == nil {
		return
	}
	conn := initConnection()
	defer closeConnection(conn)
	for _, code := range codes {
		conn.HSet(KeyMap, strconv.Itoa(code.Code), code.Name)
		fmt.Printf("%v -> %v \n", code.Code, code.Name)
	}
}

func PrintKeyMap(params []string) {
	device := getDevice(params)
	if device == nil {
		return
	}

	fmt.Println(device)
	fmt.Printf("\n%vkey map:  %v", cuibase.LightGreen, cuibase.End)
	fmt.Println(device.Capabilities)
}

func getDevice(params []string) *InputDevice {
	connection := initConnection()
	defer closeConnection(connection)

	event := ""
	if len(params) < 3 {
		last := connection.Get(LastInputEvent)
		if last == nil {
			return nil
		}
		event = last.Val()
	} else {
		event = params[2]
	}
	if event == "" {
		fmt.Printf("%vPlease select inputDevice %v\n", cuibase.Red, cuibase.End)
		return nil
	}

	device, _ := Open("/dev/input/" + event)
	if device == nil {
		return nil
	}

	return device
}

func findActualBoardMap(dev *InputDevice) (*InputDevice, []CapabilityCode) {
	for _, codes := range dev.Capabilities {
		for _, code := range codes {
			if code.Name == "KEY_ESC" {
				return dev, codes
			}
		}
	}
	return nil, nil
}

func printKeyDevice(dev *InputDevice) {
	device, _ := findActualBoardMap(dev)
	if device != nil {
		fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
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
	today := time.Now()
	todayStr := today.Format("2006:01:02")
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
		conn.ZIncr(GetRankKey(today), redis.Z{Score: 1, Member: event.Scancode})
		conn.ZIncr(TotalCount, redis.Z{Score: 1, Member: todayStr})

		// actual store us not ns
		conn.ZAdd(GetDetailKey(today), redis.Z{Score: float64(event.Scancode), Member: inputEvent.Time.Nano() / 1000})
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
