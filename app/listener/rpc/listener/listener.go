// Code generated by goctl. DO NOT EDIT!
// Source: listener.proto

package listener

import (
	"context"

	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AdminSeeListenerMoveCash              = pb.AdminSeeListenerMoveCash
	AdminSeeListenerProfile               = pb.AdminSeeListenerProfile
	AdminSeeListenerProfileDraft          = pb.AdminSeeListenerProfileDraft
	ChangeWorkStateReq                    = pb.ChangeWorkStateReq
	ChangeWorkStateResp                   = pb.ChangeWorkStateResp
	CheckListenerProfile                  = pb.CheckListenerProfile
	CheckListenerProfileReq               = pb.CheckListenerProfileReq
	CheckListenerProfileResp              = pb.CheckListenerProfileResp
	CommentTagPair                        = pb.CommentTagPair
	CommitCheckNewListenerProfileReq      = pb.CommitCheckNewListenerProfileReq
	CommitCheckNewListenerProfileResp     = pb.CommitCheckNewListenerProfileResp
	EditListenerProfileDraftReq           = pb.EditListenerProfileDraftReq
	EditListenerProfileDraftResp          = pb.EditListenerProfileDraftResp
	EditListenerWordsReq                  = pb.EditListenerWordsReq
	EditListenerWordsResp                 = pb.EditListenerWordsResp
	FindListenerListRangeByUpdateTimeReq  = pb.FindListenerListRangeByUpdateTimeReq
	FindListenerListRangeByUpdateTimeResp = pb.FindListenerListRangeByUpdateTimeResp
	GenListenerContractReq                = pb.GenListenerContractReq
	GenListenerContractResp               = pb.GenListenerContractResp
	GetBankCardReq                        = pb.GetBankCardReq
	GetBankCardResp                       = pb.GetBankCardResp
	GetCommitMoveCashReq                  = pb.GetCommitMoveCashReq
	GetCommitMoveCashResp                 = pb.GetCommitMoveCashResp
	GetListenerBasicInfoReq               = pb.GetListenerBasicInfoReq
	GetListenerBasicInfoResp              = pb.GetListenerBasicInfoResp
	GetListenerCashLogReq                 = pb.GetListenerCashLogReq
	GetListenerCashLogResp                = pb.GetListenerCashLogResp
	GetListenerHomePageDashboardReq       = pb.GetListenerHomePageDashboardReq
	GetListenerHomePageDashboardResp      = pb.GetListenerHomePageDashboardResp
	GetListenerIncomeLogReq               = pb.GetListenerIncomeLogReq
	GetListenerIncomeLogResp              = pb.GetListenerIncomeLogResp
	GetListenerMoveCashListByAdminReq     = pb.GetListenerMoveCashListByAdminReq
	GetListenerMoveCashListByAdminResp    = pb.GetListenerMoveCashListByAdminResp
	GetListenerPriceReq                   = pb.GetListenerPriceReq
	GetListenerPriceResp                  = pb.GetListenerPriceResp
	GetListenerProfileByOwnReq            = pb.GetListenerProfileByOwnReq
	GetListenerProfileByOwnResp           = pb.GetListenerProfileByOwnResp
	GetListenerProfileByUserReq           = pb.GetListenerProfileByUserReq
	GetListenerProfileByUserResp          = pb.GetListenerProfileByUserResp
	GetListenerProfileListReq             = pb.GetListenerProfileListReq
	GetListenerProfileListResp            = pb.GetListenerProfileListResp
	GetListenerRatingStatReq              = pb.GetListenerRatingStatReq
	GetListenerRatingStatResp             = pb.GetListenerRatingStatResp
	GetListenerRemarkUserReq              = pb.GetListenerRemarkUserReq
	GetListenerRemarkUserResp             = pb.GetListenerRemarkUserResp
	GetListenerWalletReq                  = pb.GetListenerWalletReq
	GetListenerWalletResp                 = pb.GetListenerWalletResp
	GetListenerWordsReq                   = pb.GetListenerWordsReq
	GetListenerWordsResp                  = pb.GetListenerWordsResp
	GetRecommendListenerByUserReq         = pb.GetRecommendListenerByUserReq
	GetRecommendListenerByUserResp        = pb.GetRecommendListenerByUserResp
	GetRecommendListenerReq               = pb.GetRecommendListenerReq
	GetRecommendListenerResp              = pb.GetRecommendListenerResp
	GetWorkStateReq                       = pb.GetWorkStateReq
	GetWorkStateResp                      = pb.GetWorkStateResp
	ListenerCashLog                       = pb.ListenerCashLog
	ListenerIncomeLog                     = pb.ListenerIncomeLog
	ListenerRemarkUserReq                 = pb.ListenerRemarkUserReq
	ListenerRemarkUserResp                = pb.ListenerRemarkUserResp
	ListenerSeeOwnProfile                 = pb.ListenerSeeOwnProfile
	ListenerShortProfile                  = pb.ListenerShortProfile
	RecListenerWhenUserLoginReq           = pb.RecListenerWhenUserLoginReq
	RecListenerWhenUserLoginResp          = pb.RecListenerWhenUserLoginResp
	SetBankCardReq                        = pb.SetBankCardReq
	SetBankCardResp                       = pb.SetBankCardResp
	SnapshotLastDaysListenerStatReq       = pb.SnapshotLastDaysListenerStatReq
	SnapshotLastDaysListenerStatResp      = pb.SnapshotLastDaysListenerStatResp
	UpdateListenerDashboardStatReq        = pb.UpdateListenerDashboardStatReq
	UpdateListenerDashboardStatResp       = pb.UpdateListenerDashboardStatResp
	UpdateListenerEveryDayAverageStatReq  = pb.UpdateListenerEveryDayAverageStatReq
	UpdateListenerEveryDayAverageStatResp = pb.UpdateListenerEveryDayAverageStatResp
	UpdateListenerOnlineStateReq          = pb.UpdateListenerOnlineStateReq
	UpdateListenerOnlineStateResp         = pb.UpdateListenerOnlineStateResp
	UpdateListenerOrderStatReq            = pb.UpdateListenerOrderStatReq
	UpdateListenerOrderStatResp           = pb.UpdateListenerOrderStatResp
	UpdateListenerSuggestionReq           = pb.UpdateListenerSuggestionReq
	UpdateListenerSuggestionResp          = pb.UpdateListenerSuggestionResp
	UpdateListenerUserStatReq             = pb.UpdateListenerUserStatReq
	UpdateListenerUserStatResp            = pb.UpdateListenerUserStatResp
	UpdateListenerWalletReq               = pb.UpdateListenerWalletReq
	UpdateListenerWalletResp              = pb.UpdateListenerWalletResp
	UpdateRecommendListenerPoolReq        = pb.UpdateRecommendListenerPoolReq
	UpdateRecommendListenerPoolResp       = pb.UpdateRecommendListenerPoolResp
	UpdateTodayListenerUserStatReq        = pb.UpdateTodayListenerUserStatReq
	UpdateTodayListenerUserStatResp       = pb.UpdateTodayListenerUserStatResp
	UserSeeListenerProfile                = pb.UserSeeListenerProfile
	UserSeeRecommendListenerProfile       = pb.UserSeeRecommendListenerProfile

	Listener interface {
		// XXX首次填写或更新自己的资料
		AddOrUpdateListenerProfileDraft(ctx context.Context, in *EditListenerProfileDraftReq, opts ...grpc.CallOption) (*EditListenerProfileDraftResp, error)
		// 新申请XXX提交审核
		CommitCheckNewListenerProfile(ctx context.Context, in *CommitCheckNewListenerProfileReq, opts ...grpc.CallOption) (*CommitCheckNewListenerProfileResp, error)
		// 用户查看XXX资料（TODO 后台内部禁止调用）
		GetListenerProfileByUser(ctx context.Context, in *GetListenerProfileByUserReq, opts ...grpc.CallOption) (*GetListenerProfileByUserResp, error)
		// XXX获取自己的资料
		GetListenerProfileByOwn(ctx context.Context, in *GetListenerProfileByOwnReq, opts ...grpc.CallOption) (*GetListenerProfileByOwnResp, error)
		// XXX备注用户
		ListenerRemarkUser(ctx context.Context, in *ListenerRemarkUserReq, opts ...grpc.CallOption) (*ListenerRemarkUserResp, error)
		// 获取XXX备注的用户
		GetListenerRemarkUser(ctx context.Context, in *GetListenerRemarkUserReq, opts ...grpc.CallOption) (*GetListenerRemarkUserResp, error)
		// 用户获取XXX推荐列表
		GetRecommendListenerListByUser(ctx context.Context, in *GetRecommendListenerByUserReq, opts ...grpc.CallOption) (*GetRecommendListenerByUserResp, error)
		// 获取推荐的XXX
		GetRecommendListenerList(ctx context.Context, in *GetRecommendListenerReq, opts ...grpc.CallOption) (*GetRecommendListenerResp, error)
		// 获取XXX工作状态设置
		GetWorkState(ctx context.Context, in *GetWorkStateReq, opts ...grpc.CallOption) (*GetWorkStateResp, error)
		// 修改XXX工作状态
		ChangeWorkState(ctx context.Context, in *ChangeWorkStateReq, opts ...grpc.CallOption) (*ChangeWorkStateResp, error)
		// 获取XXX定价和价格方案
		GetListenerPrice(ctx context.Context, in *GetListenerPriceReq, opts ...grpc.CallOption) (*GetListenerPriceResp, error)
		// XXX钱包金额更新
		UpdateListenerWallet(ctx context.Context, in *UpdateListenerWalletReq, opts ...grpc.CallOption) (*UpdateListenerWalletResp, error)
		// 更新XXX的订单统计数据
		UpdateListenerOrderStat(ctx context.Context, in *UpdateListenerOrderStatReq, opts ...grpc.CallOption) (*UpdateListenerOrderStatResp, error)
		// 绑定银行卡
		SetBankCard(ctx context.Context, in *SetBankCardReq, opts ...grpc.CallOption) (*SetBankCardResp, error)
		// 获取银行卡
		GetBankCard(ctx context.Context, in *GetBankCardReq, opts ...grpc.CallOption) (*GetBankCardResp, error)
		// 获取XXX钱包详情
		GetListenerWallet(ctx context.Context, in *GetListenerWalletReq, opts ...grpc.CallOption) (*GetListenerWalletResp, error)
		// 获取提现记录
		GetListenerCashLog(ctx context.Context, in *GetListenerCashLogReq, opts ...grpc.CallOption) (*GetListenerCashLogResp, error)
		// 获取收益记录
		GetListenerIncomeLog(ctx context.Context, in *GetListenerIncomeLogReq, opts ...grpc.CallOption) (*GetListenerIncomeLogResp, error)
		// 获取XXX评价统计情况
		GetListenerRatingStat(ctx context.Context, in *GetListenerRatingStatReq, opts ...grpc.CallOption) (*GetListenerRatingStatResp, error)
		// 获取XXX常用语
		GetListenerWords(ctx context.Context, in *GetListenerWordsReq, opts ...grpc.CallOption) (*GetListenerWordsResp, error)
		// 编辑XXX常用语
		EditListenerWords(ctx context.Context, in *EditListenerWordsReq, opts ...grpc.CallOption) (*EditListenerWordsResp, error)
		// 更新XXX与用户的交互情况
		UpdateListenerUserStat(ctx context.Context, in *UpdateListenerUserStatReq, opts ...grpc.CallOption) (*UpdateListenerUserStatResp, error)
		// 更新统计今日推荐和浏览XXX资料页统计
		UpdateTodayListenerUserStat(ctx context.Context, in *UpdateTodayListenerUserStatReq, opts ...grpc.CallOption) (*UpdateTodayListenerUserStatResp, error)
		// 更新XXX首页数据统计看板
		UpdateListenerDashboardStat(ctx context.Context, in *UpdateListenerDashboardStatReq, opts ...grpc.CallOption) (*UpdateListenerDashboardStatResp, error)
		// 查询几天内更新过的XXX列表
		FindListenerListRangeByUpdateTime(ctx context.Context, in *FindListenerListRangeByUpdateTimeReq, opts ...grpc.CallOption) (*FindListenerListRangeByUpdateTimeResp, error)
		// 获取XXX首页统计数据
		GetListenerHomePageDashboard(ctx context.Context, in *GetListenerHomePageDashboardReq, opts ...grpc.CallOption) (*GetListenerHomePageDashboardResp, error)
		// 保存最近多少天的统计数据（一天更新一次，不能覆盖每日更新的数据)
		SnapshotLastDaysListenerStat(ctx context.Context, in *SnapshotLastDaysListenerStatReq, opts ...grpc.CallOption) (*SnapshotLastDaysListenerStatResp, error)
		// 更新XXX每日统计数据的平均值
		UpdateListenerEveryDayAverageStat(ctx context.Context, in *UpdateListenerEveryDayAverageStatReq, opts ...grpc.CallOption) (*UpdateListenerEveryDayAverageStatResp, error)
		// 更新XXX建议
		UpdateListenerSuggestion(ctx context.Context, in *UpdateListenerSuggestionReq, opts ...grpc.CallOption) (*UpdateListenerSuggestionResp, error)
		// 更新XXX的状态
		UpdateListenerOnlineState(ctx context.Context, in *UpdateListenerOnlineStateReq, opts ...grpc.CallOption) (*UpdateListenerOnlineStateResp, error)
		// 生成合同
		GenListenerContract(ctx context.Context, in *GenListenerContractReq, opts ...grpc.CallOption) (*GenListenerContractResp, error)
		// 管理员接口 获取XXX列表
		AdminGetListenerProfileList(ctx context.Context, in *GetListenerProfileListReq, opts ...grpc.CallOption) (*GetListenerProfileListResp, error)
		// 管理员接口 审核XXX
		AdminCheckListenerProfile(ctx context.Context, in *CheckListenerProfileReq, opts ...grpc.CallOption) (*CheckListenerProfileResp, error)
		// 获取提交转账信息并更新状态
		GetCommitMoveCash(ctx context.Context, in *GetCommitMoveCashReq, opts ...grpc.CallOption) (*GetCommitMoveCashResp, error)
		// 获取XXX基本资料（后台内部）
		GetListenerBasicInfo(ctx context.Context, in *GetListenerBasicInfoReq, opts ...grpc.CallOption) (*GetListenerBasicInfoResp, error)
		// 更新新用户推荐XXX
		UpdateRecommendListenerPool(ctx context.Context, in *UpdateRecommendListenerPoolReq, opts ...grpc.CallOption) (*UpdateRecommendListenerPoolResp, error)
		// 获取新用户推荐XXX
		RecListenerWhenUserLogin(ctx context.Context, in *RecListenerWhenUserLoginReq, opts ...grpc.CallOption) (*RecListenerWhenUserLoginResp, error)
		// 管理后台获取XXX提现申请列表
		GetListenerMoveCashListByAdmin(ctx context.Context, in *GetListenerMoveCashListByAdminReq, opts ...grpc.CallOption) (*GetListenerMoveCashListByAdminResp, error)
	}

	defaultListener struct {
		cli zrpc.Client
	}
)

