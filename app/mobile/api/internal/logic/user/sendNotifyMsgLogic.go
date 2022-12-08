package user

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
)

type SendNotifyMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendNotifyMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNotifyMsgLogic {
	return &SendNotifyMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendNotifyMsgLogic) SendNotifyMsg(req *types.SendNotifyMsgReq) (resp *types.SendNotifyMsgResp, err error) {
	//var kqMsg *kqueue.SendImDefineMessage
	//switch req.MsgType {
	//case notify.DefineMsgTypeViewMsg17:
	//	return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "error msg type")
	//	kqMsg = &kqueue.SendImDefineMessage{
	//		FromUid:           notify.TimViewNotifyUid,
	//		ToUid:             req.ListenerUid,
	//		MsgType:           notify.DefineMsgTypeViewMsg17,
	//		Title:             notify.DefineNotifyMsgTemplateViewMsgTitle17,
	//		Text:              fmt.Sprintf(notify.DefineNotifyMsgTemplateViewMsg17, req.NickName),
	//		Val1:              strconv.FormatInt(req.Uid, 10),
	//		Val2:              req.Avatar,
	//		Sync:              tim.TimMsgSyncFromNo,
	//		RepeatMsgCheckId:  fmt.Sprintf("%d-%d", notify.TimOrderNotifyUid, req.ListenerUid),
	//		RepeatMsgCheckSec: notify.SendViewNotifyLimitMin * 60,
	//	}
	//default:
	//	return nil, xerr.NewErrCodeMsg(xerr.RequestParamError, "wrong msg type")
	//}
	//
	//buf, err := json.Marshal(&kqMsg)
	//if err != nil {
	//	return nil, err
	//}
	//err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	//if err != nil {
	//	return nil, err
	//}
	resp = &types.SendNotifyMsgResp{}
	return
}
