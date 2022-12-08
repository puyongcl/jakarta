package tim

import (
	"context"
	"encoding/json"
	"jakarta/common/key/userkey"
	"jakarta/common/kqueue"
	tim2 "jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"strconv"

	"jakarta/app/im/api/internal/svc"
	"jakarta/app/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackLogic) StateChangeCallback(req *types.TIMCallbackStateChangeReq) (resp *types.TIMCallbackResp, err error) {
	uid, err := strconv.ParseInt(req.Info.ToAccount, 10, 64)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("StateChangeCallback error uid req:%+v", req)
		return nil, err
	}
	// Login 表示上线（TCP 建立），Logout 表示下线（TCP 断开），Disconnect 表示网络断开（TCP 断开）
	val := kqueue.ImStateChangeMessage{
		Uid:       uid,
		State:     getOnlineStatus(req.Info.Action),
		EventTime: req.EventTime,
	}
	buf, err := json.Marshal(&val)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("StateChangeCallback json marshal error req:%+v", req)
		return nil, err
	}

	err = l.svcCtx.KqueueImStateChangeMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("StateChangeCallback kafka push error req:%+v", req)
		return nil, err
	}
	resp = &types.TIMCallbackResp{}
	logx.WithContext(l.ctx).Infof("StateChangeCallback kafka push req:%+v", req)
	return
}

func getOnlineStatus(event string) int64 {
	switch event {
	case tim2.Login:
		return userkey.Login
	case tim2.Logout:
		return userkey.Logout
	case tim2.Disconnect:
		return userkey.Disconnect
	default:
		return 0
	}
}

func (l *CallbackLogic) AfterSendMsgCallback(req *types.TIMCallbackAfterSendMsgReq) (resp *types.TIMCallbackResp, err error) {
	resp = &types.TIMCallbackResp{}
	if req.SendMsgResult != 0 { // 没有发送成功
		return
	}
	if len(req.MsgBody) <= 0 {
		return
	}
	if !tool.IsStringArrayExist(req.MsgBody[0].MsgType, []string{tim2.TIMTextElem, tim2.TIMSoundElem, tim2.TIMFaceElem, tim2.TIMImageElem}) { // 文本消息 语音消息 表情消息
		return
	}
	if req.FromAccount == "" || req.ToAccount == "" {
		return
	}

	val := kqueue.ImAfterSendMsg{
		MsgType: req.MsgBody[0].MsgType,
		FromUid: req.FromAccount,
		ToUid:   req.ToAccount,
		Text:    req.MsgBody[0].MsgContent.Text,
		MsgSeq:  int64(req.MsgSeq),
		MsgTime: int64(req.MsgTime),
		MsgKey:  req.MsgKey,
	}
	buf, err := json.Marshal(&val)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AfterSendMsgCallback json marshal error req:%+v", req)
		return nil, err
	}

	err = l.svcCtx.KqueueImAfterSendMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AfterSendMsgCallback kafka push error req:%+v", req)
		return nil, err
	}
	return
}
