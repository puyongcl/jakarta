// Code generated by goctl. DO NOT EDIT!
// Source: order.proto

package order

import (
	"context"

	"jakarta/app/order/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AutoProcessOrder                   = pb.AutoProcessOrder
	BusinessChatPricePlan              = pb.BusinessChatPricePlan
	ChatOrder                          = pb.ChatOrder
	ChatOrderStateLog                  = pb.ChatOrderStateLog
	CreateChatOrderReq                 = pb.CreateChatOrderReq
	CreateChatOrderResp                = pb.CreateChatOrderResp
	DoChatOrderActionReq               = pb.DoChatOrderActionReq
	DoChatOrderActionResp              = pb.DoChatOrderActionResp
	EstimateOrderReq                   = pb.EstimateOrderReq
	EstimateOrderResp                  = pb.EstimateOrderResp
	ExpireVoiceChatOrder               = pb.ExpireVoiceChatOrder
	FeedbackOrderReq                   = pb.FeedbackOrderReq
	FeedbackOrderResp                  = pb.FeedbackOrderResp
	GetAutoProcessOrderReq             = pb.GetAutoProcessOrderReq
	GetAutoProcessOrderResp            = pb.GetAutoProcessOrderResp
	GetBusinessChatPricePlanReq        = pb.GetBusinessChatPricePlanReq
	GetBusinessChatPricePlanResp       = pb.GetBusinessChatPricePlanResp
	GetChatOrderDetailReq              = pb.GetChatOrderDetailReq
	GetChatOrderDetailResp             = pb.GetChatOrderDetailResp
	GetChatOrderFeedbackListByUserReq  = pb.GetChatOrderFeedbackListByUserReq
	GetChatOrderFeedbackListByUserResp = pb.GetChatOrderFeedbackListByUserResp
	GetChatOrderListReq                = pb.GetChatOrderListReq
	GetChatOrderListResp               = pb.GetChatOrderListResp
	GetChatOrderStateLogReq            = pb.GetChatOrderStateLogReq
	GetChatOrderStateLogResp           = pb.GetChatOrderStateLogResp
	GetExpireVoiceChatOrderReq         = pb.GetExpireVoiceChatOrderReq
	GetExpireVoiceChatOrderResp        = pb.GetExpireVoiceChatOrderResp
	GetLastCommentOrderReq             = pb.GetLastCommentOrderReq
	GetLastCommentOrderResp            = pb.GetLastCommentOrderResp
	GetListenerCommentListReq          = pb.GetListenerCommentListReq
	GetListenerCommentListResp         = pb.GetListenerCommentListResp
	GetListenerRecentGoodCommentReq    = pb.GetListenerRecentGoodCommentReq
	GetListenerRecentGoodCommentResp   = pb.GetListenerRecentGoodCommentResp
	GetOrderListByAdminReq             = pb.GetOrderListByAdminReq
	GetOrderListByAdminResp            = pb.GetOrderListByAdminResp
	GetRecentGoodCommentReq            = pb.GetRecentGoodCommentReq
	GetRecentGoodCommentResp           = pb.GetRecentGoodCommentResp
	GetRecentPaidUserCntReq            = pb.GetRecentPaidUserCntReq
	GetRecentPaidUserCntResp           = pb.GetRecentPaidUserCntResp
	ListenerOrderOpinion               = pb.ListenerOrderOpinion
	RecentGoodComment                  = pb.RecentGoodComment
	RefundChatOrder                    = pb.RefundChatOrder
	ReplyOrderCommentReq               = pb.ReplyOrderCommentReq
	ReplyOrderCommentResp              = pb.ReplyOrderCommentResp
	SettleChatOrderReq                 = pb.SettleChatOrderReq
	SettleChatOrderResp                = pb.SettleChatOrderResp
	ShortChatOrder                     = pb.ShortChatOrder
	ShortListenerProfile               = pb.ShortListenerProfile
	UpdateChatOrderUseReq              = pb.UpdateChatOrderUseReq
	UpdateChatOrderUseResp             = pb.UpdateChatOrderUseResp
	UpdateOrderLastDaysStatReq         = pb.UpdateOrderLastDaysStatReq
	UpdateOrderLastDaysStatResp        = pb.UpdateOrderLastDaysStatResp
	UserSeeChatOrderFeedback           = pb.UserSeeChatOrderFeedback

	Order interface {
		// 获取订单费用估算
		EstimateOrder(ctx context.Context, in *EstimateOrderReq, opts ...grpc.CallOption) (*EstimateOrderResp, error)
		// 咨询服务下订单
		CreateChatOrder(ctx context.Context, in *CreateChatOrderReq, opts ...grpc.CallOption) (*CreateChatOrderResp, error)
		// 获取订单详情
		GetChatOrderDetail(ctx context.Context, in *GetChatOrderDetailReq, opts ...grpc.CallOption) (*GetChatOrderDetailResp, error)
		// 订单操作
		DoChatOrderAction(ctx context.Context, in *DoChatOrderActionReq, opts ...grpc.CallOption) (*DoChatOrderActionResp, error)
		// 订单列表
		GetChatOrderList(ctx context.Context, in *GetChatOrderListReq, opts ...grpc.CallOption) (*GetChatOrderListResp, error)
		// 获取聊天服务价格配置
		GetBusinessChatPricePlan(ctx context.Context, in *GetBusinessChatPricePlanReq, opts ...grpc.CallOption) (*GetBusinessChatPricePlanResp, error)
		// 获取用户订单状态变化记录
		GetChatOrderStateLog(ctx context.Context, in *GetChatOrderStateLogReq, opts ...grpc.CallOption) (*GetChatOrderStateLogResp, error)
		// 更新订单的使用情况
		UpdateChatOrderUse(ctx context.Context, in *UpdateChatOrderUseReq, opts ...grpc.CallOption) (*UpdateChatOrderUseResp, error)
		// 获取过期的语音订单
		GetExpireVoiceChatOrder(ctx context.Context, in *GetExpireVoiceChatOrderReq, opts ...grpc.CallOption) (*GetExpireVoiceChatOrderResp, error)
		// 结算订单
		SettleChatOrder(ctx context.Context, in *SettleChatOrderReq, opts ...grpc.CallOption) (*SettleChatOrderResp, error)
		// 获取XXX评价列表
		GetListenerCommentList(ctx context.Context, in *GetListenerCommentListReq, opts ...grpc.CallOption) (*GetListenerCommentListResp, error)
		// XXX回复用户的订单评价
		ReplyOrderComment(ctx context.Context, in *ReplyOrderCommentReq, opts ...grpc.CallOption) (*ReplyOrderCommentResp, error)
		// XXX反馈
		FeedbackOrder(ctx context.Context, in *FeedbackOrderReq, opts ...grpc.CallOption) (*FeedbackOrderResp, error)
		// 管理后台获取订单列表
		GetOrderListByAdmin(ctx context.Context, in *GetOrderListByAdminReq, opts ...grpc.CallOption) (*GetOrderListByAdminResp, error)
		// 获取需要自动处理的订单
		GetAutoProcessOrder(ctx context.Context, in *GetAutoProcessOrderReq, opts ...grpc.CallOption) (*GetAutoProcessOrderResp, error)
		// 更新XXX订单统计数据
		UpdateOrderLastDaysStat(ctx context.Context, in *UpdateOrderLastDaysStatReq, opts ...grpc.CallOption) (*UpdateOrderLastDaysStatResp, error)
		// 获取最近的好评
		GetRecentGoodComment(ctx context.Context, in *GetRecentGoodCommentReq, opts ...grpc.CallOption) (*GetRecentGoodCommentResp, error)
		// 获取指定XXX的好评
		GetListenerRecentGoodComment(ctx context.Context, in *GetListenerRecentGoodCommentReq, opts ...grpc.CallOption) (*GetListenerRecentGoodCommentResp, error)
		// 获取最近一条评价
		GetLastCommentOrder(ctx context.Context, in *GetLastCommentOrderReq, opts ...grpc.CallOption) (*GetLastCommentOrderResp, error)
		// 获取最近时间段付费用户数
		GetRecentPaidUserCnt(ctx context.Context, in *GetRecentPaidUserCntReq, opts ...grpc.CallOption) (*GetRecentPaidUserCntResp, error)
		// 用户获取反馈列表
		GetChatOrderFeedbackListByUser(ctx context.Context, in *GetChatOrderFeedbackListByUserReq, opts ...grpc.CallOption) (*GetChatOrderFeedbackListByUserResp, error)
	}

	defaultOrder struct {
		cli zrpc.Client
	}
)

