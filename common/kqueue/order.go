package kqueue

// 订单状态操作变更之后的后续操作
type UpdateChatOrderActionMessage struct {
	OrderId          string  `json:"orderId"`
	ListenerUid      int64   `json:"listenerUid"`
	ListenerNickName string  `json:"listenerNickName"`
	Uid              int64   `json:"uid"`
	NickName         string  `json:"nickName"`
	Action           int64   `json:"action"`
	PaidAmount       int64   `json:"paidAmount"`
	Star             int64   `json:"star"`
	OrderType        int64   `json:"orderType"`
	UsedMinute       int64   `json:"usedMinute"`
	BuyMinute        int64   `json:"buyMinute"`
	OrderCreateTime  string  `json:"orderCreateTime"` // 订单创建日期
	Reason           string  `json:"reason"`          // 退款理由
	SendMsg          int64   `json:"sendMsg"`
	Comment          string  `json:"comment"`        // 评价内容
	CommentTag       []int64 `json:"commentTag"`     // 评价标签
	AddRepeat        int64   `json:"addRepeat"`      // 复购人数加1
	AddUser          int64   `json:"addUser"`        // 服务人数加1
	TextExpireTime   string  `json:"textExpireTime"` // 文字订单过期时间
	UserChannel      string  `json:"userChannel"`
	ListenerChannel  string  `json:"listenerChannel"`
}
