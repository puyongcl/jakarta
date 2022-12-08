package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/orderkey"
	"jakarta/common/money"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBusinessChatPricePlanLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBusinessChatPricePlanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBusinessChatPricePlanLogic {
	return &GetBusinessChatPricePlanLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取聊天服务价格配置
func (l *GetBusinessChatPricePlanLogic) GetBusinessChatPricePlan(in *pb.GetBusinessChatPricePlanReq) (*pb.GetBusinessChatPricePlanResp, error) {
	return l.getListenerPrice(in.Uid, in.Id)
}

func (l *GetBusinessChatPricePlanLogic) getListenerPrice(uid int64, id int64) (*pb.GetBusinessChatPricePlanResp, error) {
	// 新用户
	cnt, err := l.svcCtx.ChatOrderModel.FindCount2(l.ctx, uid, 0, 0, false, orderkey.AbnormalOrderState, 0)
	if err != nil {
		return nil, err
	}
	//
	var rsp *orderPgModel.ChatOrderPricingPlan
	if id > 0 {
		rsp, err = l.svcCtx.ChatOrderPricingPlanModel.FindOne(l.ctx, id)
		if err != nil {
			return nil, err
		}
	} else {
		rsp, err = l.svcCtx.ChatOrderPricingPlanModel.FindPriceConfig(l.ctx)
		if err != nil {
			return nil, err
		}
	}

	var val pb.BusinessChatPricePlan
	_ = copier.Copy(&val, rsp)
	val.CreateTime = rsp.CreateTime.Format(db.DateTimeFormat)
	val.OrderCnt = cnt
	// 新用户折扣
	if cnt > 0 { // 不是新用户
		val.NewUserDiscount = money.DivNumber
		val.FreeMinute = 0
	}
	return &pb.GetBusinessChatPricePlanResp{Config: &val}, nil
}