func NewListener(cli zrpc.Client) Listener {
	return &defaultListener{
		cli: cli,
	}
}

// XXX首次填写或更新自己的资料
func (m *defaultListener) AddOrUpdateListenerProfileDraft(ctx context.Context, in *EditListenerProfileDraftReq, opts ...grpc.CallOption) (*EditListenerProfileDraftResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.AddOrUpdateListenerProfileDraft(ctx, in, opts...)
}

// 新申请XXX提交审核
func (m *defaultListener) CommitCheckNewListenerProfile(ctx context.Context, in *CommitCheckNewListenerProfileReq, opts ...grpc.CallOption) (*CommitCheckNewListenerProfileResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.CommitCheckNewListenerProfile(ctx, in, opts...)
}

// 用户查看XXX资料（TODO 后台内部禁止调用）
func (m *defaultListener) GetListenerProfileByUser(ctx context.Context, in *GetListenerProfileByUserReq, opts ...grpc.CallOption) (*GetListenerProfileByUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerProfileByUser(ctx, in, opts...)
}

// XXX获取自己的资料
func (m *defaultListener) GetListenerProfileByOwn(ctx context.Context, in *GetListenerProfileByOwnReq, opts ...grpc.CallOption) (*GetListenerProfileByOwnResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerProfileByOwn(ctx, in, opts...)
}

