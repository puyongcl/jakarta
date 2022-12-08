package tim

import (
	"encoding/json"
	"jakarta/common/xerr"
	"strconv"
)

const sendmsgUrl = "v4/openim/sendmsg"

// MsgType的值	类型
const (
	TIMTextElem      = "TIMTextElem"      // TIMTextElem	文本消息。
	TIMLocationElem  = "TIMLocationElem"  // TIMLocationElem	地理位置消息。
	TIMFaceElem      = "TIMFaceElem"      // TIMFaceElem	表情消息。
	TIMCustomElem    = "TIMCustomElem"    // TIMCustomElem	自定义消息，当接收方为 iOS 系统且应用处在后台时，此消息类型可携带除文本以外的字段到 APNs。一条组合消息中只能包含一个 TIMCustomElem 自定义消息元素。
	TIMSoundElem     = "TIMSoundElem"     // TIMSoundElem	语音消息。
	TIMImageElem     = "TIMImageElem"     // TIMImageElem	图像消息。
	TIMFileElem      = "TIMFileElem"      // TIMFileElem	文件消息。
	TIMVideoFileElem = "TIMVideoFileElem" // TIMVideoFileElem 视频消息
)

type DefineMsgPack struct {
	MsgType    string        `json:"MsgType"`    // 必填 消息类型
	MsgContent DefineMsgBody `json:"MsgContent"` // Object	必填	对于每种 MsgType 用不同的 MsgContent 格式，具体可参考 消息格式描述
}

type TextMsgPack struct {
	MsgType    string         `json:"MsgType"`    // 必填 消息类型
	MsgContent TextMsgContent `json:"MsgContent"` // Object	必填	对于每种 MsgType 用不同的 MsgContent 格式，具体可参考 消息格式描述
}

// 文本消息
type TextMsgContent struct {
	Text string `json:"Text"`
}

// 自定义消息
type DefineMsgContent struct {
	DefineMsgType int64  `json:"defineMsgType"` // 自定义消息类型
	Title         string `json:"title"`         // 消息标题
	Text          string `json:"text"`          // 消息内容
	Val1          string `json:"val_1"`         // 附带值 根据类型确定
	Val2          string `json:"val_2"`
	Val3          string `json:"val_3"` // 附带值 根据类型确定
	Val4          string `json:"val_4"`
	Val5          string `json:"val_5"`
	Val6          string `json:"val_6"`
}

// 是否把消息同步到发送端
const (
	TimMsgSyncFromYes = 1
	TimMsgSyncFromNo  = 2
)

// 消息体
type SendDefineTimMsgReq struct {
	SyncOtherMachine int64  `json:"SyncOtherMachine"` // 选填	1：把消息同步到 From_Account 在线终端和漫游上； 2：消息不同步至 From_Account； 若不填写默认情况下会将消息存 From_Account 漫游
	FromAccount      string `json:"From_Account"`     // 选填	消息发送方 UserID（用于指定发送消息方帐号）
	ToAccount        string `json:"To_Account"`       // 必填	消息接收方 UserID
	MsgLifeTime      int    `json:"MsgLifeTime"`      // 选填	消息离线保存时长（单位：秒），最长为7天（604800秒） 若设置该字段为0，则消息只发在线用户，不保存离线 若设置该字段超过7天（604800秒），仍只保存7天 若不设置该字段，则默认保存7天
	//MsgSeq                int       `json:"MsgSeq"`                // 选填	消息序列号（32位无符号整数），后台会根据该字段去重及进行同秒内消息的排序，详细规则请看本接口的功能说明。若不填该字段，则由后台填入随机数
	MsgRandom             uint32   `json:"MsgRandom"`             // 必填	消息随机数（32位无符号整数），后台用于同一秒内的消息去重。请确保该字段填的是随机
	ForbidCallbackControl []string `json:"ForbidCallbackControl"` //Array	选填	消息回调禁止开关，只对本条消息有效，ForbidBeforeSendMsgCallback 表示禁止发消息前回调，ForbidAfterSendMsgCallback 表示禁止发消息后回调
	//SendMsgControl        []string  `json:"SendMsgControl"`        //Array	选填	消息发送控制选项，是一个 String 数组，只对本条消息有效。"NoUnread"表示该条消息不计入未读数。"NoLastMsg"表示该条消息不更新会话列表。"WithMuteNotifications"表示该条消息的接收方对发送方设置的免打扰选项生效（默认不生效）。示例："SendMsgControl": ["NoUnread","NoLastMsg","WithMuteNotifications"]
	MsgBody []DefineMsgPack `json:"MsgBody"` //必填	消息内容，具体格式请参考 消息格式描述（注意，一条消息可包括多种消息元素，MsgBody 为 Array 类型）
}

