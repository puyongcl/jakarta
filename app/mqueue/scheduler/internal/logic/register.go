package logic

import (
	"context"
	"jakarta/app/mqueue/scheduler/internal/svc"
)

type MqueueScheduler struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronScheduler(ctx context.Context, svcCtx *svc.ServiceContext) *MqueueScheduler {
	return &MqueueScheduler{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MqueueScheduler) Register() {
	// 更新过期语音订单
	l.updateChatOrderExpiry()
	// 检查即将过期语音订单
	l.checkVoiceChatOrderExpiry()
	// 语音通话倒计时
	l.checkCurrentVoiceChat()
	// 文字聊天倒计时
	l.checkCurrentTextChat()
	// 自动评价
	l.autoCommentChatOrder()
	// 自动确认完成
	l.autoConfirmChatOrder()
	// 自动对已经同意退款的订单发起退款
	l.autoStartRefundChatOrder()
	// 自动同意退款 超过一天未处理的退款申请
	l.autoAgreeNotProcessRefundApplyChatOrder()
	// 间隔几分钟自动更新今日推荐用户数、今日访问个人资料用户数、今日进入聊天界面用户数
	l.autoUpdateListenerUserStat()
	// 每日凌晨更新近7天、30天的订单统计数据
	l.updateListenerLastDayStat()
	// 间隔几分钟自动更新XXX统计数据看板
	l.autoUpdateListenerDashboard()
	// 自动清理不需要的redis数据 今日统计数据
	l.autoClearListenerTodayStatHistoryData()
	// 每日统计
	l.updateDailyStat()
	// 每日更新推荐XXX列表
	l.updateRecommendListenerPool()
	// 自动更新用户的状态统计
	l.autoUpdateUserStat()

	// TODO 根据设置的休息时间对XXX休息状态自动转换(暂不做)
}
