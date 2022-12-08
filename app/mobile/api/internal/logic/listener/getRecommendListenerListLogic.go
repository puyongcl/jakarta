package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/db"
	"jakarta/common/money"
	"jakarta/common/xerr"
)

type GetRecommendListenerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecommendListenerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendListenerListLogic {
	return &GetRecommendListenerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecommendListenerListLogic) GetRecommendListenerList(req *types.RecommendListenerReq) (resp *types.RecommendListenerResp, err error) {
	authKey := ctxdata.GetUserAuthKeyFromCtx(l.ctx)
	return l.DoGetRecommendListenerList(req, authKey)
}

func (l *GetRecommendListenerListLogic) DoGetRecommendListenerList(req *types.RecommendListenerReq, authKey string) (resp *types.RecommendListenerResp, err error) {
	if req.Uid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}

	var pp *pbOrder.GetBusinessChatPricePlanResp
	var rp *pbListener.GetRecommendListenerByUserResp
	err = mr.Finish(func() error {
		var in pbListener.GetRecommendListenerByUserReq
		_ = copier.Copy(&in, req)
		in.AuthKey = authKey
		var err2 error
		rp, err2 = l.svcCtx.ListenerRpc.GetRecommendListenerListByUser(l.ctx, &in)
		if err2 != nil {
			return err2
		}
		return nil
	},
		func() error {
			var in pbOrder.GetBusinessChatPricePlanReq
			_ = copier.Copy(&in, req)
			var err2 error
			pp, err2 = l.svcCtx.OrderRpc.GetBusinessChatPricePlan(l.ctx, &in)
			if err2 != nil {
				return err2
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	if len(rp.Listener) <= 0 {
		return &types.RecommendListenerResp{}, nil
	}
	resp = &types.RecommendListenerResp{
		Listener: make([]*types.UserSeeRecommendListenerProfile, 0),
	}
	for idx := 0; idx < len(rp.Listener); idx++ {
		val := AddListenerData(rp.Listener[idx], pp.Config)
		resp.Listener = append(resp.Listener, val)
	}
	return
}

func AddListenerData(lp *pbListener.UserSeeRecommendListenerProfile, bp *pbOrder.BusinessChatPricePlan) *types.UserSeeRecommendListenerProfile {
	var val types.UserSeeRecommendListenerProfile
	_ = copier.Copy(&val, lp)
	//
	var arg money.CurrentListenerChatPriceArg
	_ = copier.Copy(&arg, lp)
	_ = copier.Copy(&arg, bp)
	res := money.GetCurrentListenerChatPrice(&arg)
	_ = copier.Copy(&val, res)
	// 新用户展示XX
	if bp.NewUserDiscount != money.DivNumber && bp.NewUserDiscount != 0 {
		val.FreeFlag = "XX"
	}

	// 展示价格
	val.ShowPrice = val.TextChatPrice
	val.ShowActualPrice = val.TextChatActualPrice
	if val.TextChatSwitch != db.Enable {
		val.ShowPrice = val.VoiceChatPrice
		val.ShowActualPrice = val.VoiceChatActualPrice
	}

	return &val
}
