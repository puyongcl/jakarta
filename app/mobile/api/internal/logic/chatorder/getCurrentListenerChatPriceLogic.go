package chatorder

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	pbBusiness "jakarta/app/listener/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/money"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentListenerChatPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentListenerChatPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentListenerChatPriceLogic {
	return &GetCurrentListenerChatPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentListenerChatPriceLogic) GetCurrentListenerChatPrice(req *types.GetCurrentListenerChatPriceReq) (resp *types.GetCurrentListenerChatPriceResp, err error) {
	if req.ListenerUid == 0 || req.Uid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	authKey := ctxdata.GetUserAuthKeyFromCtx(l.ctx)
	var pp *pbOrder.GetBusinessChatPricePlanResp
	var lp *pbBusiness.GetListenerPriceResp
	err = mr.Finish(func() error {
		var in pbOrder.GetBusinessChatPricePlanReq
		_ = copier.Copy(&in, req)
		var err2 error
		pp, err2 = l.svcCtx.OrderRpc.GetBusinessChatPricePlan(l.ctx, &in)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var in pbBusiness.GetListenerPriceReq
		_ = copier.Copy(&in, req)
		in.AuthKey = authKey
		var err2 error
		lp, err2 = l.svcCtx.ListenerRpc.GetListenerPrice(l.ctx, &in)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	//
	var arg money.CurrentListenerChatPriceArg
	_ = copier.Copy(&arg, lp)
	_ = copier.Copy(&arg, pp.Config)
	res := money.GetCurrentListenerChatPrice(&arg)

	resp = &types.GetCurrentListenerChatPriceResp{
		PricingPlanId: pp.Config.Id,
	}
	_ = copier.Copy(resp, res)
	if resp.FreeMinute > 0 {
		resp.VoiceChatActualPrice = 0
		resp.TextChatActualPrice = 0
	}
	return
}
