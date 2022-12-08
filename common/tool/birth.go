package tool

import (
	"time"
)

const DateLayout = "2006-01-02"

var constellation = []string{
	"水瓶座",
	"双鱼座",
	"白羊座",
	"金牛座",
	"双子座",
	"巨蟹座",
	"狮子座",
	"处女座",
	"天秤座",
	"天蝎座",
	"射手座",
	"魔羯座",
}

func GetTimeFromStrDate(date string) (year, month, day int) {
	d, err := time.Parse(DateLayout, date)
	if err != nil {
		return 0, 0, 0
	}
	year = d.Year()
	month = int(d.Month())
	day = d.Day()
	return
}

func GetZodiac(date string) (zodiac string) {
	d, err := time.Parse(DateLayout, date)
	if err != nil {
		return ""
	}
	year := d.Year()
	if year <= 0 {
		zodiac = "-1"
	}
	start := 1901
	x := (start - year) % 12
	if x < 0 {
		x += 12
	}

	switch x {
	case 1:
		zodiac = "鼠"
	case 0:
		zodiac = "牛"
	case 11:
		zodiac = "虎"
	case 10:
		zodiac = "兔"
	case 9:
		zodiac = "龙"
	case 8:
		zodiac = "蛇"
	case 7:
		zodiac = "马"
	case 6:
		zodiac = "羊"
	case 5:
		zodiac = "猴"
	case 4:
		zodiac = "鸡"
	case 3:
		zodiac = "狗"
	case 2:
		zodiac = "猪"
	}
	return
}

func GetAge(date string) (age int) {
	d, err := time.Parse(DateLayout, date)
	if err != nil {
		return 0
	}
	year := d.Year()
	if year <= 0 {
		age = -1
	}
	currentYear := time.Now().Year()
	age = currentYear - year
	return
}

func GetAge2(date time.Time) (age int64) {
	year := date.Year()
	if year <= 0 {
		return -1
	}
	currentYear := time.Now().Year()
	age = int64(currentYear - year)
	return
}

func GetConstellation(date string) (idx int64, star string) {
	d, err := time.Parse(DateLayout, date)
	if err != nil {
		return -1, ""
	}
	month := d.Month()
	day := d.Day()
	switch {
	case month <= 0, month >= 13, day <= 0, day >= 32:
		return -1, ""
	case month == 1 && day >= 20, month == 2 && day <= 18:
		return 1, "水瓶座"
	case month == 2 && day >= 19, month == 3 && day <= 20:
		return 2, "双鱼座"
	case month == 3 && day >= 21, month == 4 && day <= 19:
		return 3, "白羊座"
	case month == 4 && day >= 20, month == 5 && day <= 20:
		return 4, "金牛座"
	case month == 5 && day >= 21, month == 6 && day <= 21:
		return 5, "双子座"
	case month == 6 && day >= 22, month == 7 && day <= 22:
		return 6, "巨蟹座"
	case month == 7 && day >= 23, month == 8 && day <= 22:
		return 7, "狮子座"
	case month == 8 && day >= 23, month == 9 && day <= 22:
		return 8, "处女座"
	case month == 9 && day >= 23, month == 10 && day <= 22:
		return 9, "天秤座"
	case month == 10 && day >= 23, month == 11 && day <= 21:
		return 10, "天蝎座"
	case month == 11 && day >= 22, month == 12 && day <= 21:
		return 11, "射手座"
	case month == 12 && day >= 22, month == 1 && day <= 19:
		return 12, "魔蝎座"
	}
	return
}
