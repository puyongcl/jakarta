package listenerkey

import (
	"fmt"
	"jakarta/common/key/orderkey"
)

func GetChatOrderDescription(t, minute int64) string {
	switch t {
	case orderkey.ListenerOrderTypeTextChat:
		return fmt.Sprintf("%d分钟文字XX服务", minute)
	case orderkey.ListenerOrderTypeVoiceChat:
		return fmt.Sprintf("%d分钟语音XX服务", minute)
	default:

	}
	return ""
}
