package db

import "time"

const DateTimeFormat = "2006-01-02 15:04:05"
const DateTimeFormat2 = "2006-01-02 15:04"
const DateFormat = "2006-01-02"
const DateFormat2 = "20060102"

func ParseTimeString(s string) *time.Time {
	t, _ := time.ParseInLocation(DateTimeFormat, s, time.Local)
	return &t
}
