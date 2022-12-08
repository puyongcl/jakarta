package kqueue

// 发送im自定义消息
type SendImDefineMessage struct {
	FromUid           int64  `json:"fromUid"`
	ToUid             int64  `json:"toUid"`
	MsgType           int64  `json:"msgType"`
	Title             string `json:"title"`
	Text              string `json:"text"`
	Val1              string `json:"val1"`
	Val2              string `json:"val2"`
	Val3              string `json:"val3,optional"`
	Val4              string `json:"val4,optional"`
	Val5              string `json:"val5,optional"`
	Val6              string `json:"val6,optional"`
	Sync              int64  `json:"sync"`
	RepeatMsgCheckId  string `json:"repeatMsgCheckId"`  // 去重检查
	RepeatMsgCheckSec int    `json:"repeatMsgCheckSec"` // 多少秒内不能发送第二条
}

// im回调 在线状态变更
type ImStateChangeMessage struct {
	Uid       int64 `json:"uid"`
	State     int64 `json:"state"`
	EventTime int64 `json:"eventTime"`
}

// 用户发送消息之后回调
type ImAfterSendMsg struct {
	MsgType string `json:"msgType"`
	FromUid string `json:"fromUid"`
	ToUid   string `json:"toUid"`
	Text    string `json:"text"`
	MsgSeq  int64  `json:"MsgSeq"`
	MsgTime int64  `json:"MsgTime"`
	MsgKey  string `json:"MsgKey"`
}
