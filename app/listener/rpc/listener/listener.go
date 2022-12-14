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
		// XXX????????????????????????????????????
		AddOrUpdateListenerProfileDraft(ctx context.Context, in *EditListenerProfileDraftReq, opts ...grpc.CallOption) (*EditListenerProfileDraftResp, error)
		// ?????????XXX????????????
		CommitCheckNewListenerProfile(ctx context.Context, in *CommitCheckNewListenerProfileReq, opts ...grpc.CallOption) (*CommitCheckNewListenerProfileResp, error)
		// ????????????XXX?????????TODO ???????????????????????????
		GetListenerProfileByUser(ctx context.Context, in *GetListenerProfileByUserReq, opts ...grpc.CallOption) (*GetListenerProfileByUserResp, error)
		// XXX?????????????????????
		GetListenerProfileByOwn(ctx context.Context, in *GetListenerProfileByOwnReq, opts ...grpc.CallOption) (*GetListenerProfileByOwnResp, error)
		// XXX????????????
		ListenerRemarkUser(ctx context.Context, in *ListenerRemarkUserReq, opts ...grpc.CallOption) (*ListenerRemarkUserResp, error)
		// ??????XXX???????????????
		GetListenerRemarkUser(ctx context.Context, in *GetListenerRemarkUserReq, opts ...grpc.CallOption) (*GetListenerRemarkUserResp, error)
		// ????????????XXX????????????
		GetRecommendListenerListByUser(ctx context.Context, in *GetRecommendListenerByUserReq, opts ...grpc.CallOption) (*GetRecommendListenerByUserResp, error)
		// ???????????????XXX
		GetRecommendListenerList(ctx context.Context, in *GetRecommendListenerReq, opts ...grpc.CallOption) (*GetRecommendListenerResp, error)
		// ??????XXX??????????????????
		GetWorkState(ctx context.Context, in *GetWorkStateReq, opts ...grpc.CallOption) (*GetWorkStateResp, error)
		// ??????XXX????????????
		ChangeWorkState(ctx context.Context, in *ChangeWorkStateReq, opts ...grpc.CallOption) (*ChangeWorkStateResp, error)
		// ??????XXX?????????????????????
		GetListenerPrice(ctx context.Context, in *GetListenerPriceReq, opts ...grpc.CallOption) (*GetListenerPriceResp, error)
		// XXX??????????????????
		UpdateListenerWallet(ctx context.Context, in *UpdateListenerWalletReq, opts ...grpc.CallOption) (*UpdateListenerWalletResp, error)
		// ??????XXX?????????????????????
		UpdateListenerOrderStat(ctx context.Context, in *UpdateListenerOrderStatReq, opts ...grpc.CallOption) (*UpdateListenerOrderStatResp, error)
		// ???????????????
		SetBankCard(ctx context.Context, in *SetBankCardReq, opts ...grpc.CallOption) (*SetBankCardResp, error)
		// ???????????????
		GetBankCard(ctx context.Context, in *GetBankCardReq, opts ...grpc.CallOption) (*GetBankCardResp, error)
		// ??????XXX????????????
		GetListenerWallet(ctx context.Context, in *GetListenerWalletReq, opts ...grpc.CallOption) (*GetListenerWalletResp, error)
		// ??????????????????
		GetListenerCashLog(ctx context.Context, in *GetListenerCashLogReq, opts ...grpc.CallOption) (*GetListenerCashLogResp, error)
		// ??????????????????
		GetListenerIncomeLog(ctx context.Context, in *GetListenerIncomeLogReq, opts ...grpc.CallOption) (*GetListenerIncomeLogResp, error)
		// ??????XXX??????????????????
		GetListenerRatingStat(ctx context.Context, in *GetListenerRatingStatReq, opts ...grpc.CallOption) (*GetListenerRatingStatResp, error)
		// ??????XXX?????????
		GetListenerWords(ctx context.Context, in *GetListenerWordsReq, opts ...grpc.CallOption) (*GetListenerWordsResp, error)
		// ??????XXX?????????
		EditListenerWords(ctx context.Context, in *EditListenerWordsReq, opts ...grpc.CallOption) (*EditListenerWordsResp, error)
		// ??????XXX????????????????????????
		UpdateListenerUserStat(ctx context.Context, in *UpdateListenerUserStatReq, opts ...grpc.CallOption) (*UpdateListenerUserStatResp, error)
		// ?????????????????????????????????XXX???????????????
		UpdateTodayListenerUserStat(ctx context.Context, in *UpdateTodayListenerUserStatReq, opts ...grpc.CallOption) (*UpdateTodayListenerUserStatResp, error)
		// ??????XXX????????????????????????
		UpdateListenerDashboardStat(ctx context.Context, in *UpdateListenerDashboardStatReq, opts ...grpc.CallOption) (*UpdateListenerDashboardStatResp, error)
		// ???????????????????????????XXX??????
		FindListenerListRangeByUpdateTime(ctx context.Context, in *FindListenerListRangeByUpdateTimeReq, opts ...grpc.CallOption) (*FindListenerListRangeByUpdateTimeResp, error)
		// ??????XXX??????????????????
		GetListenerHomePageDashboard(ctx context.Context, in *GetListenerHomePageDashboardReq, opts ...grpc.CallOption) (*GetListenerHomePageDashboardResp, error)
		// ?????????????????????????????????????????????????????????????????????????????????????????????)
		SnapshotLastDaysListenerStat(ctx context.Context, in *SnapshotLastDaysListenerStatReq, opts ...grpc.CallOption) (*SnapshotLastDaysListenerStatResp, error)
		// ??????XXX??????????????????????????????
		UpdateListenerEveryDayAverageStat(ctx context.Context, in *UpdateListenerEveryDayAverageStatReq, opts ...grpc.CallOption) (*UpdateListenerEveryDayAverageStatResp, error)
		// ??????XXX??????
		UpdateListenerSuggestion(ctx context.Context, in *UpdateListenerSuggestionReq, opts ...grpc.CallOption) (*UpdateListenerSuggestionResp, error)
		// ??????XXX?????????
		UpdateListenerOnlineState(ctx context.Context, in *UpdateListenerOnlineStateReq, opts ...grpc.CallOption) (*UpdateListenerOnlineStateResp, error)
		// ????????????
		GenListenerContract(ctx context.Context, in *GenListenerContractReq, opts ...grpc.CallOption) (*GenListenerContractResp, error)
		// ??????????????? ??????XXX??????
		AdminGetListenerProfileList(ctx context.Context, in *GetListenerProfileListReq, opts ...grpc.CallOption) (*GetListenerProfileListResp, error)
		// ??????????????? ??????XXX
		AdminCheckListenerProfile(ctx context.Context, in *CheckListenerProfileReq, opts ...grpc.CallOption) (*CheckListenerProfileResp, error)
		// ???????????????????????????????????????
		GetCommitMoveCash(ctx context.Context, in *GetCommitMoveCashReq, opts ...grpc.CallOption) (*GetCommitMoveCashResp, error)
		// ??????XXX??????????????????????????????
		GetListenerBasicInfo(ctx context.Context, in *GetListenerBasicInfoReq, opts ...grpc.CallOption) (*GetListenerBasicInfoResp, error)
		// ?????????????????????XXX
		UpdateRecommendListenerPool(ctx context.Context, in *UpdateRecommendListenerPoolReq, opts ...grpc.CallOption) (*UpdateRecommendListenerPoolResp, error)
		// ?????????????????????XXX
		RecListenerWhenUserLogin(ctx context.Context, in *RecListenerWhenUserLoginReq, opts ...grpc.CallOption) (*RecListenerWhenUserLoginResp, error)
		// ??????????????????XXX??????????????????
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

