package kqueue

// 用户订阅消息
type SubscribeNotifyMsgMessage struct {
	Uid       int64                            `json:"uid"`             // 订阅消息者uid
	TargetUid int64                            `json:"targetUid"`       // 订阅消息对象Uid
	MsgType   int64                            `json:"msgType"`         // 消息类型
	SendCnt   int64                            `json:"sendCnt"`         // 发送次数 1 1次 2 永久订阅
	Action    int64                            `json:"action"`          // 见订阅消息事件类型
	IMMsg     *SendImDefineMessage             `json:"imMsg,omitempty"` // 自定义消息
	MpMsg     *SendMiniProgramSubscribeMessage `json:"mpMsg,omitempty"` // 小程序通知消息
}
