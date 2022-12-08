package kqueue

// 小程序订阅消息
type SendMiniProgramSubscribeMessage struct {
	Thing4  string `json:"thing4"`
	Thing5  string `json:"thing5"`
	Time3   string `json:"time3"`
	Thing1  string `json:"thing1"`
	Thing3  string `json:"thing3"`
	Time2   string `json:"time2"`
	Date4   string `json:"date4"`
	MsgType int64  `json:"msgType"`
	ToUid   int64  `json:"toUid"` // 接收消息用户
	Page    string `json:"page"`
}
