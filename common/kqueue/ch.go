package kqueue

// 上报事件
type UploadUserEventMessage struct {
	Uid   int64  `json:"uid"`
	Cb    string `json:"cb"`
	Event string `json:"event"`
	Value string `json:"value"`
	Stamp string `json:"stamp"`
}
