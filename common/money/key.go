package money

import "fmt"

// 千分
const DivNumber = 1000

// 分 转成 元
const MoneyYuan = 100

func GetYuan(a int64) string {
	return fmt.Sprintf("%.2f", float64(a)/MoneyYuan)
}
