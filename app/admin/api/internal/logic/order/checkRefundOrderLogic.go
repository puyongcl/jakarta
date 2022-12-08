package order

import (
	"context"
	"jakarta/app/admin/api/internal/logic/adminlog"
	pbOrder "jakarta/app/order/rpc/pb"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckRefundOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckRefundOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckRefundOrderLogic {
	return &CheckRefundOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckRefundOrderLogic) CheckRefundOrder(req *types.CheckRefundOrderReq) (resp *types.CheckRefundOrderResp, err error) {
	defer func() {
		adminlog.SaveAdminLog(l.ctx, l.svcCtx.AdminLogModel, "CheckRefundOrder", req.AdminUid, err, req, resp)
	}()
	var in pbOrder.DoChatOrderActionReq
	in.OrderId = req.OrderId
	in.Action = req.Action
	in.Remark = req.Remark
	in.OperatorUid = req.AdminUid
	_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &in)
	if err != nil {
		return nil, err
	}
	return
}
