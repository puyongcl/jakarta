package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"jakarta/app/mqueue/job/internal/logic/delay"
	"jakarta/app/mqueue/job/internal/logic/schedule/chat"
	"jakarta/app/mqueue/job/internal/logic/schedule/chatorder"
	"jakarta/app/mqueue/job/internal/logic/schedule/listener"
	"jakarta/app/mqueue/job/internal/logic/schedule/stat"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/app/mqueue/job/jobtype"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// register job
func (l *CronJob) Register() *asynq.ServeMux {

	mux := asynq.NewServeMux()

	//scheduler job
	mux.Handle(jobtype.ScheduleCheckVoiceChatOrderExpiry, chatorder.NewCheckVoiceChatOrderExpiryHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleUpdateVoiceChatOrderExpiry, chatorder.NewUpdateVoiceChatOrderExpiryHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleAutoCommentAndFinishChatOrder, chatorder.NewAutoCommentChatOrderHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleAutoConfirmFinishChatOrder, chatorder.NewAutoConfirmFinishChatOrderHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleAutoStartRefundChatOrder, chatorder.NewAutoStartRefundChatOrderHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleAutoAgreeNotProcessRefundApplyChatOrder, chatorder.NewAutoAgreeNotProcessRefundApplyChatOrderHandler(l.svcCtx))
	//
	mux.Handle(jobtype.ScheduleCheckCurrentVoiceChat, chat.NewCheckCurrentVoiceChatHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleCheckCurrentTextChat, chat.NewCheckCurrentTextChatHandler(l.svcCtx))
	//
	mux.Handle(jobtype.ScheduleUpdateListenerUserStat, listener.NewAutoUpdateListenerUserStatHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleUpdateListenerDashboardStat, listener.NewAutoUpdateListenerDashboardHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleAutoClearListenerTodayStatHistoryData, listener.NewAutoClearHistoryDataHandler(l.svcCtx))
	//
	mux.Handle(jobtype.ScheduleUpdateListenerLastDayStat, listener.NewUpdateListenerLastDayStatHandler(l.svcCtx))
	//
	mux.Handle(jobtype.ScheduleUpdateRecommendListenerPool, listener.NewUpdateNewUserRecommendListenerHandler(l.svcCtx))
	//
	mux.Handle(jobtype.ScheduleUpdateDailyStat, stat.NewUpdateDailyStatHandler(l.svcCtx))
	mux.Handle(jobtype.ScheduleAutoUpdateUserStat, stat.NewAutoUpdateUserStatHandler(l.svcCtx))
	//defer job
	mux.Handle(jobtype.DeferCloseChatOrder, delay.NewCloseChatOrderHandler(l.svcCtx))
	mux.Handle(jobtype.DeferSendImMsg, delay.NewDeferSendImMsgHandler(l.svcCtx))
	mux.Handle(jobtype.DeferCheckChatState, delay.NewDeferCheckChatStateHandler(l.svcCtx))

	//queue job , asynq support queue job
	// wait you fill..

	return mux
}
