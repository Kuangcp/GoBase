package store

import (
	"time"

	"github.com/kuangcp/gobase/pkg/ctk"
)

const (
	Prefix       = "keyboard:"
	DateFormat   = "2006:01:02"
	TimeFormat   = ctk.HH_MM_SS
	MsTimeFormat = ctk.HH_MM_SS_MS
)

const (
	LastInputEvent = Prefix + "last-event" // last use event. STRING
	TotalCount     = Prefix + "total"      // total count. ZSET
	KeyMap         = Prefix + "key-map"    // all key code map. HASH
	CoreLive       = Prefix + "core-live"  // core process live heart beat
)

// string max kpm in today
func GetTodayMaxKPMKey(time time.Time) string {
	return GetTodayMaxKPMKeyByString(time.Format(DateFormat))
}
func GetTodayTempKPMKey(time time.Time) string {
	return GetTodayTempKPMKeyByString(time.Format(DateFormat))
}

func GetTodayTempKPMKeyByString(timeStr string) string {
	return Prefix + timeStr + ":temp-kpm"
}

func GetTodayMaxKPMKeyByString(timeStr string) string {
	return Prefix + timeStr + ":kpm"
}

// GetRankKey by time
// zset member keyCode score 按键数
func GetRankKey(time time.Time) string {
	return GetRankKeyByString(time.Format(DateFormat))
}

// GetRankKeyByString
func GetRankKeyByString(timeStr string) string {
	return Prefix + timeStr + ":rank"
}

// GetDetailKey by time
// zset member 时间戳 score keyCode
func GetDetailKey(time time.Time) string {
	return GetDetailKeyByString(time.Format(DateFormat))
}

func GetDetailKeyByString(timeStr string) string {
	return Prefix + timeStr + ":detail"
}
