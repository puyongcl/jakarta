package chatorder

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/core/stores/redis"
	pbBusiness "jakarta/app/listener/rpc/pb"
	"jakarta/app/order/rpc/order"
	pbOrder "jakarta/app/order/rpc/pb"
	"jakarta/app/payment/rpc/payment"
	pbPayment "jakarta/app/payment/rpc/pb"
	pbUser "jakarta/app/usercenter/rpc/pb"
	"jakarta/common/ctxdata"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/paykey"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateChatOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateChatOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateChatOrderLogic {
	return &CreateChatOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateChatOrderLogic) CreateChatOrder(req *types.CreateChatOrderReq) (resp *types.CreateChatOrderResp, err error) {
	if req.Uid == req.ListenerUid || req.Uid == 0 || req.ListenerUid == 0 || (req.OrderType != orderkey.ListenerOrderTypeTextChat && req.OrderType != orderkey.ListenerOrderTypeVoiceChat) {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}

	// 加分布式锁
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid != req.Uid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("uid not match %d-%d", uid, req.Uid))
	}
	rkey := fmt.Sprintf(rediskey.RedisLockUser, uid)
	rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
	rl.SetExpire(2)
	b, err := rl.AcquireCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
	}
	defer func() {
		_, err2 := rl.ReleaseCtx(l.ctx)
		if err2 != nil {
			logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
			return
		}
	}()

	authKey := ctxdata.GetUserAuthKeyFromCtx(l.ctx)

	// 获得XXX价格
	var rs1 *pbBusiness.GetListenerPriceResp
	var rs2 *pbUser.GetUserAuthyUserIdResp
	err = mr.Finish(func() error {
		var in1 pbBusiness.GetListenerPriceReq
		_ = copier.Copy(&in1, req)
		in1.AuthKey = authKey
		var err2 error
		rs1, err2 = l.svcCtx.ListenerRpc.GetListenerPrice(l.ctx, &in1)
		if err2 != nil {
			return err2
		}
		return nil
	}, func() error {
		var in2 pbUser.GetUserAuthByUserIdReq
		in2.Uid = req.Uid
		var err2 error
		rs2, err2 = l.svcCtx.UsercenterRpc.GetUserAuthByUserId(l.ctx, &in2)
		if err2 != nil {
			return err2
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// 创建订单
	var in3 pbOrder.CreateChatOrderReq
	var rs3 *pbOrder.CreateChatOrderResp
	_ = copier.Copy(&in3, req)
	_ = copier.Copy(&in3, &rs1)
	in3.UserChannel = req.Channel
	in3.ListenerChannel = rs2.UserAuth.Channel

	rs3, err = l.svcCtx.OrderRpc.CreateChatOrder(l.ctx, &in3)
	if err != nil {
		return nil, err
	}
	resp = &types.CreateChatOrderResp{}
	_ = copier.Copy(resp, rs3)

	// 调用预支付
	if rs3.ActualAmount > 0 {
		desc := listenerkey.GetChatOrderDescription(rs3.OrderType, rs3.BuyUnit*resp.ChatUnitMinute)
		var rs4 *pbPayment.CreatePaymentResp
		rs4, err = l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
			Uid:         req.Uid,
			PayModel:    paykey.ThirdPaymentPayModelWechatPay,
			PayAmount:   rs3.ActualAmount,
			OrderId:     rs3.OrderId,
			OrderType:   rs3.OrderType,
			Description: desc,
			OpenId:      rs2.UserAuth.AuthKey,
		})
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(resp, rs4)
	} else { // 免费订单
		// 更新订单支付状态
		_, err = l.svcCtx.OrderRpc.DoChatOrderAction(l.ctx, &order.DoChatOrderActionReq{
			OrderId:     rs3.OrderId,
			Action:      orderkey.GetPaySuccessOrderState(rs3.OrderType),
			OperatorUid: orderkey.DefaultSystemOperatorUid,
		})
		if err != nil {
			return nil, err
		}
	}

	return
}
