package tool

import (
	"fmt"
	"time"
)

// 获取今日还剩多少秒
func GetTodayRemainSecond() int64 {
	//todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	//todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	now := time.Now()
	todayLastTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	return todayLastTime.Unix() - now.Unix()
}

// 获取今日开始和结束的时间
func GetTodayStartAndEndTime() (*time.Time, *time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	return &start, &end
}

// 获取昨日开始和结束的时间
func GetYesterdayStartAndEndTime() (*time.Time, *time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.Local)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return &start, &end
}

// 获取最近多少完整天的开始和结束的时间 不包括今天
func GetLastDayStartAndEndTime(d int) (*time.Time, *time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 0, -1)
	start = end.AddDate(0, 0, -d)
	return &start, &end
}

func GetTimeDurationText(td time.Duration) string {
	m := int64(td.Minutes())
	h := int64(td.Hours())
	d := h / 24

	m = m - h*60
	h = h - d*24

	var s string
	if d > 0 {
		s = fmt.Sprintf("%d天", d)
	} else if h > 0 {
		s = fmt.Sprintf("%d小时", h)
	} else if m > 0 {
		s = fmt.Sprintf("%d分钟", m)
	} else {
		s = "1分钟"
	}
	return s
}
