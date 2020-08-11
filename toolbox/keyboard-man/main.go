package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

var info = cuibase.HelpInfo{
	Description: "Record key input, show rank",
	Version:     "1.0.2",
	VerbLen:     -3,
	ParamLen:    -14,
	Params: []cuibase.ParamInfo{
		{
			Verb:    "-h",
			Param:   "",
			Comment: "Help info",
		}, {
			Verb:    "-l",
			Param:   "",
			Comment: "[root] List keyboard device",
		}, {
			Verb:    "-ld",
			Param:   "",
			Comment: "[root] List all device",
		}, {
			Verb:    "-p",
			Param:   "",
			Comment: "[root] Print key map",
		}, {
			Verb:    "-ca",
			Param:   "",
			Comment: "[root] Cache key map",
		}, {
			Verb:    "-s",
			Param:   "",
			Comment: "[root] Listen keyboard with last device or specific device",
		}, {
			Verb:    "-d",
			Param:   "",
			Comment: "Print daily total by before x day ago and duration",
		}, {
			Verb:    "-dr",
			Param:   "",
			Comment: "Print daily rank by before x day ago and duration",
		}, {
			Verb:    "-t",
			Param:   "<x>,<duration>",
			Comment: "Before x day ago and duration. For -d and -dr",
		}, {
			Verb:    "-e",
			Param:   "device",
			Comment: "Device. For -p -ca -s",
		},
	}}

var (
	help               bool
	printKeyMap        bool
	cacheKeyMap        bool
	listKeyboardDevice bool
	listAllDevice      bool
	listenDevice       bool
	day                bool
	dayRank            bool

	targetDevice string
	timePair     string
)

func init() {
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&printKeyMap, "p", false, "")
	flag.StringVar(&targetDevice, "e", "", "specific device")
	flag.BoolVar(&cacheKeyMap, "ca", false, "")
	flag.BoolVar(&listKeyboardDevice, "l", false, "")
	flag.BoolVar(&listAllDevice, "la", false, "")
	flag.BoolVar(&listenDevice, "s", false, "")
	flag.BoolVar(&day, "d", false, "")
	flag.BoolVar(&dayRank, "dr", false, "")
	flag.StringVar(&timePair, "t", "1", "")
}

func main() {
	flag.Parse()

	targetDevice = formatEvent(targetDevice)

	if help {
		info.PrintHelp()
		return
	}

	if listKeyboardDevice {
		devices, _ := ListInputDevices()
		for _, dev := range devices {
			printKeyDevice(dev)
		}
		return
	}

	if listAllDevice {
		devices, _ := ListInputDevices()
		for _, dev := range devices {
			fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
		}
		return
	}

	if cacheKeyMap {
		CacheKeyMap(targetDevice)
		return
	}

	if listenDevice {
		ListenDevice(targetDevice)
	}

	// simple query info

	if printKeyMap {
		PrintKeyMap(targetDevice)
	}

	if day {
		printDayInfoByDay(timePair, printTotalByDate)
	}

	if dayRank {
		printDayInfoByDay(timePair, printRankByDate)
	}
}

func printDayInfoByDay(timeSegment string, action func(time time.Time, conn *redis.Client)) {
	timePairs := strings.Split(timeSegment, ",")

	connection := initConnection()
	if connection == nil {
		return
	}
	defer closeConnection(connection)

	now := time.Now()
	indexDay := 0
	durationDay := 1
	if len(timePairs) == 1 {
		day, err := strconv.Atoi(timePairs[0])
		cuibase.CheckIfError(err)
		indexDay = day - 1
		durationDay = day
	} else if len(timePairs) == 2 {
		day, err := strconv.Atoi(timePairs[0])
		cuibase.CheckIfError(err)
		indexDay = day

		durationDay, err = strconv.Atoi(timePairs[1])
		cuibase.CheckIfError(err)
	}
	for i := 0; i < durationDay; i++ {
		action(now.AddDate(0, 0, -indexDay+i), connection)
	}
}

