package store

import "time"

const (
	Prefix     = "keyboard:"
	DateFormat = "2006:01:02"
	TimeFormat = "15:04:05"
)

const (
	LastInputEvent = Prefix + "last-event" // last use event. STRING
	TotalCount     = Prefix + "total"      // total count. ZSET
	KeyMap         = Prefix + "key-map"    // all key code map. HASH
)

// string max kpm in today
func GetTodayMaxKPMKey(time time.Time) string {
	return GetTodayMaxKPMKeyByString(time.Format(DateFormat))
}

func GetTodayTempKPMKeyByString(timeStr string) string {
	return Prefix + timeStr + ":temp-kpm"
}

func GetTodayMaxKPMKeyByString(timeStr string) string {
	return Prefix + timeStr + ":kpm"
}

//GetRankKey by time
// zset member keyCode score 按键数
func GetRankKey(time time.Time) string {
	return GetRankKeyByString(time.Format(DateFormat))
}

//GetRankKeyByString
func GetRankKeyByString(timeStr string) string {
	return Prefix + timeStr + ":rank"
}

//GetDetailKey by time
// zset member 时间戳 score keyCode
func GetDetailKey(time time.Time) string {
	return GetDetailKeyByString(time.Format(DateFormat))
}

func GetDetailKeyByString(timeStr string) string {
	return Prefix + timeStr + ":detail"
}
