package app

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/kuangcp/gobase/pkg/ctk"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/queue"
	"github.com/kuangcp/gobase/toolbox/keylogger/app/store"

	"github.com/kuangcp/logger"
	"github.com/manifoldco/promptui"

	"github.com/go-redis/redis"
	. "github.com/gvalkov/golang-evdev"
)

const (
	slideWindowMs        = 60_000                  // KPM 所以滑动窗口是一分钟
	calculateKPMPeriod   = time.Millisecond * 888  // 从KPM队列，计算得到最大KPM 操作的周期
	printKPMWindowPeriod = time.Millisecond * 4500 // 输出日志频率控制
)

var (
	targetDevice string
	timeSegment  string
	slideQueue   = queue.New() // key stroke timestamp(ms)
	currentKPM   = 0
)

func SetFormatTargetDevice(input string) {
	if input != "" {
		if !strings.Contains(input, "event") {
			targetDevice = "event" + input
		} else {
			targetDevice = input
		}
	}
}

func SetTimePair(timePair string) {
	timeSegment = timePair
}

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
func ListenDevice() {
	connection := store.GetConnection()
	if connection == nil {
		return
	}

	if targetDevice == "" {
		last := connection.Get(store.LastInputEvent)
		if last == nil {
			return
		}
		targetDevice = last.Val()
		if targetDevice == "" {
			return
		}
	}

	fmt.Println("Try to listen " + ctk.Yellow.Print(targetDevice) + " ...")

	device, err := Open("/dev/input/" + targetDevice)
	defer closeDevice(device)
	if device == nil || err != nil {
		logger.Error(err)
		return
	}

	go calculateKPM()

	hasSuccess := false
	for true {
		inputEvents, err := device.Read()
		if err != nil {
			logger.Error(err)
			return
		}

		if inputEvents == nil || len(inputEvents) == 0 {
			continue
		}

		handleResult := handleEvents(inputEvents)
		if !hasSuccess && handleResult {
			hasSuccess = true
			fmt.Println(ctk.Green.Print("\n    Listen success."))
			connection.Set(store.LastInputEvent, targetDevice, 0)
		}
	}
}

func calculateKPM() {
	conn := store.GetConnection()
	ticker := time.NewTicker(calculateKPMPeriod)
	for now := range ticker.C {
		day := now.Format(store.DateFormat)
		maxKPMKey := store.GetTodayMaxKPMKeyByString(day)

		// remove element that outer of time window
		nowMs := now.UnixNano() / 1000_000
		for {
			peek := slideQueue.Peek()
			if peek == nil {
				break
			}
			peekVal := (*peek).(int64)
			if nowMs-peekVal < slideWindowMs {
				break
			}
			slideQueue.Pop()
		}

		latestKPM := slideQueue.Len()
		if latestKPM == currentKPM {
			continue
		}
		// refresh redis current kpm
		currentKPM = latestKPM
		tempKPMKey := store.GetTodayTempKPMKeyByString(day)
		conn.Set(tempKPMKey, currentKPM, time.Hour*12)

		// init redis max kpm
		todayMaxKPM, err := conn.Get(maxKPMKey).Int()
		if err != nil {
			todayMaxKPM = 0
		}

		// set max kpm
		if todayMaxKPM >= currentKPM {
			continue
		}
		todayMaxKPM = currentKPM
		conn.Set(maxKPMKey, todayMaxKPM, 0)

		// delay log
		go func(tempKPM int) {
			time.Sleep(printKPMWindowPeriod)
			//fmt.Println("cache:", tempKPM, "now:", slideQueue.Len())
			if tempKPM >= slideQueue.Len() {
				logger.Info("Today max kpm up to", todayMaxKPM)
			}
		}(currentKPM)
	}
}

// ns us ms s
func handleEvents(inputEvents []InputEvent) bool {
	today := time.Now()
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
		store.IncrRankKey(today, event.Scancode)
		store.IncrTotalCount(today)

		keyNs := inputEvent.Time.Nano()
		store.AddKeyDetail(today, keyNs, event.Scancode)

		// push ms
		slideQueue.Push(keyNs / 1000_000)
	}
	return matchFlag
}

func OpenDevice() *InputDevice {
	event := ""
	if targetDevice == "" {
		last := store.GetConnection().Get(store.LastInputEvent)
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
		fmt.Println(ctk.Red.Print("Please select inputDevice"))
		return nil
	}

	device, _ := Open("/dev/input/" + event)
	if device == nil {
		return nil
	}

	return device
}

func PrintDay() {
	now := time.Now()
	indexDay, durationDay := parseTime(timeSegment)
	for i := 0; i < durationDay; i++ {
		handleTotalByDate(now.AddDate(0, 0, -indexDay+i), store.GetConnection())
	}
}

func PrintDayRank() {
	now := time.Now()
	indexDay, durationDay := parseTime(timeSegment)
	for i := 0; i < durationDay; i++ {
		handleRankByDate(now.AddDate(0, 0, -indexDay+i), store.GetConnection())
	}
}

