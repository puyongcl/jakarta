package user

import (
	"context"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubscribeMultiNotifyMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubscribeMultiNotifyMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubscribeMultiNotifyMsgLogic {
	return &SubscribeMultiNotifyMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubscribeMultiNotifyMsgLogic) SubscribeMultiNotifyMsg(req *types.SubscribeMultiNotifyMsgReq) (resp *types.SubscribeMultiNotifyMsgResp, err error) {
	if len(req.Sub) <= 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}
	sng := NewSubscribeNotifyMsgLogic(l.ctx, l.svcCtx)
	for k, _ := range req.Sub {
		_, err = sng.SubscribeNotifyMsg(req.Sub[k])
		if err != nil {
			return nil, err
		}
	}
	resp = &types.SubscribeMultiNotifyMsgResp{}
	return
}
