package kq

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/svc"
	pbIm "jakarta/app/im/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/xerr"
)

// 发送小程序订阅消息
type SendMiniProgramSubscribeMsgMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMiniProgramSubscribeMsgMq(ctx context.Context, svcCtx *svc.ServiceContext) *SendMiniProgramSubscribeMsgMq {
	return &SendMiniProgramSubscribeMsgMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMiniProgramSubscribeMsgMq) Consume(_, val string) error {
	var message kqueue.SendMiniProgramSubscribeMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendMiniProgramSubscribeMsgMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("SendMiniProgramSubscribeMsgMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *SendMiniProgramSubscribeMsgMq) execService(message *kqueue.SendMiniProgramSubscribeMessage) error {
	var in pbIm.SendMiniProgramSubscribeMsgReq
	_ = copier.Copy(&in, message)
	//
	rsp, err := l.svcCtx.UserRpc.GetUserWxOpenId(l.ctx, &pbUser.GetUserWxOpenIdReq{Uid: message.ToUid})
	if err != nil {
		return err
	}
	if rsp.MpOpenId == "" {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "mp open ip empty")
	}
	//
	in.OpenId = rsp.MpOpenId

	// 模版id
	switch message.MsgType {
	case notify.DefineNotifyMsgTypeMiniProgramMsg1:
		in.TemplateId = notify.MiniProgramNotifyTemplateId1

	case notify.DefineNotifyMsgTypeMiniProgramMsg2:
		in.TemplateId = notify.MiniProgramNotifyTemplateId2

	case notify.DefineNotifyMsgTypeMiniProgramMsg3:
		in.TemplateId = notify.MiniProgramNotifyTemplateId3

	default:
		return xerr.NewErrCodeMsg(xerr.RequestParamError, "wrong msg type")
	}

	//
	_, err = l.svcCtx.ImRpc.SendMiniProgramSubscribeMsg(l.ctx, &in)
	if err != nil {
		return err
	}
	return nil
}