// XXX????????????????????????????????????
func (m *defaultListener) AddOrUpdateListenerProfileDraft(ctx context.Context, in *EditListenerProfileDraftReq, opts ...grpc.CallOption) (*EditListenerProfileDraftResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.AddOrUpdateListenerProfileDraft(ctx, in, opts...)
}

// ?????????XXX????????????
func (m *defaultListener) CommitCheckNewListenerProfile(ctx context.Context, in *CommitCheckNewListenerProfileReq, opts ...grpc.CallOption) (*CommitCheckNewListenerProfileResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.CommitCheckNewListenerProfile(ctx, in, opts...)
}

// ????????????XXX?????????TODO ???????????????????????????
func (m *defaultListener) GetListenerProfileByUser(ctx context.Context, in *GetListenerProfileByUserReq, opts ...grpc.CallOption) (*GetListenerProfileByUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerProfileByUser(ctx, in, opts...)
}

// XXX?????????????????????
func (m *defaultListener) GetListenerProfileByOwn(ctx context.Context, in *GetListenerProfileByOwnReq, opts ...grpc.CallOption) (*GetListenerProfileByOwnResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerProfileByOwn(ctx, in, opts...)
}

// XXX????????????
func (m *defaultListener) ListenerRemarkUser(ctx context.Context, in *ListenerRemarkUserReq, opts ...grpc.CallOption) (*ListenerRemarkUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.ListenerRemarkUser(ctx, in, opts...)
}

