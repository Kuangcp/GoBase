package app

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/wonderivan/logger"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
	"github.com/kuangcp/gobase/pkg/cuibase"
)

func closeDevice(device *InputDevice) {
	if device == nil {
		return
	}
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
		if targetDevice == "" {
			return
		}
	}

	targetDevice = FormatEvent(targetDevice)
	fmt.Println("Try to listen " + cuibase.Yellow.Print(targetDevice) + " ...")

	device, _ := Open("/dev/input/" + targetDevice)
	defer closeDevice(device)
	if device == nil {
		log.Println("device not exist")
		return
	}

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
			fmt.Println(cuibase.Green.Print("\n    Listen success."))
			connection.Set(LastInputEvent, targetDevice, 0)
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
	todayStr := today.Format(DateFormat)
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
	now := time.Now()
	indexDay, durationDay := parseTime(timeSegment)
	for i := 0; i < durationDay; i++ {
		handleTotalByDate(now.AddDate(0, 0, -indexDay+i), GetConnection())
	}
}

func PrintDayRank(timeSegment string) {
	now := time.Now()
	indexDay, durationDay := parseTime(timeSegment)
	for i := 0; i < durationDay; i++ {
		handleRankByDate(now.AddDate(0, 0, -indexDay+i), GetConnection())
	}
}

func PrintTotalRank(timeSegment string) {
	now := time.Now()
	indexDay, durationDay := parseTime(timeSegment)
	conn := GetConnection()
	all := conn.HGetAll(KeyMap)
	var keyMap map[string]string
	if all != nil {
		keyMap = all.Val()
	}

	result := make(map[string]float64)
	firstDay := now.AddDate(0, 0, -indexDay)
	lastDay := now.AddDate(0, 0, -indexDay+durationDay-1)
	for i := 0; i < durationDay; i++ {
		timeIndex := now.AddDate(0, 0, -indexDay+i)

		keyRank := conn.ZRevRangeByScoreWithScores(GetRankKey(timeIndex), redis.ZRangeBy{Min: "0", Max: "50000"})
		for _, v := range keyRank.Val() {
			keyCode := v.Member.(string)
			val, ok := result[keyCode]
			if ok {
				result[keyCode] = val + v.Score
			} else {
				result[keyCode] = v.Score
			}
		}
	}
	type kv struct {
		Key   string
		Value float64
	}

	var sortList []kv
	for k, v := range result {
		sortList = append(sortList, kv{k, v})
	}

	sort.Slice(sortList, func(i, j int) bool {
		return sortList[i].Value > sortList[j].Value // 降序
	})

	fmt.Printf("    %s → %s\n", firstDay.Format("2006-01-02"), lastDay.Format("2006-01-02"))

	if len(keyMap) != 0 {
		printWithTwoColumn(len(sortList), func(index int) string {
			val := sortList[index]
			return fmt.Sprintf("%7v → %-28v", val.Value, cuibase.LightGreen.Print(keyMap[val.Key]))
		})
	} else {
		printWithTwoColumn(len(sortList), func(index int) string {
			val := sortList[index]
			return fmt.Sprintf("%7v → %-28v", val.Value, cuibase.LightGreen.Print(val.Key))
		})
	}
}

func printWithTwoColumn(dataLen int, toString func(index int) string) {
	var lines []string
	row := dataLen/2 + 1
	for i := 0; i < dataLen; i++ {
		var lineIdx = i % row
		halfLine := toString(i)
		if len(lines) <= lineIdx {
			lines = append(lines, halfLine)
		} else {
			lines[lineIdx] = lines[lineIdx] + halfLine
		}
	}
	fmt.Println()
	for _, s := range lines {
		fmt.Println(s)
	}
}

func handleRankByDate(time time.Time, conn *redis.Client) {
	today := time.Format(DateFormat)

	all := conn.HGetAll(KeyMap)
	var keyMap map[string]string
	if all != nil {
		keyMap = all.Val()
	}
	totalScore := conn.ZScore(TotalCount, today)

	fmt.Printf("%s | %s | Total: %v\n", cuibase.Green.Printf("%-8s", time.Weekday()),
		time.Format("2006-01-02"), cuibase.Yellow.Printf("%d", int64(totalScore.Val())))

	keyRank := conn.ZRevRangeByScoreWithScores(GetRankKey(time), redis.ZRangeBy{Min: "0", Max: "50000"})
	if len(keyMap) != 0 {
		valList := keyRank.Val()
		printWithTwoColumn(len(valList), func(index int) string {
			val := valList[index]
			return fmt.Sprintf("%4v → %-26v", val.Score, cuibase.LightGreen.Print(keyMap[val.Member.(string)]))
		})
	} else {
		valList := keyRank.Val()
		printWithTwoColumn(len(valList), func(index int) string {
			val := valList[index]
			return fmt.Sprintf("%4v → %-26v", val.Score, cuibase.LightGreen.Print(val.Member.(string)))
		})
	}
}

func parseTime(timeSegment string) (int, int) {
	timePairs := strings.Split(timeSegment, ",")

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
	return indexDay, durationDay
}

func handleTotalByDate(time time.Time, conn *redis.Client) {
	today := time.Format(DateFormat)
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
	list := buildKeyBoardDeviceList()
	for _, dev := range list {
		fmt.Printf("%s %s %s\n", dev.Fn, dev.Name, dev.Phys)
	}
}

func SelectDevice() (string, error) {
	type option struct {
		Device  string
		Desc    string
		Peppers int
	}

	var peppers []option
	dev := buildKeyBoardDeviceList()
	for i, device := range dev {
		peppers = append(peppers, option{Device: device.Fn[11:], Desc: device.Name, Peppers: i})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "🔓{{ .Device | green }} {{ .Desc }}",
		Inactive: "  {{ .Device | cyan }} {{ .Desc }}",
		Selected: "🔒️{{ .Device | green | cyan }}",
	}

	searcher := func(input string, index int) bool {
		pepper := peppers[index]
		name := strings.Replace(strings.ToLower(pepper.Device), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	promptSelect := promptui.Select{
		Label:     "Listen which device",
		Items:     peppers,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := promptSelect.Run()
	if err != nil {
		log.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return peppers[i].Device, nil
}

func buildKeyBoardDeviceList() []*InputDevice {
	var result []*InputDevice
	devices, _ := ListInputDevices()
	for _, dev := range devices {
		device, _ := findActualBoardMap(dev)
		if device != nil {
			result = append(result, dev)
		}
	}
	return result
}
