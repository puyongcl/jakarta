package userkey

// 在线状态
// Login 表示上线（TCP 建立），Logout 表示下线（TCP 断开），Disconnect 表示网络断开（TCP 断开）
const (
	Unknown    = 0
	Login      = 1
	Logout     = 2
	Disconnect = 3
)