// XXX备注用户
func (m *defaultListener) ListenerRemarkUser(ctx context.Context, in *ListenerRemarkUserReq, opts ...grpc.CallOption) (*ListenerRemarkUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.ListenerRemarkUser(ctx, in, opts...)
}

// 获取XXX备注的用户
func (m *defaultListener) GetListenerRemarkUser(ctx context.Context, in *GetListenerRemarkUserReq, opts ...grpc.CallOption) (*GetListenerRemarkUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerRemarkUser(ctx, in, opts...)
}

// 用户获取XXX推荐列表
func (m *defaultListener) GetRecommendListenerListByUser(ctx context.Context, in *GetRecommendListenerByUserReq, opts ...grpc.CallOption) (*GetRecommendListenerByUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetRecommendListenerListByUser(ctx, in, opts...)
}

// 获取推荐的XXX
func (m *defaultListener) GetRecommendListenerList(ctx context.Context, in *GetRecommendListenerReq, opts ...grpc.CallOption) (*GetRecommendListenerResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetRecommendListenerList(ctx, in, opts...)
}

// 获取XXX工作状态设置
func (m *defaultListener) GetWorkState(ctx context.Context, in *GetWorkStateReq, opts ...grpc.CallOption) (*GetWorkStateResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetWorkState(ctx, in, opts...)
}