func PrintTotalRank() {
	now := time.Now()
	indexDay, durationDay := parseTime(timeSegment)
	conn := store.GetConnection()
	all := conn.HGetAll(store.KeyMap)
	var keyMap map[string]string
	if all != nil {
		keyMap = all.Val()
	}

	result := make(map[string]float64)
	firstDay := now.AddDate(0, 0, -indexDay)
	lastDay := now.AddDate(0, 0, -indexDay+durationDay-1)
	for i := 0; i < durationDay; i++ {
		timeIndex := now.AddDate(0, 0, -indexDay+i)

		keyRank := conn.ZRevRangeByScoreWithScores(store.GetRankKey(timeIndex), redis.ZRangeBy{Min: "0", Max: "50000"})
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

	fmt.Printf("    %s → %s\n", firstDay.Format(ctk.YYYY_MM_DD), lastDay.Format(ctk.YYYY_MM_DD))

	if len(keyMap) != 0 {
		printByFourColumn(len(sortList), func(index int) string {
			val := sortList[index]
			return fmt.Sprintf("%7v → %-28v", val.Value, ctk.LightGreen.Print(keyMap[val.Key]))
		})
	} else {
		printByFourColumn(len(sortList), func(index int) string {
			val := sortList[index]
			return fmt.Sprintf("%7v → %-28v", val.Value, ctk.LightGreen.Print(val.Key))
		})
	}
}

func printByFourColumn(dataLen int, toString func(index int) string) {
	printByColumn(4, dataLen, toString)
}

// printByColumn 从上往下，从左往右 多列展示
func printByColumn(columnCount int, dataLen int, toString func(index int) string) {
	var lines []string
	row := dataLen/columnCount + 1
	for i := 0; i < dataLen; i++ {
		var lineIdx = i % row
		halfLine := toString(i)
		if halfLine == "" {
			continue
		}
		if len(lines) <= lineIdx {
			lines = append(lines, halfLine)
		} else {
			if lines != nil {
				lines[lineIdx] = lines[lineIdx] + halfLine
			}
		}
	}
	fmt.Println()
	for _, s := range lines {
		fmt.Println(s)
	}
}

func handleRankByDate(time time.Time, conn *redis.Client) {
	all := conn.HGetAll(store.KeyMap)
	var keyMap map[string]string
	if all != nil {
		keyMap = all.Val()
	}

	totalScore := store.TotalCountVal(time)
	maxKPM := store.MaxKPMVal(time)

	fmt.Printf("\n%s | %s | %-3s | Total: %s \n",
		ctk.Green.Printf("%-9s", time.Weekday()),
		time.Format(ctk.YYYY_MM_DD),
		ctk.Yellow.Printf("%3s", maxKPM),
		ctk.Green.Printf("%-5d", totalScore))

	keyRank := conn.ZRevRangeByScoreWithScores(store.GetRankKey(time), redis.ZRangeBy{Min: "0", Max: "50000"})
	if len(keyMap) != 0 {
		valList := keyRank.Val()
		printByFourColumn(len(valList), func(index int) string {
			val := valList[index]
			return fmt.Sprintf("%4v → %-26v", val.Score, ctk.LightGreen.Print(keyMap[val.Member.(string)]))
		})
	} else {
		valList := keyRank.Val()
		printByFourColumn(len(valList), func(index int) string {
			val := valList[index]
			return fmt.Sprintf("%4v → %-26v", val.Score, ctk.LightGreen.Print(val.Member.(string)))
		})
	}
}

func parseTime(timeSegment string) (int, int) {
	timePairs := strings.Split(timeSegment, ",")

	indexDay := 0
	durationDay := 1
	if len(timePairs) == 1 {
		day, err := strconv.Atoi(timePairs[0])
		ctk.CheckIfError(err)
		indexDay = day - 1
		durationDay = day
	} else if len(timePairs) == 2 {
		day, err := strconv.Atoi(timePairs[0])
		ctk.CheckIfError(err)
		indexDay = day

		durationDay, err = strconv.Atoi(timePairs[1])
		ctk.CheckIfError(err)
	}
	return indexDay, durationDay
}

func handleTotalByDate(time time.Time, conn *redis.Client) {
	today := time.Format(store.DateFormat)
	score := conn.ZScore(store.TotalCount, today)
	maxKPM := store.MaxKPMVal(time)
	fmt.Printf("%s %s %s %6v\n", time.Format(ctk.YYYY_MM_DD),
		ctk.Green.Printf("%-9s", time.Weekday()),
		ctk.Yellow.Printf("%4s", maxKPM),
		int64(score.Val()))
}

// CacheKeyMap to redis
func CacheKeyMap() {
	device := OpenDevice()
	if device == nil {
		return
	}
	_, codes := findActualBoardMap(device)
	if codes == nil {
		return
	}
	for _, code := range codes {
		store.GetConnection().HSet(store.KeyMap, strconv.Itoa(code.Code), code.Name[4:])
		fmt.Printf("%v -> %v \n", code.Code, code.Name)
	}
}

// PrintKeyMap show
func PrintKeyMap() {
	device := OpenDevice()
	if device == nil {
		return
	}

	fmt.Println(device)
	for capType, codes := range device.Capabilities {
		fmt.Printf("\n\n %s%v %v%s\n", ctk.Purple, capType.Type, capType.Name, ctk.End)
		printByColumn(6, len(codes), func(index int) string {
			if len(codes[index].Name) == 0 {
				return ""
			}
			return fmt.Sprintf("%s%4d%s %20s┃", ctk.LightGreen, codes[index].Code, ctk.End, codes[index].Name)
		})
	}
}

func findActualBoardMap(dev *InputDevice) (*InputDevice, []CapabilityCode) {
	for _, codes := range dev.Capabilities {
		for _, code := range codes {
			//logger.Info(dev.Fn, dev.Name, code)
			// 鼠标 键盘 均有该事件code
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
	devList := buildKeyBoardDeviceList()
	for i, device := range devList {
		peppers = append(peppers, option{
			Device:  device.Fn[11:],
			Desc:    device.Name + " | " + device.Phys,
			Peppers: i,
		})
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
		Size:      7,
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