// ??????XXX???????????????
func (m *defaultListener) GetListenerRemarkUser(ctx context.Context, in *GetListenerRemarkUserReq, opts ...grpc.CallOption) (*GetListenerRemarkUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerRemarkUser(ctx, in, opts...)
}

// ????????????XXX????????????
func (m *defaultListener) GetRecommendListenerListByUser(ctx context.Context, in *GetRecommendListenerByUserReq, opts ...grpc.CallOption) (*GetRecommendListenerByUserResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetRecommendListenerListByUser(ctx, in, opts...)
}

// ???????????????XXX
func (m *defaultListener) GetRecommendListenerList(ctx context.Context, in *GetRecommendListenerReq, opts ...grpc.CallOption) (*GetRecommendListenerResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetRecommendListenerList(ctx, in, opts...)
}

// ??????XXX??????????????????
func (m *defaultListener) GetWorkState(ctx context.Context, in *GetWorkStateReq, opts ...grpc.CallOption) (*GetWorkStateResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetWorkState(ctx, in, opts...)
}

// ??????XXX????????????
func (m *defaultListener) ChangeWorkState(ctx context.Context, in *ChangeWorkStateReq, opts ...grpc.CallOption) (*ChangeWorkStateResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.ChangeWorkState(ctx, in, opts...)
}

// ??????XXX?????????????????????
func (m *defaultListener) GetListenerPrice(ctx context.Context, in *GetListenerPriceReq, opts ...grpc.CallOption) (*GetListenerPriceResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerPrice(ctx, in, opts...)
}

// XXX??????????????????
func (m *defaultListener) UpdateListenerWallet(ctx context.Context, in *UpdateListenerWalletReq, opts ...grpc.CallOption) (*UpdateListenerWalletResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerWallet(ctx, in, opts...)
}

// ??????XXX?????????????????????
func (m *defaultListener) UpdateListenerOrderStat(ctx context.Context, in *UpdateListenerOrderStatReq, opts ...grpc.CallOption) (*UpdateListenerOrderStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerOrderStat(ctx, in, opts...)
}

// ???????????????
func (m *defaultListener) SetBankCard(ctx context.Context, in *SetBankCardReq, opts ...grpc.CallOption) (*SetBankCardResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.SetBankCard(ctx, in, opts...)
}

// ???????????????
func (m *defaultListener) GetBankCard(ctx context.Context, in *GetBankCardReq, opts ...grpc.CallOption) (*GetBankCardResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetBankCard(ctx, in, opts...)
}

// ??????XXX????????????
func (m *defaultListener) GetListenerWallet(ctx context.Context, in *GetListenerWalletReq, opts ...grpc.CallOption) (*GetListenerWalletResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerWallet(ctx, in, opts...)
}

// ??????????????????
func (m *defaultListener) GetListenerCashLog(ctx context.Context, in *GetListenerCashLogReq, opts ...grpc.CallOption) (*GetListenerCashLogResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerCashLog(ctx, in, opts...)
}

// ??????????????????
func (m *defaultListener) GetListenerIncomeLog(ctx context.Context, in *GetListenerIncomeLogReq, opts ...grpc.CallOption) (*GetListenerIncomeLogResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerIncomeLog(ctx, in, opts...)
}

// ??????XXX??????????????????
func (m *defaultListener) GetListenerRatingStat(ctx context.Context, in *GetListenerRatingStatReq, opts ...grpc.CallOption) (*GetListenerRatingStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerRatingStat(ctx, in, opts...)
}

// ??????XXX?????????
func (m *defaultListener) GetListenerWords(ctx context.Context, in *GetListenerWordsReq, opts ...grpc.CallOption) (*GetListenerWordsResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerWords(ctx, in, opts...)
}

// ??????XXX?????????
func (m *defaultListener) EditListenerWords(ctx context.Context, in *EditListenerWordsReq, opts ...grpc.CallOption) (*EditListenerWordsResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.EditListenerWords(ctx, in, opts...)
}

// ??????XXX????????????????????????
func (m *defaultListener) UpdateListenerUserStat(ctx context.Context, in *UpdateListenerUserStatReq, opts ...grpc.CallOption) (*UpdateListenerUserStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerUserStat(ctx, in, opts...)
}

// ?????????????????????????????????XXX???????????????
func (m *defaultListener) UpdateTodayListenerUserStat(ctx context.Context, in *UpdateTodayListenerUserStatReq, opts ...grpc.CallOption) (*UpdateTodayListenerUserStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateTodayListenerUserStat(ctx, in, opts...)
}

