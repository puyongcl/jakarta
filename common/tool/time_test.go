package tool

import (
	"fmt"
	"jakarta/common/dashboard"
	"jakarta/common/key/db"
	"testing"
	"time"
)

func TestGetTodayRemainSecond(t *testing.T) {
	got := GetTodayRemainSecond()
	fmt.Println(got / 60)
}

func TestAdd(t *testing.T) {
	now := time.Now()
	buyUnit := int64(10)
	addMin := time.Duration(buyUnit*15) * time.Minute
	fmt.Println(now.Add(addMin))
	fmt.Println(now.AddDate(0, 0, 100000))
}

func TestSub(t *testing.T) {
	now := time.Now()
	after := now.Add(34 * time.Hour)
	td := after.Sub(now)
	fmt.Println(td)
	m := td.Minutes()
	fmt.Println(m)
	fmt.Println(int(m))
	fmt.Println(td.Minutes())
}

func TestGetTodayStartAndEndTime(t *testing.T) {
	_, tmp := GetTodayStartAndEndTime()
	end := tmp.AddDate(0, 0, -1)
	fmt.Println(end)
	start := end.AddDate(0, 0, -dashboard.LastStat30Days)
	fmt.Println(start)
	fmt.Println(end.Sub(start).Hours() / 24)
}

func TestGetLastDayStartAndEndTime(t *testing.T) {
	a, b := GetLastDayStartAndEndTime(30)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(b.Sub(*a).Hours() / 24)
}

func TestLen(t *testing.T) {
	str := "系统通知"
	fmt.Println(len(str))

	testLen(nil)
}

func testLen(a []int64) {
	fmt.Println(len(a))
}

func TestGetTimeDurationText(t *testing.T) {
	got := GetTimeDurationText(10026 * time.Minute)
	fmt.Println(got)
}

func TestTimeAdd(t *testing.T) {
	t1 := time.Now()
	t2 := &t1
	t3 := *t2
	addTime(&t3)
	fmt.Println(t3)
}

func addTime(t *time.Time) {
	*t = t.Add(18 * time.Minute)
}

func TestDateFormat(t *testing.T) {
	tt := "20220905"
	t1, err := time.ParseInLocation(db.DateFormat2, tt, time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t1)
}

func TestTimeDuration(t *testing.T) {
	const t1 = 21 * 24 * 60 * 60
	t2 := t1 * time.Second
	fmt.Println(t1)
	fmt.Println(t2)
}

func TestTimeUnixMill(t *testing.T) {
	var um int64 = 1665976726538
	t3 := time.UnixMilli(um)
	fmt.Println(t3)
}
