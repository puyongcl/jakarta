package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"jakarta/app/mqueue/job/jobtype"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/common/key/orderkey"
	"jakarta/common/money"
	"jakarta/common/tool"
	"jakarta/common/uniqueid"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/order/rpc/internal/svc"
	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateChatOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateChatOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateChatOrderLogic {
	return &CreateChatOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 咨询服务下订单
func (l *CreateChatOrderLogic) CreateChatOrder(in *pb.CreateChatOrderReq) (*pb.CreateChatOrderResp, error) {
	// 查找是否存在未支付的订单
	var err error
	cnt, err := l.svcCtx.ChatOrderModel.FindCount(l.ctx, in.Uid, 0, 0, []int64{orderkey.ChatOrderStateWaitPay1}, orderkey.NoPayOrderLimitSecond)
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "请求频率过快，请等待10秒再试。")
	}
	// 定价方案
	var pp *pb.GetBusinessChatPricePlanResp
	gpp := NewGetBusinessChatPricePlanLogic(l.ctx, l.svcCtx)
	pp, err = gpp.getListenerPrice(in.Uid, in.PricingPlanId)
	if err != nil {
		return nil, err
	}
	// 计算
	var order *orderPgModel.ChatOrder
	now := time.Now()
	order, err = GetOrder(&now, GenOrderId(in.OrderType), in, pp.Config)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ChatOrderModel.Insert(l.ctx, order)
	if err != nil {
		return nil, err
	}

	//
	var rsp pb.CreateChatOrderResp
	_ = copier.Copy(&rsp, order)

	// 延迟关闭订单 暂不需要后台关闭
	//l.deferCloseOrder(order.OrderId)

	return &rsp, nil
}

func (l *CreateChatOrderLogic) deferCloseOrder(orderId string) {
	payload, err := json.Marshal(jobtype.DeferCloseChatOrderPayload{OrderId: orderId})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("CreateChatOrder payload json marshal err:%+v", err)
		return
	}
	_, err = l.svcCtx.AsynqClient.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferCloseChatOrder, payload), asynq.ProcessIn(orderkey.CloseOrderTimeMinutes*time.Minute))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("CreateChatOrder AsynqClient.Enqueue err:%+v", err)
		return
	}
}

func GenOrderId(orderType int64) string {
	switch orderType {
	case orderkey.ListenerOrderTypeTextChat:
		return uniqueid.GenSn(uniqueid.SnPrefixTextChatOrderId)
	case orderkey.ListenerOrderTypeVoiceChat:
		return uniqueid.GenSn(uniqueid.SnPrefixVoiceChatOrderId)
	default:
		return ""
	}
}

func GetOrder(createTime *time.Time, orderId string, in *pb.CreateChatOrderReq, pp *pb.BusinessChatPricePlan) (*orderPgModel.ChatOrder, error) {
	var order orderPgModel.ChatOrder
	// 订单号
	order.OrderId = orderId
	order.ListenerUid = in.ListenerUid
	order.ListenerNickName = in.ListenerNickName
	order.ListenerAvatar = in.ListenerAvatar
	order.Uid = in.Uid
	order.NickName = in.NickName
	order.Avatar = in.Avatar
	order.ChatUnitMinute = pp.ChatUnitMinute
	order.BuyUnit = in.BuyUnit
	order.OrderType = in.OrderType
	order.PricingPlanId = in.PricingPlanId
	order.OrderState = orderkey.ChatOrderStateWaitPay1
	// 费用计算
	rs := calculate(in.TextChatPrice, in.VoiceChatPrice, in.OrderType, in.BuyUnit, createTime, pp)
	order.ChatUnitMinute = rs.ChatUnitMinute
	order.BuyUnit = rs.BuyUnit
	if in.BuyUnit != rs.BuyUnit { // 防止恶意刷免费单
		in.BuyUnit = rs.BuyUnit
	}
	order.UnitPrice = rs.UnitPrice
	order.BaseAmount = rs.BaseAmount
	order.TaxAmount = rs.TaxAmount
	order.NightAddAmount = rs.NightAddAmount
	order.SaveAmount = rs.SaveAmount
	order.ActualAmount = rs.ActualAmount
	order.UserChannel = in.UserChannel
	order.ListenerChannel = in.ListenerChannel

	return &order, nil
}

func calculate(textChatPrice, voiceChatPrice, orderType int64, buyUnit int64, now *time.Time, pp *pb.BusinessChatPricePlan) (result *pb.EstimateOrderResp) {
	result = &pb.EstimateOrderResp{BuyUnit: buyUnit, ChatUnitMinute: pp.ChatUnitMinute}
	// 计算展示单价
	var chatPrice int64
	var arg money.CurrentListenerChatPriceArg
	_ = copier.Copy(&arg, pp)
	arg.TextChatPrice = textChatPrice
	arg.VoiceChatPrice = voiceChatPrice
	cp := money.GetCurrentListenerChatPrice(&arg)
	switch orderType {
	case orderkey.ListenerOrderTypeTextChat:
		chatPrice = textChatPrice
		result.UnitPrice = cp.TextChatActualPrice
	case orderkey.ListenerOrderTypeVoiceChat:
		chatPrice = voiceChatPrice
		result.UnitPrice = cp.VoiceChatActualPrice
	default:
		return nil
	}

	// 优惠减免金额
	if pp.NewUserDiscount == 0 { // 0 折 免费
		// 校验免费时间
		result.BuyUnit = 1
		result.SaveAmount = tool.DivideInt64(chatPrice, pp.ChatUnitMinute) * cp.FreeMinute
		result.BaseAmount = 0
		result.ChatUnitMinute = cp.FreeMinute
		result.UnitPrice = 0

	} else if pp.NewUserDiscount < money.DivNumber { // 打折
		saveChatPrice := money.RoundDivideInt64(chatPrice*(money.DivNumber-pp.NewUserDiscount), money.DivNumber)
		result.SaveAmount = saveChatPrice * buyUnit
		result.BaseAmount = (chatPrice - saveChatPrice) * buyUnit
	} else {
		// 减免费用
		result.SaveAmount = 0
		// 基本服务费
		result.BaseAmount = chatPrice * buyUnit
	}

	// 夜间服务费
	nowH := int64(now.Hour())
	if nowH >= pp.NightAddPriceHourStart && nowH <= pp.NightAddPriceHourEnd {
		result.SaveAmount += money.RoundDivideInt64(result.SaveAmount*pp.NightAddPriceRate, money.DivNumber)
		result.NightAddAmount = money.RoundDivideInt64(result.BaseAmount*pp.NightAddPriceRate, money.DivNumber)
	} else {
		result.NightAddAmount = 0
	}

	// 税费
	result.TaxAmount = money.RoundDivideInt64((result.BaseAmount+result.NightAddAmount)*pp.TaxRate, money.DivNumber)
	result.SaveAmount += money.RoundDivideInt64(result.SaveAmount*pp.TaxRate, money.DivNumber)
	// 实际支付金额
	actualAmount := result.BaseAmount + result.NightAddAmount + result.TaxAmount
	result.ActualAmount = money.RoundYuan(actualAmount)
	if (result.ActualAmount > actualAmount) && (actualAmount > 0) && (pp.TaxRate > 0) { // 重新计算税费
		result.TaxAmount = result.ActualAmount - tool.DivideInt64(result.ActualAmount*money.DivNumber, money.DivNumber+pp.TaxRate)
	}
	return
}
