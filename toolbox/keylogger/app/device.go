package app

import (
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/cuibase"
)

func closeDevice(device *InputDevice) {
	err := device.Release()
	if err != nil {
		fmt.Println("release device error: ", err)
	}
}

// ListenDevice listen and record
func ListenDevice(targetDevice string) {
	if connection == nil {
		return
	}

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

	fmt.Println("try listen", targetDevice)
	targetDevice = FormatEvent(targetDevice)

	device, _ := Open("/dev/input/" + targetDevice)
	if device == nil {
		return
	}
	defer closeDevice(device)

	success := false
	for true {
		inputEvents, err := device.Read()
		if err != nil {
			logger.Error(err)
			return
		}

		if inputEvents == nil || len(inputEvents) == 0 {
			continue
		}

		handleResult := handleEvents(inputEvents, connection)
		if !success && handleResult {
			success = handleResult
			fmt.Println(cuibase.Green.Println("listen success. "))
		}
	}
}

func FormatEvent(input string) string {
	if input != "" && !strings.Contains(input, "event") {
		return "event" + input
	}
	return input
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
		result, err := conn.ZIncr(GetRankKey(today), redis.Z{Score: 1, Member: event.Scancode}).Result()
		if err != nil {
			fmt.Println("zincr: ", result, err)
			CloseConnection()
			os.Exit(1)
		}
		result, err = conn.ZIncr(TotalCount, redis.Z{Score: 1, Member: todayStr}).Result()
		if err != nil {
			fmt.Println("zincr: ", result, err)
			CloseConnection()
			os.Exit(1)
		}
		// actual store us not ns
		var num int64 = 0
		num, err = conn.ZAdd(GetDetailKey(today), redis.Z{Score: float64(event.Scancode), Member: inputEvent.Time.Nano() / 1000}).Result()
		if err != nil {
			fmt.Println("zadd: ", num, err)
			CloseConnection()
			os.Exit(1)
		}
	}
	return matchFlag
}

func OpenDevice(targetDevice string) *InputDevice {
	event := ""
	if targetDevice == "" {
		last := GetConnection().Get(LastInputEvent)
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

func PrintDay(timeSegment string) {
	printDayInfoByDay(timeSegment, handleTotalByDate)
}

func PrintDayRank(timeSegment string) {
	printDayInfoByDay(timeSegment, handleRankByDate)
}

func printDayInfoByDay(timeSegment string, action func(time time.Time, conn *redis.Client)) {
	timePairs := strings.Split(timeSegment, ",")

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
		action(now.AddDate(0, 0, -indexDay+i), GetConnection())
	}
}

func handleRankByDate(time time.Time, conn *redis.Client) {
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

func handleTotalByDate(time time.Time, conn *redis.Client) {
	today := time.Format("2006:01:02")
	score := conn.ZScore(TotalCount, today)
	fmt.Printf("%s %s %v\n", time.Format("2006-01-02"),
		cuibase.Green.Printf("%-9s", time.Weekday()), int64(score.Val()))
}

//CacheKeyMap to redis
func CacheKeyMap(targetDevice string) {
	device := OpenDevice(targetDevice)
	if device == nil {
		return
	}
	_, codes := findActualBoardMap(device)
	if codes == nil {
		return
	}
	for _, code := range codes {
		GetConnection().HSet(KeyMap, strconv.Itoa(code.Code), code.Name[4:])
		fmt.Printf("%v -> %v \n", code.Code, code.Name)
	}
}

//PrintKeyMap show
func PrintKeyMap(targetDevice string) {
	device := OpenDevice(targetDevice)
	if device == nil {
		return
	}

	fmt.Println(device)
	fmt.Printf("\n%vkey map:  %v", cuibase.LightGreen, cuibase.End)
	fmt.Println(device.Capabilities)
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

func ListAllDevice() {
	devices, _ := ListInputDevices()
	for _, dev := range devices {
		fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
	}
}

func ListAllKeyBoardDevice() {
	devices, _ := ListInputDevices()
	for _, dev := range devices {
		device, _ := findActualBoardMap(dev)
		if device != nil {
			fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
		}
	}
}
