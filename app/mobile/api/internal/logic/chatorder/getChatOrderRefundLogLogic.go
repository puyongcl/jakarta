package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/order/rpc/pb"
	"jakarta/common/key/orderkey"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatOrderRefundLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatOrderRefundLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatOrderRefundLogLogic {
	return &GetChatOrderRefundLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatOrderRefundLogLogic) GetChatOrderRefundLog(req *types.GetChatOrderRefundLogReq) (resp *types.GetChatOrderRefundLogResp, err error) {
	var in pb.GetChatOrderStateLogReq
	_ = copier.Copy(&in, req)
	in.State = orderkey.ChatOrderRefundState
	rsp, err := l.svcCtx.OrderRpc.GetChatOrderStateLog(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetChatOrderRefundLogResp{List: make([]*types.ChatOrderRefundLog, 0)}
	_ = copier.Copy(resp, rsp)
	return
}
