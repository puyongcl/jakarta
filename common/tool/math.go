package tool

// 绝对值
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

// 除法
func DivideInt64(a int64, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a / b
}
