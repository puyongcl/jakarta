package notify

// 延迟发送新用户IM消息 秒 防止用户尚未注册IM成功
const DeferSendImMsgSecond = 5

// 延迟发送新用户推荐XXX消息
const DeferSendNewUserRecommendListenerImMsgSecond = 6

// 提前通知订单到期 小时
const ScheduleNotifyWillExpiryOrderHour = 24

// 同一个人发送浏览通知的频率 分钟
const SendViewNotifyLimitMin = 3

// 发消息之后延长时间
const SendViewNotifyLimitRenewMin = 2

// 订阅消息发送频率
const (
	SubscribeNotifyMsgSendCntOne    = 1 // 只发送一次
	SubscribeNotifyMsgSendCntAlways = 2 // 总是发送
)

// 发送聊天氛围消息的间隔
const SendChatMsg31IntervalSecond = 2 * 24 * 60 * 60

// 展示最近几天帮助的人
const RecentHelpUserDay = 3

// 发送个人介绍的间隔
const SendListenerIntroMsgIntervalSecond = 7 * 24 * 60 * 60

// 发送问候消息
const SendListenerHelloMsgIntervalSecond = 30 * 60

// 延迟发送新用户推荐XXX消息
const DeferSendUserRecommendListenerImMsgSecond int64 = 10

const DeferSendUserRecommendListenerImMsgIntervalSecond int64 = 8

// 发送氛围提示语延长时间 秒
const DeferSendChatMsg31Second = 3
