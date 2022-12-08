package kq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/mq/internal/svc"
	pbIm "jakarta/app/im/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/fwh"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"
)

// 发送自定义消息
type WxFwhCallbackEventMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWxFwhCallbackEventMq(ctx context.Context, svcCtx *svc.ServiceContext) *WxFwhCallbackEventMq {
	return &WxFwhCallbackEventMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WxFwhCallbackEventMq) Consume(_, val string) error {
	var message kqueue.WxFwhCallbackEventMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Errorf("WxFwhCallbackEventMq->Consume Unmarshal err : %+v , val : %s", err, val)
		return err
	}

	if err := l.execService(&message); err != nil {
		logx.WithContext(l.ctx).Errorf("WxFwhCallbackEventMq->execService err : %+v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

func (l *WxFwhCallbackEventMq) execService(message *kqueue.WxFwhCallbackEventMessage) error {
	var state int64
	var unionId string
	switch message.Event {
	case fwh.EventSubscribe:
		state = db.Enable
		rsp, err := l.svcCtx.ImRpc.GetUserUnionIdByFwhOpenId(l.ctx, &pbIm.GetUserUnionIdByFwhOpenIdReq{OpenId: message.OpenId})
		if err != nil {
			return err
		}
		unionId = rsp.UnionId

	case fwh.EventUnsubscribe:
		state = db.Disable

	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数错误 %s", message.Event))
	}

	//
	_, err := l.svcCtx.UserRpc.UpdateUserWxFwhState(l.ctx, &pbUser.UpdateUserWxFwhStateReq{
		OpenId:  message.OpenId,
		UnionId: unionId,
		State:   state,
	})

	return err
}