func NewOrder(cli zrpc.Client) Order {
	return &defaultOrder{
		cli: cli,
	}
}

// 获取订单费用估算
func (m *defaultOrder) EstimateOrder(ctx context.Context, in *EstimateOrderReq, opts ...grpc.CallOption) (*EstimateOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.EstimateOrder(ctx, in, opts...)
}

// 咨询服务下订单
func (m *defaultOrder) CreateChatOrder(ctx context.Context, in *CreateChatOrderReq, opts ...grpc.CallOption) (*CreateChatOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.CreateChatOrder(ctx, in, opts...)
}

// 获取订单详情
func (m *defaultOrder) GetChatOrderDetail(ctx context.Context, in *GetChatOrderDetailReq, opts ...grpc.CallOption) (*GetChatOrderDetailResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetChatOrderDetail(ctx, in, opts...)
}

// 订单操作
func (m *defaultOrder) DoChatOrderAction(ctx context.Context, in *DoChatOrderActionReq, opts ...grpc.CallOption) (*DoChatOrderActionResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.DoChatOrderAction(ctx, in, opts...)
}

// 订单列表
func (m *defaultOrder) GetChatOrderList(ctx context.Context, in *GetChatOrderListReq, opts ...grpc.CallOption) (*GetChatOrderListResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetChatOrderList(ctx, in, opts...)
}

