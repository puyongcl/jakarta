package user

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"jakarta/common/kqueue"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeNotifyMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeNotifyMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeNotifyMsgLogic {
	return &SubscribeNotifyMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeNotifyMsgLogic) SubscribeNotifyMsg(req *types.SubscribeNotifyMsgReq) (resp *types.SubscribeNotifyMsgResp, err error) {
	if req.Uid == 0 || req.MsgType == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	var kqMsg kqueue.SubscribeNotifyMsgMessage
	_ = copier.Copy(&kqMsg, req)
	buf, err := json.Marshal(&kqMsg)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.KqSubscribeNotifyMsgClient.Push(string(buf))
	if err != nil {
		return nil, err
	}
	return
}
