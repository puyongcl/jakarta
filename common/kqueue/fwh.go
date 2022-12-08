package kqueue

// 服务号事件
type WxFwhCallbackEventMessage struct {
	OpenId string `json:"openId"` // 服务号openid
	Event  string `json:"event"`  // 2 关注 4 取消关注
}

// 服务号模版消息
type SendFwhNotifyMessage struct {
	First    string `json:"first"` // 标题
	Keyword1 string `json:"keyword1"`
	Keyword2 string `json:"keyword2"`
	Keyword3 string `json:"keyword3"`
	Keyword4 string `json:"keyword4"`
	Remark   string `json:"remark"` // 提示
	MsgType  int64  `json:"msgType"`
	ToUid    int64  `json:"toUid"` // 接收消息用户
	Color    string `json:"color"`
	Path     string `json:"path"`
}
