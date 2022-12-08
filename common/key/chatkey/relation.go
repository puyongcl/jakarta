package chatkey

// 用户和XXX互动加分
const (
	InteractiveEventTypeSendMsg1           = 1 // 用户发送消息
	InteractiveEventTypeSendMsg2           = 2 // XXX发送消息
	InteractiveEventTypePayOrder3          = 3 // 用户下文字订单
	InteractiveEventTypePayOrder4          = 4 // 用户下语音订单
	InteractiveEventTypeCommentOrder5Star5 = 5 // 评价订单满意
	InteractiveEventTypeCommentOrder3Star6 = 6 // 评价订单一般
	InteractiveEventTypeCommentOrder1Star7 = 7 // 评价订单不满意
)

func GetAddScore(eventType int64) int64 {
	switch eventType {
	case InteractiveEventTypeSendMsg1, InteractiveEventTypeSendMsg2:
		return 1
	case InteractiveEventTypePayOrder3, InteractiveEventTypePayOrder4:
		return 200
	case InteractiveEventTypeCommentOrder5Star5:
		return 50
	case InteractiveEventTypeCommentOrder3Star6:
		return 30
	case InteractiveEventTypeCommentOrder1Star7:
		return 0
	default:
		return 0
	}
}
