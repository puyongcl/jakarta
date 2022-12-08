package kq

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"jakarta/app/im/mq/internal/svc"
	pbIm "jakarta/app/im/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/xerr"
)

// 发送微信服务号订阅消息
type SendWxFwhNotifyMsgMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendWxFwhNotifyMsgMq(ctx context.Context, svcCtx *svc.ServiceContext) *SendWxFwhNotifyMsgMq {
	return &SendWxFwhNotifyMsgMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendWxFwhNotifyMsgMq) Consume(_, val string) error {
	var message kqueue.SendFwhNotifyMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendWxFwhNotifyMsgMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendWxFwhNotifyMsgMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *SendWxFwhNotifyMsgMq) execService(message *kqueue.SendFwhNotifyMessage) error {
	if l.svcCtx.Config.Mode == service.DevMode {
		logx.WithContext(l.ctx).Infof("SendWxFwhNotifyMsgMq 开发环境，暂不支持")
		return nil
	}
	var in pbIm.SendFwhTemplateMsgReq
	_ = copier.Copy(&in, message)
	//
	rsp, err := l.svcCtx.UserRpc.GetUserWxOpenId(l.ctx, &pbUser.GetUserWxOpenIdReq{Uid: message.ToUid})
	if err != nil {
		return err
	}
	//
	if rsp.FwhOpenId == "" {
		logx.WithContext(l.ctx).Infof("SendWxFwhNotifyMsgMq open id is empty uid :%d", message.ToUid)
		return nil //xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "fwh open ip empty")
	}
	in.OpenId = rsp.FwhOpenId

	// 模版id
	switch message.MsgType {
	case notify.DefineNotifyMsgTypeFwhMsg1:
		in.TemplateId = notify.FwhTemplateMsg1
		in.Path = notify.DefineNotifyMsgTypeFwhMsg1Path
		in.Color = notify.DefineNotifyMsgTypeFwhMsg1Color

	case notify.DefineNotifyMsgTypeFwhMsg2:
		in.TemplateId = notify.FwhTemplateMsg2
		in.Path = notify.DefineNotifyMsgTypeFwhMsg2Path
		in.Color = notify.DefineNotifyMsgTypeFwhMsg2Color

	case notify.DefineNotifyMsgTypeFwhMsg3:
		in.TemplateId = notify.FwhTemplateMsg3
		in.Path = notify.DefineNotifyMsgTypeFwhMsg3Path
		in.Color = notify.DefineNotifyMsgTypeFwhMsg3Color

	case notify.DefineNotifyMsgTypeFwhMsg4:
		in.TemplateId = notify.FwhTemplateMsg4
		in.Color = notify.DefineNotifyMsgTypeFwhMsg4Color

	case notify.DefineNotifyMsgTypeFwhMsg5:
		in.TemplateId = notify.FwhTemplateMsg5
		in.Path = notify.DefineNotifyMsgTypeFwhMsg5Path
		in.Color = notify.DefineNotifyMsgTypeFwhMsg5Color

	default:
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "wrong msg type")
	}

	//
	_, err = l.svcCtx.ImRpc.SendFwhTemplateMsg(l.ctx, &in)
	if err != nil {
		return err
	}
	return nil
}
