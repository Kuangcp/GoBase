package app

import "time"

//Prefix redis prefix
const Prefix = "keyboard:"

//LastInputEvent last use event
const LastInputEvent = Prefix + "last-event"

//TotalCount total count ZSET
const TotalCount = Prefix + "total"
const DateFormat = "2006:01:02"
const TimeFormat = "15:04:05"

//KeyMap cache key code map HASH
const KeyMap = Prefix + "key-map"

// string max kpm in today
func GetTodayMaxKPMKey(time time.Time) string {
	return GetTodayMaxKPMKeyByString(time.Format(DateFormat))
}

// string temp kpm in today
func GetTodayTempKPMKey(time time.Time) string {
	return GetTodayTempKPMKeyByString(time.Format(DateFormat))
}

func GetTodayTempKPMKeyByString(time string) string {
	return Prefix + time + ":temp-kpm"
}

func GetTodayMaxKPMKeyByString(time string) string {
	return Prefix + time + ":kpm"
}

//GetRankKey by time
// zset member keyCode score 按键数
func GetRankKey(time time.Time) string {
	return GetRankKeyByString(time.Format(DateFormat))
}

//GetRankKeyByString
func GetRankKeyByString(time string) string {
	return Prefix + time + ":rank"
}

//GetDetailKey by time
// zset member 时间戳 score keyCode
func GetDetailKey(time time.Time) string {
	return GetDetailKeyByString(time.Format(DateFormat))
}

func GetDetailKeyByString(time string) string {
	return Prefix + time + ":detail"
}