// 获取聊天服务价格配置
func (m *defaultOrder) GetBusinessChatPricePlan(ctx context.Context, in *GetBusinessChatPricePlanReq, opts ...grpc.CallOption) (*GetBusinessChatPricePlanResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetBusinessChatPricePlan(ctx, in, opts...)
}

// 获取用户订单状态变化记录
func (m *defaultOrder) GetChatOrderStateLog(ctx context.Context, in *GetChatOrderStateLogReq, opts ...grpc.CallOption) (*GetChatOrderStateLogResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetChatOrderStateLog(ctx, in, opts...)
}

// 更新订单的使用情况
func (m *defaultOrder) UpdateChatOrderUse(ctx context.Context, in *UpdateChatOrderUseReq, opts ...grpc.CallOption) (*UpdateChatOrderUseResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.UpdateChatOrderUse(ctx, in, opts...)
}

// 获取过期的语音订单
func (m *defaultOrder) GetExpireVoiceChatOrder(ctx context.Context, in *GetExpireVoiceChatOrderReq, opts ...grpc.CallOption) (*GetExpireVoiceChatOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetExpireVoiceChatOrder(ctx, in, opts...)
}

// 结算订单
func (m *defaultOrder) SettleChatOrder(ctx context.Context, in *SettleChatOrderReq, opts ...grpc.CallOption) (*SettleChatOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.SettleChatOrder(ctx, in, opts...)
}

// 获取XXX评价列表
func (m *defaultOrder) GetListenerCommentList(ctx context.Context, in *GetListenerCommentListReq, opts ...grpc.CallOption) (*GetListenerCommentListResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetListenerCommentList(ctx, in, opts...)
}

// XXX回复用户的订单评价
func (m *defaultOrder) ReplyOrderComment(ctx context.Context, in *ReplyOrderCommentReq, opts ...grpc.CallOption) (*ReplyOrderCommentResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.ReplyOrderComment(ctx, in, opts...)
}

// XXX反馈
func (m *defaultOrder) FeedbackOrder(ctx context.Context, in *FeedbackOrderReq, opts ...grpc.CallOption) (*FeedbackOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.FeedbackOrder(ctx, in, opts...)
}

// 管理后台获取订单列表
func (m *defaultOrder) GetOrderListByAdmin(ctx context.Context, in *GetOrderListByAdminReq, opts ...grpc.CallOption) (*GetOrderListByAdminResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetOrderListByAdmin(ctx, in, opts...)
}

// 获取需要自动处理的订单
func (m *defaultOrder) GetAutoProcessOrder(ctx context.Context, in *GetAutoProcessOrderReq, opts ...grpc.CallOption) (*GetAutoProcessOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetAutoProcessOrder(ctx, in, opts...)
}

// 更新XXX订单统计数据
func (m *defaultOrder) UpdateOrderLastDaysStat(ctx context.Context, in *UpdateOrderLastDaysStatReq, opts ...grpc.CallOption) (*UpdateOrderLastDaysStatResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.UpdateOrderLastDaysStat(ctx, in, opts...)
}

// 获取最近的好评
func (m *defaultOrder) GetRecentGoodComment(ctx context.Context, in *GetRecentGoodCommentReq, opts ...grpc.CallOption) (*GetRecentGoodCommentResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetRecentGoodComment(ctx, in, opts...)
}

// 获取指定XXX的好评
func (m *defaultOrder) GetListenerRecentGoodComment(ctx context.Context, in *GetListenerRecentGoodCommentReq, opts ...grpc.CallOption) (*GetListenerRecentGoodCommentResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetListenerRecentGoodComment(ctx, in, opts...)
}

// 获取最近一条评价
func (m *defaultOrder) GetLastCommentOrder(ctx context.Context, in *GetLastCommentOrderReq, opts ...grpc.CallOption) (*GetLastCommentOrderResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetLastCommentOrder(ctx, in, opts...)
}

// 获取最近时间段付费用户数
func (m *defaultOrder) GetRecentPaidUserCnt(ctx context.Context, in *GetRecentPaidUserCntReq, opts ...grpc.CallOption) (*GetRecentPaidUserCntResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetRecentPaidUserCnt(ctx, in, opts...)
}

// 用户获取反馈列表
func (m *defaultOrder) GetChatOrderFeedbackListByUser(ctx context.Context, in *GetChatOrderFeedbackListByUserReq, opts ...grpc.CallOption) (*GetChatOrderFeedbackListByUserResp, error) {
	client := pb.NewOrderClient(m.cli.Conn())
	return client.GetChatOrderFeedbackListByUser(ctx, in, opts...)
}