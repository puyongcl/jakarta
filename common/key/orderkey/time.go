package orderkey

// 有效期
const VoiceChatOrderExpireHour = 72

// 定时更新订单有效期间隔 分钟
const UpdateOrderExpiryIntervalMinutes = 10

// 定时检查即将过期订单
const CheckOrderExpiryIntervalMinutes = 6

// 关闭未支付订单的时间 分钟
const CloseOrderTimeMinutes = 2

// 限制多少时间内不能存在未支付订单 防止频繁错误下单
const NoPayOrderLimitSecond = 10

// XXX拒绝退款 不确认完成 几天自动转已完成
const AutoFinishAfterListenerRefuseRefund = 1

// 评价后不确认完成 几天后自动转已完成
const AutoFinishAfterNotGoodCommentDay = 1

// 自动确认任务频率 小时
const AutoConfirmFinishScheduleIntervalHour = 1

// 订单服务结束后不评价 几天后自动好评并完成
const AutoGoodCommentAfterStopDay = 1

// 自动好评任务频率 小时
const AutoGoodCommentScheduleIntervalHour = 1

// 订单结束后 几天内允许申请退款
const CanApplyRefundAfterStopDay = 1

// XXX几天内不处理退款 自动同意
const AutoAgreeNotProcessRefundApplyChatOrderDay = 1

// 自动同意XXX未处理退款申请任务频率 小时
const AutoAgreeNotProcessRefundApplyChatOrderDayIntervalHour = 1

// XXX拒绝退款后 几天内 可以再次发起客服介入
const CanSecondApplyRefundAfterRefuseDay = 1

// 发起客服介入退款 几天内不处理 自动拒绝退款 自动确认完成
const AutoFinishAfterSecondApplyRefundDay = 1

// 对同意退款订单 几天后发起退款支付
const AutoStartRefundDay = 1

// 自动发起退款任务频率 小时
const AutoStartRefundScheduleIntervalHour = 1

// 定时更新近几天XXX统计数据 定时任务时间点
const UpdateLastDayListenerStatHour = 0

// 定时清理过去距今第60天的今日XXX统计数据 定时任务时间点
const ClearListenerTodayStatHistoryDataHour = 2

// 文字订单由于取整分钟时间 额外补偿用户的时间 分钟 （设置为0 因为到期时间向后取整）
const TextOrderExtraAddMinute = 0