// 修改XXX工作状态
func (m *defaultListener) ChangeWorkState(ctx context.Context, in *ChangeWorkStateReq, opts ...grpc.CallOption) (*ChangeWorkStateResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.ChangeWorkState(ctx, in, opts...)
}

// 获取XXX定价和价格方案
func (m *defaultListener) GetListenerPrice(ctx context.Context, in *GetListenerPriceReq, opts ...grpc.CallOption) (*GetListenerPriceResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerPrice(ctx, in, opts...)
}

// XXX钱包金额更新
func (m *defaultListener) UpdateListenerWallet(ctx context.Context, in *UpdateListenerWalletReq, opts ...grpc.CallOption) (*UpdateListenerWalletResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerWallet(ctx, in, opts...)
}

// 更新XXX的订单统计数据
func (m *defaultListener) UpdateListenerOrderStat(ctx context.Context, in *UpdateListenerOrderStatReq, opts ...grpc.CallOption) (*UpdateListenerOrderStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerOrderStat(ctx, in, opts...)
}

// 绑定银行卡
func (m *defaultListener) SetBankCard(ctx context.Context, in *SetBankCardReq, opts ...grpc.CallOption) (*SetBankCardResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.SetBankCard(ctx, in, opts...)
}

// 获取银行卡
func (m *defaultListener) GetBankCard(ctx context.Context, in *GetBankCardReq, opts ...grpc.CallOption) (*GetBankCardResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetBankCard(ctx, in, opts...)
}

