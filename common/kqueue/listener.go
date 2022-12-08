package kqueue

type UpdateListenerUserStatMessage struct {
	Time        string  `json:"time"`
	Event       int64   `json:"event"`
	Uid         int64   `json:"uid"`
	NickName    string  `json:"nickName"`
	Avatar      string  `json:"avatar"`
	ListenerUid []int64 `json:"listenerUid"`
}

type SendHelloWhenUserLoginMessage struct {
	Uid       int64 `json:"uid"`
	IsNewUser int64 `json:"isNewUser"`
}
