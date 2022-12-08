package orderkey

// 服务类型
const (
	ListenerOrderTypeTextChat  = 2 // 文字聊天
	ListenerOrderTypeVoiceChat = 4 // 语音聊天
)

func GetPaySuccessOrderState(orderType int64) int64 {
	if orderType == ListenerOrderTypeVoiceChat {
		return ChatOrderStatePaySuccess3
	}
	return ChatOrderStateTextOrderPaySuccess24
}
