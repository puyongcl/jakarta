package jobtype

// DeferCloseChatOrder 延迟关闭订单
const DeferCloseChatOrder = "defer:chat_order:close"

// 延时发送消息
const DeferSendImMsg = "defer:im:msg:send"

// 延时检查用户聊天状态是否结束
const DeferCheckChatState = "defer:check:chat:state"

// 文字聊天倒计时
const ScheduleCheckCurrentTextChat = "schedule:check:chat:text"

// 语音聊天倒计时
const ScheduleCheckCurrentVoiceChat = "schedule:check:chat:voice"

// 刷新语音聊天订单有效期
const ScheduleUpdateVoiceChatOrderExpiry = "schedule:update:chat_order:voice:expiry"

// 语音订单到期一天前通知
const ScheduleCheckVoiceChatOrderExpiry = "schedule:check:chat_order:voice:expiry"

// 自动好评并确认 每小时查询一次需要自动评价订单
const ScheduleAutoCommentAndFinishChatOrder = "schedule:auto:chat_order:good_comment_and_finish"

// 自动确认 每小时查询一次需要自动确认的订单
const ScheduleAutoConfirmFinishChatOrder = "schedule:auto:chat_order:confirm_finish"

// 自动开始退款
const ScheduleAutoStartRefundChatOrder = "schedule:auto:chat_order:start_refund"

// 自动同意XXX未处理的退款申请
const ScheduleAutoAgreeNotProcessRefundApplyChatOrder = "schedule:auto:chat_order:refund_apply:agree"

// 自动刷新XXX与用户的数据 推荐、访问个人资料、进入聊天页面用户数
const ScheduleUpdateListenerUserStat = "schedule:update:listener:user:stat"

// 定时更新XXX近多少天的统计数据
const ScheduleUpdateListenerLastDayStat = "schedule:update:listener:last_day_stat"

// 自动刷新XXX首页看板统计数据
const ScheduleUpdateListenerDashboardStat = "schedule:update:listener:dashboard:stat"

// 定时清除距离现在第60天的今日XXX数据
const ScheduleAutoClearListenerTodayStatHistoryData = "schedule:auto:clear:listener:today_stat_history"

// 每日统计
const ScheduleUpdateDailyStat = "schedule:update:daily_stat"

// 每日更新推荐XXX列表
const ScheduleUpdateRecommendListenerPool = "schedule:update:recommend:listener:pool"

// 生成真实的评论
const ScheduleGenRealComment = "schedule:gen_real_comment"

// 自动统计用户和XXX的状态
const ScheduleAutoUpdateUserStat = "schedule:auto_update_user_stat"