// 获取XXX钱包详情
func (m *defaultListener) GetListenerWallet(ctx context.Context, in *GetListenerWalletReq, opts ...grpc.CallOption) (*GetListenerWalletResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerWallet(ctx, in, opts...)
}

// 获取提现记录
func (m *defaultListener) GetListenerCashLog(ctx context.Context, in *GetListenerCashLogReq, opts ...grpc.CallOption) (*GetListenerCashLogResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerCashLog(ctx, in, opts...)
}

// 获取收益记录
func (m *defaultListener) GetListenerIncomeLog(ctx context.Context, in *GetListenerIncomeLogReq, opts ...grpc.CallOption) (*GetListenerIncomeLogResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerIncomeLog(ctx, in, opts...)
}

// 获取XXX评价统计情况
func (m *defaultListener) GetListenerRatingStat(ctx context.Context, in *GetListenerRatingStatReq, opts ...grpc.CallOption) (*GetListenerRatingStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerRatingStat(ctx, in, opts...)
}

// 获取XXX常用语
func (m *defaultListener) GetListenerWords(ctx context.Context, in *GetListenerWordsReq, opts ...grpc.CallOption) (*GetListenerWordsResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerWords(ctx, in, opts...)
}

// 编辑XXX常用语
func (m *defaultListener) EditListenerWords(ctx context.Context, in *EditListenerWordsReq, opts ...grpc.CallOption) (*EditListenerWordsResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.EditListenerWords(ctx, in, opts...)
}

// 更新XXX与用户的交互情况
func (m *defaultListener) UpdateListenerUserStat(ctx context.Context, in *UpdateListenerUserStatReq, opts ...grpc.CallOption) (*UpdateListenerUserStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerUserStat(ctx, in, opts...)
}

// 更新统计今日推荐和浏览XXX资料页统计
func (m *defaultListener) UpdateTodayListenerUserStat(ctx context.Context, in *UpdateTodayListenerUserStatReq, opts ...grpc.CallOption) (*UpdateTodayListenerUserStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateTodayListenerUserStat(ctx, in, opts...)
}

// 更新XXX首页数据统计看板
func (m *defaultListener) UpdateListenerDashboardStat(ctx context.Context, in *UpdateListenerDashboardStatReq, opts ...grpc.CallOption) (*UpdateListenerDashboardStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerDashboardStat(ctx, in, opts...)
}

