package listener

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	pbChat "jakarta/app/chat/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/money"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailListenerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailListenerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailListenerLogic {
	return &DetailListenerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailListenerLogic) DetailListener(req *types.GetListenerInfoReq) (resp *types.GetListenerInfoResp, err error) {
	if req.Uid == 0 || req.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}

	var pp *pbOrder.GetBusinessChatPricePlanResp
	var rp *pbListener.GetListenerProfileByUserResp
	var sp *pbChat.SyncChatStateResp
	err = mr.Finish(func() error {
		var err2 error
		var in pbListener.GetListenerProfileByUserReq
		_ = copier.Copy(&in, req)
		rp, err2 = l.svcCtx.ListenerRpc.GetListenerProfileByUser(l.ctx, &in)
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
		},
		func() error {
			var in pbChat.SyncChatStateReq
			in.Uid = req.Uid
			in.ListenerUid = req.ListenerUid
			in.Action = chatkey.ChatAction12
			var err2 error
			sp, err2 = l.svcCtx.ChatRpc.SyncChatState(l.ctx, &in)
			if err2 != nil {
				return err2
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	var val types.UserSeeListenerProfile
	_ = copier.Copy(&val, rp.Profile)
	//
	var arg money.CurrentListenerChatPriceArg
	_ = copier.Copy(&arg, rp.Profile)
	_ = copier.Copy(&arg, pp.Config)
	res := money.GetCurrentListenerChatPrice(&arg)
	// 展示价格
	val.ShowPrice = val.TextChatPrice
	val.ShowActualPrice = res.TextChatActualPrice
	val.FreeFlag = res.FreeFlag
	val.NewUserDiscount = res.NewUserDiscount
	val.FreeMinute = res.FreeMinute

	if val.TextChatSwitch != db.Enable {
		val.ShowPrice = val.VoiceChatPrice
		val.ShowActualPrice = res.VoiceChatActualPrice
	}
	//
	var val2 types.ListenerChatState
	_ = copier.Copy(&val2, sp)

	resp = &types.GetListenerInfoResp{Info: &val, State: &val2}
	return
}