// 消息体
type SendTextTimMsgReq struct {
	SyncOtherMachine int64  `json:"SyncOtherMachine"` // 选填	1：把消息同步到 From_Account 在线终端和漫游上； 2：消息不同步至 From_Account； 若不填写默认情况下会将消息存 From_Account 漫游
	FromAccount      string `json:"From_Account"`     // 选填	消息发送方 UserID（用于指定发送消息方帐号）
	ToAccount        string `json:"To_Account"`       // 必填	消息接收方 UserID
	MsgLifeTime      int    `json:"MsgLifeTime"`      // 选填	消息离线保存时长（单位：秒），最长为7天（604800秒） 若设置该字段为0，则消息只发在线用户，不保存离线 若设置该字段超过7天（604800秒），仍只保存7天 若不设置该字段，则默认保存7天
	//MsgSeq                int       `json:"MsgSeq"`                // 选填	消息序列号（32位无符号整数），后台会根据该字段去重及进行同秒内消息的排序，详细规则请看本接口的功能说明。若不填该字段，则由后台填入随机数
	MsgRandom             uint32   `json:"MsgRandom"`             // 必填	消息随机数（32位无符号整数），后台用于同一秒内的消息去重。请确保该字段填的是随机
	ForbidCallbackControl []string `json:"ForbidCallbackControl"` //Array	选填	消息回调禁止开关，只对本条消息有效，ForbidBeforeSendMsgCallback 表示禁止发消息前回调，ForbidAfterSendMsgCallback 表示禁止发消息后回调
	//SendMsgControl        []string  `json:"SendMsgControl"`        //Array	选填	消息发送控制选项，是一个 String 数组，只对本条消息有效。"NoUnread"表示该条消息不计入未读数。"NoLastMsg"表示该条消息不更新会话列表。"WithMuteNotifications"表示该条消息的接收方对发送方设置的免打扰选项生效（默认不生效）。示例："SendMsgControl": ["NoUnread","NoLastMsg","WithMuteNotifications"]
	MsgBody []TextMsgPack `json:"MsgBody"` //必填	消息内容，具体格式请参考 消息格式描述（注意，一条消息可包括多种消息元素，MsgBody 为 Array 类型）
}

type SendTimMsgResp struct {
	ActionStatus string `json:"ActionStatus"` //请求处理的结果，OK 表示处理成功，FAIL 表示失败
	ErrorCode    int64  `json:"ErrorCode"`    // 错误码，0表示成功，非0表示失败
	ErrorInfo    string `json:"ErrorInfo"`    //	错误信息
	MsgTime      int    `json:"MsgTime"`      //	消息时间戳，UNIX 时间戳
	MsgKey       string `json:"MsgKey"`       // 消息唯一标识，用于撤回。长度不超过50个字符
}

type DefineMsgBody struct {
	Data  string `json:"Data"`
	Desc  string `json:"Desc"`
	Ext   string `json:"Ext"`
	Sound string `json:"Sound"`
}

// 消息保存天数
const MsgKeepTime = 7 * 24 * 60 * 60

func (ra *RestApiClient) SendDefineMsg(fromUid, toUid int64, msgC *DefineMsgContent, sync int64) (rsp *SendTimMsgResp, err error) {
	buf, err := json.Marshal(msgC)
	if err != nil {
		return nil, err
	}
	msg := SendDefineTimMsgReq{
		SyncOtherMachine:      sync,
		FromAccount:           strconv.FormatInt(fromUid, 10),
		ToAccount:             strconv.FormatInt(toUid, 10),
		MsgLifeTime:           MsgKeepTime,
		MsgRandom:             ra.randInt.Uint32(),
		ForbidCallbackControl: []string{"ForbidBeforeSendMsgCallback", "ForbidAfterSendMsgCallback"},
		//SendMsgControl:        nil,
		MsgBody: []DefineMsgPack{
			{
				MsgType: TIMCustomElem,
				MsgContent: DefineMsgBody{
					Data: string(buf),
				},
			},
		},
	}

	rsp = &SendTimMsgResp{}
	resp, err := ra.Client.R().SetBody(&msg).SetResult(&rsp).Post(ra.GetTimRestApiUrl(sendmsgUrl))
	if err != nil {
		return
	}
	if resp.IsError() {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, resp.Status())
	}
	if rsp.ActionStatus == "OK" {
		return
	}
	return nil, xerr.NewGrpcErrCodeMsg(uint32(rsp.ErrorCode), rsp.ErrorInfo)
}

func (ra *RestApiClient) SendTextMsg(fromUid, toUid int64, text string, sync int64) (rsp *SendTimMsgResp, err error) {
	msg := SendTextTimMsgReq{
		SyncOtherMachine:      sync,
		FromAccount:           strconv.FormatInt(fromUid, 10),
		ToAccount:             strconv.FormatInt(toUid, 10),
		MsgLifeTime:           MsgKeepTime,
		MsgRandom:             ra.randInt.Uint32(),
		ForbidCallbackControl: []string{"ForbidBeforeSendMsgCallback", "ForbidAfterSendMsgCallback"},
		//SendMsgControl:        nil,
		MsgBody: []TextMsgPack{
			{
				MsgType: TIMTextElem,
				MsgContent: TextMsgContent{
					Text: text,
				},
			},
		},
	}

	rsp = &SendTimMsgResp{}
	resp, err := ra.Client.R().SetBody(&msg).SetResult(&rsp).Post(ra.GetTimRestApiUrl(sendmsgUrl))
	if err != nil {
		return
	}
	if resp.IsError() {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, resp.Status())
	}
	if rsp.ActionStatus == "OK" {
		return
	}
	return nil, xerr.NewGrpcErrCodeMsg(uint32(rsp.ErrorCode), rsp.ErrorInfo)
}
