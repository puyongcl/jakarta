package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	pbListener "jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCashListReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCashListReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCashListReqLogic {
	return &GetCashListReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCashListReqLogic) GetCashListReq(req *types.GetCashListReq) (resp *types.GetCashListResp, err error) {
	var in pbListener.GetListenerMoveCashListByAdminReq
	_ = copier.Copy(&in, req)
	rs, err := l.svcCtx.ListenerRpc.GetListenerMoveCashListByAdmin(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	resp = &types.GetCashListResp{List: make([]*types.CashLogDetail, 0)}
	_ = copier.Copy(resp, rs)

	return
}
