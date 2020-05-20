package main

import "time"

const Prefix = "keyboard:"
const LastInputEvent = Prefix + "last-event"
const TotalCount = Prefix + "total"
const KeyMap = Prefix+"key-map"

func GetRankKey(time time.Time) string {
	today := time.Format("2006:01:02")
	return Prefix + today + ":rank"
}

func GetDetailKey(time time.Time) string {
	today := time.Format("2006:01:02")
	return Prefix + today + ":detail"
}