// 查询几天内更新过的XXX列表
func (m *defaultListener) FindListenerListRangeByUpdateTime(ctx context.Context, in *FindListenerListRangeByUpdateTimeReq, opts ...grpc.CallOption) (*FindListenerListRangeByUpdateTimeResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.FindListenerListRangeByUpdateTime(ctx, in, opts...)
}

// 获取XXX首页统计数据
func (m *defaultListener) GetListenerHomePageDashboard(ctx context.Context, in *GetListenerHomePageDashboardReq, opts ...grpc.CallOption) (*GetListenerHomePageDashboardResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerHomePageDashboard(ctx, in, opts...)
}

// 保存最近多少天的统计数据（一天更新一次，不能覆盖每日更新的数据)
func (m *defaultListener) SnapshotLastDaysListenerStat(ctx context.Context, in *SnapshotLastDaysListenerStatReq, opts ...grpc.CallOption) (*SnapshotLastDaysListenerStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.SnapshotLastDaysListenerStat(ctx, in, opts...)
}

// 更新XXX每日统计数据的平均值
func (m *defaultListener) UpdateListenerEveryDayAverageStat(ctx context.Context, in *UpdateListenerEveryDayAverageStatReq, opts ...grpc.CallOption) (*UpdateListenerEveryDayAverageStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerEveryDayAverageStat(ctx, in, opts...)
}

// 更新XXX建议
func (m *defaultListener) UpdateListenerSuggestion(ctx context.Context, in *UpdateListenerSuggestionReq, opts ...grpc.CallOption) (*UpdateListenerSuggestionResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerSuggestion(ctx, in, opts...)
}

// 更新XXX的状态
func (m *defaultListener) UpdateListenerOnlineState(ctx context.Context, in *UpdateListenerOnlineStateReq, opts ...grpc.CallOption) (*UpdateListenerOnlineStateResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerOnlineState(ctx, in, opts...)
}

// 生成合同
func (m *defaultListener) GenListenerContract(ctx context.Context, in *GenListenerContractReq, opts ...grpc.CallOption) (*GenListenerContractResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GenListenerContract(ctx, in, opts...)
}

// 管理员接口 获取XXX列表
func (m *defaultListener) AdminGetListenerProfileList(ctx context.Context, in *GetListenerProfileListReq, opts ...grpc.CallOption) (*GetListenerProfileListResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.AdminGetListenerProfileList(ctx, in, opts...)
}

// 管理员接口 审核XXX
func (m *defaultListener) AdminCheckListenerProfile(ctx context.Context, in *CheckListenerProfileReq, opts ...grpc.CallOption) (*CheckListenerProfileResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.AdminCheckListenerProfile(ctx, in, opts...)
}

// 获取提交转账信息并更新状态
func (m *defaultListener) GetCommitMoveCash(ctx context.Context, in *GetCommitMoveCashReq, opts ...grpc.CallOption) (*GetCommitMoveCashResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetCommitMoveCash(ctx, in, opts...)
}

// 获取XXX基本资料（后台内部）
func (m *defaultListener) GetListenerBasicInfo(ctx context.Context, in *GetListenerBasicInfoReq, opts ...grpc.CallOption) (*GetListenerBasicInfoResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerBasicInfo(ctx, in, opts...)
}

// 更新新用户推荐XXX
func (m *defaultListener) UpdateRecommendListenerPool(ctx context.Context, in *UpdateRecommendListenerPoolReq, opts ...grpc.CallOption) (*UpdateRecommendListenerPoolResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateRecommendListenerPool(ctx, in, opts...)
}

// 获取新用户推荐XXX
func (m *defaultListener) RecListenerWhenUserLogin(ctx context.Context, in *RecListenerWhenUserLoginReq, opts ...grpc.CallOption) (*RecListenerWhenUserLoginResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.RecListenerWhenUserLogin(ctx, in, opts...)
}

// 管理后台获取XXX提现申请列表
func (m *defaultListener) GetListenerMoveCashListByAdmin(ctx context.Context, in *GetListenerMoveCashListByAdminReq, opts ...grpc.CallOption) (*GetListenerMoveCashListByAdminResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerMoveCashListByAdmin(ctx, in, opts...)
}