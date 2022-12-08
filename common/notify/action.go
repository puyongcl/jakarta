package notify

// 订阅消息事件类型
const (
	SubscribeUserNotifyMsgEventAdd     = 1 // 订阅某用户的某一种消息
	SubscribeUserNotifyMsgEventCancel  = 2 // 取消订阅用户的某一种消息
	SubscribeUserNotifyMsgEventSend    = 3 // 给用户的订阅者发送消息
	SubscribeOneTimeNotifyMsgEventAdd  = 4 // 订阅一次性的某种消息
	SubscribeOneTimeNotifyMsgEventSend = 5 // 发送订阅一次性的某种消息
)
