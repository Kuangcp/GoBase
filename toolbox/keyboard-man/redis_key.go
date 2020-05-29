package main

import "time"

//Prefix redis prefix
const Prefix = "keyboard:"

//LastInputEvent last use event
const LastInputEvent = Prefix + "last-event"

//TotalCount total count ZSET
const TotalCount = Prefix + "total"

//KeyMap cache key code map HASH
const KeyMap = Prefix+"key-map"

//GetRankKey by time
func GetRankKey(time time.Time) string {
	today := time.Format("2006:01:02")
	return Prefix + today + ":rank"
}

//GetDetailKey by time
func GetDetailKey(time time.Time) string {
	today := time.Format("2006:01:02")
	return Prefix + today + ":detail"
}