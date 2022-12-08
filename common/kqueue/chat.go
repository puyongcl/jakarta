package kqueue

// 统计本次用户通话
type UpdateChatStatMessage struct {
	Uid         int64  `json:"uid"`
	ListenerUid int64  `json:"listenerUid"`
	LogId       string `json:"logId"`
	StartTime   string `json:"startTime"`
	StopTime    string `json:"stopTime"`
	OrderType   int64  `json:"orderType"`
}

type UpdateTextChatStatMessage struct {
	Uid         int64 `json:"uid"`
	ListenerUid int64 `json:"listenerUid"`
}

// 延迟检查聊天是否结束
type CheckChatStateMessage struct {
	Uid         int64 `json:"uid"`
	ListenerUid int64 `json:"listenerUid"`
	OrderType   int64 `json:"orderType"`
	DeferMinute int64 `json:"deferMinute"` // 延迟几分钟
}

// 用户首次进入聊天
type UserFirstEnterChatMessage struct {
	Uid         int64 `json:"uid"`
	ListenerUid int64 `json:"listenerUid"`
	IsFirst     int64 `json:"isFirst"`
}