// ??????XXX????????????????????????
func (m *defaultListener) UpdateListenerDashboardStat(ctx context.Context, in *UpdateListenerDashboardStatReq, opts ...grpc.CallOption) (*UpdateListenerDashboardStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerDashboardStat(ctx, in, opts...)
}

// ???????????????????????????XXX??????
func (m *defaultListener) FindListenerListRangeByUpdateTime(ctx context.Context, in *FindListenerListRangeByUpdateTimeReq, opts ...grpc.CallOption) (*FindListenerListRangeByUpdateTimeResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.FindListenerListRangeByUpdateTime(ctx, in, opts...)
}

// ??????XXX??????????????????
func (m *defaultListener) GetListenerHomePageDashboard(ctx context.Context, in *GetListenerHomePageDashboardReq, opts ...grpc.CallOption) (*GetListenerHomePageDashboardResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerHomePageDashboard(ctx, in, opts...)
}

// ?????????????????????????????????????????????????????????????????????????????????????????????)
func (m *defaultListener) SnapshotLastDaysListenerStat(ctx context.Context, in *SnapshotLastDaysListenerStatReq, opts ...grpc.CallOption) (*SnapshotLastDaysListenerStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.SnapshotLastDaysListenerStat(ctx, in, opts...)
}

// ??????XXX??????????????????????????????
func (m *defaultListener) UpdateListenerEveryDayAverageStat(ctx context.Context, in *UpdateListenerEveryDayAverageStatReq, opts ...grpc.CallOption) (*UpdateListenerEveryDayAverageStatResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerEveryDayAverageStat(ctx, in, opts...)
}

// ??????XXX??????
func (m *defaultListener) UpdateListenerSuggestion(ctx context.Context, in *UpdateListenerSuggestionReq, opts ...grpc.CallOption) (*UpdateListenerSuggestionResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerSuggestion(ctx, in, opts...)
}

// ??????XXX?????????
func (m *defaultListener) UpdateListenerOnlineState(ctx context.Context, in *UpdateListenerOnlineStateReq, opts ...grpc.CallOption) (*UpdateListenerOnlineStateResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateListenerOnlineState(ctx, in, opts...)
}

// ????????????
func (m *defaultListener) GenListenerContract(ctx context.Context, in *GenListenerContractReq, opts ...grpc.CallOption) (*GenListenerContractResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GenListenerContract(ctx, in, opts...)
}

// ??????????????? ??????XXX??????
func (m *defaultListener) AdminGetListenerProfileList(ctx context.Context, in *GetListenerProfileListReq, opts ...grpc.CallOption) (*GetListenerProfileListResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.AdminGetListenerProfileList(ctx, in, opts...)
}

// ??????????????? ??????XXX
func (m *defaultListener) AdminCheckListenerProfile(ctx context.Context, in *CheckListenerProfileReq, opts ...grpc.CallOption) (*CheckListenerProfileResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.AdminCheckListenerProfile(ctx, in, opts...)
}

// ???????????????????????????????????????
func (m *defaultListener) GetCommitMoveCash(ctx context.Context, in *GetCommitMoveCashReq, opts ...grpc.CallOption) (*GetCommitMoveCashResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetCommitMoveCash(ctx, in, opts...)
}

// ??????XXX??????????????????????????????
func (m *defaultListener) GetListenerBasicInfo(ctx context.Context, in *GetListenerBasicInfoReq, opts ...grpc.CallOption) (*GetListenerBasicInfoResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerBasicInfo(ctx, in, opts...)
}

// ?????????????????????XXX
func (m *defaultListener) UpdateRecommendListenerPool(ctx context.Context, in *UpdateRecommendListenerPoolReq, opts ...grpc.CallOption) (*UpdateRecommendListenerPoolResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.UpdateRecommendListenerPool(ctx, in, opts...)
}

// ?????????????????????XXX
func (m *defaultListener) RecListenerWhenUserLogin(ctx context.Context, in *RecListenerWhenUserLoginReq, opts ...grpc.CallOption) (*RecListenerWhenUserLoginResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.RecListenerWhenUserLogin(ctx, in, opts...)
}

// ??????????????????XXX??????????????????
func (m *defaultListener) GetListenerMoveCashListByAdmin(ctx context.Context, in *GetListenerMoveCashListByAdminReq, opts ...grpc.CallOption) (*GetListenerMoveCashListByAdminResp, error) {
	client := pb.NewListenerClient(m.cli.Conn())
	return client.GetListenerMoveCashListByAdmin(ctx, in, opts...)
}