func printRankByDate(time time.Time, conn *redis.Client) {
	today := time.Format("2006:01:02")

	all := conn.HGetAll(KeyMap)
	var keyMap map[string]string
	if all != nil {
		keyMap = all.Val()
	}
	totalScore := conn.ZScore(TotalCount, today)

	fmt.Printf("%s | %s | Total: %v\n", cuibase.Green.Printf("%-8s", time.Weekday()),
		time.Format("2006-01-02"), cuibase.Yellow.Printf("%d", int64(totalScore.Val())))

	keyRank := conn.ZRevRangeByScoreWithScores(GetRankKey(time), redis.ZRangeBy{Min: "0", Max: "10000"})
	if len(keyMap) != 0 {
		var page []string
		row := len(keyRank.Val())/2 + 1
		for index, v := range keyRank.Val() {
			var d = index % row
			element := fmt.Sprintf("%4v → %-26v", v.Score, cuibase.LightGreen.Print(keyMap[v.Member.(string)]))

			if len(page) <= d {
				page = append(page, element)
			} else {
				page[d] = page[d] + element
			}
		}
		fmt.Println()
		for _, s := range page {
			fmt.Println(s)
		}
	} else {
		for _, v := range keyRank.Val() {
			fmt.Printf("%4v %v\n", v.Score, v.Member)
		}
	}
}

func printTotalByDate(time time.Time, conn *redis.Client) {
	today := time.Format("2006:01:02")
	score := conn.ZScore(TotalCount, today)
	fmt.Printf("%s %s %v\n", time.Format("2006-01-02"),
		cuibase.Green.Printf("%-9s", time.Weekday()), int64(score.Val()))
}

//CacheKeyMap to redis
func CacheKeyMap(targetDevice string) {
	device := getDevice(targetDevice)
	if device == nil {
		return
	}
	_, codes := findActualBoardMap(device)
	if codes == nil {
		return
	}
	conn := initConnection()
	if conn == nil {
		return
	}
	defer closeConnection(conn)
	for _, code := range codes {
		conn.HSet(KeyMap, strconv.Itoa(code.Code), code.Name[4:])
		fmt.Printf("%v -> %v \n", code.Code, code.Name)
	}
}

//PrintKeyMap show
func PrintKeyMap(targetDevice string) {
	device := getDevice(targetDevice)
	if device == nil {
		return
	}

	fmt.Println(device)
	fmt.Printf("\n%vkey map:  %v", cuibase.LightGreen, cuibase.End)
	fmt.Println(device.Capabilities)
}

func getDevice(targetDevice string) *InputDevice {
	connection := initConnection()
	if connection == nil {
		return nil
	}
	defer closeConnection(connection)

	event := ""
	if targetDevice == "" {
		last := connection.Get(LastInputEvent)
		if last == nil {
			return nil
		}
		if !strings.Contains(last.Val(), "event") {
			event = "event" + last.Val()
		} else {
			event = last.Val()
		}
	} else {
		event = targetDevice
	}
	if event == "" {
		fmt.Println(cuibase.Red.Print("Please select inputDevice"))
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

func formatEvent(input string) string {
	if !strings.Contains(input, "event") {
		return "event" + input
	}
	return input
}

// ListenDevice listen and record
func ListenDevice(targetDevice string) {
	connection := initConnection()
	if connection == nil {
		return
	}
	defer closeConnection(connection)

	if targetDevice == "" {
		last := connection.Get(LastInputEvent)
		if last == nil {
			return
		}
		targetDevice = last.Val()
	} else {
		connection.GetSet(LastInputEvent, targetDevice)
	}
	if targetDevice == "" {
		return
	}

	targetDevice = formatEvent(targetDevice)

	device, _ := Open("/dev/input/" + targetDevice)
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
			fmt.Println(cuibase.Green.Println("  listen success. "))
		}
	}
}

func handleEvents(inputEvents []InputEvent, conn *redis.Client) bool {
	today := time.Now()
	todayStr := today.Format("2006:01:02")
	matchFlag := false
	for _, inputEvent := range inputEvents {
		if inputEvent.Code == 0 {
			continue
		}

		event := NewKeyEvent(&inputEvent)
		if event.State != KeyDown {
			continue
		}

		matchFlag = true
		//fmt.Printf("%v           %v\n", event, inputEvent)
		conn.ZIncr(GetRankKey(today), redis.Z{Score: 1, Member: event.Scancode})
		conn.ZIncr(TotalCount, redis.Z{Score: 1, Member: todayStr})

		// actual store us not ns
		conn.ZAdd(GetDetailKey(today), redis.Z{Score: float64(event.Scancode), Member: inputEvent.Time.Nano() / 1000})
	}
	return matchFlag
}

func initConnection() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6667",
		Password: "",
		DB:       5,
	})

	result, err := conn.Ping().Result()
	if err != nil {
		fmt.Println(result, err)
		return nil
	}

	return conn
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
