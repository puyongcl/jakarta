package chatkey

// 每日登陆每位XXX免费文字聊天次数
const FreeTextChatCntPerDay = 5

// 下单语音赠送免费消息
const GiftTextChatCnt = 10

// 同步聊天状态接口
const (
	ChatAction1  int64 = 1  // 未通话之前 进入聊天界面 获取当前免费文字聊天次数 检查通话状态
	ChatAction2  int64 = 2  // 文字聊天次数减1
	ChatAction3  int64 = 3  // 普通用户开始语音聊天
	ChatAction4  int64 = 4  // 普通用户结束语音聊天（主动结束）
	ChatAction5  int64 = 5  // 普通用户结束语音聊天（未接通）
	ChatAction6  int64 = 6  // 普通用户结束语音聊天（意外中断）
	ChatAction7  int64 = 7  // XXX开始语音聊天
	ChatAction8  int64 = 8  // XXX结束语音聊天（主动结束）
	ChatAction9  int64 = 9  // XXX结束语音聊天（未接通）
	ChatAction10 int64 = 10 // XXX结束语音聊天（意外中断)
	ChatAction11 int64 = 11 // 检测通话状态是否异常
	ChatAction12 int64 = 12 // 进入聊天界面之后 获取用户当前的通话可用时间
	ChatAction13 int64 = 13 // 拨打语音通话
	ChatAction14 int64 = 14 // 解除当前用户的用户通话对象锁定占用
)

// 当前语音聊天状态
const (
	VoiceChatStateStart  int64 = 1 // 已开始
	VoiceChatStateStop   int64 = 2 // 已结束
	VoiceChatStateSettle int64 = 3 // 已结算
)

// 当前文字聊天状态
const (
	TextChatStateStart       int64 = 1 // 已开始
	TextChatStateStop        int64 = 2 // 已结束
	TextChatStateAlreadyStop int64 = 3 // 当前已经更新结束 不需要再次更新
)

// 定时检查语音聊天时间结束 秒
const CheckVoiceChatUseOutIntervalSecond = 10

// 定时检查文字聊天时间结束 分钟
const CheckTextChatUseOutIntervalMinute = 1

// 提前几分钟通知
const NotifyOrderOverBeforeMinute = 1

// 检查聊天到期时间范围
const CheckChatUseOutRange = 2

// 聊天服务可用额度更新操作类型
const (
	ChatStatUpdateTypeOrderPaidAdd    = 1 // 支付成功增加
	ChatStatUpdateTypeOrderExpireDecr = 2 // 订单过期减少
	ChatStatUpdateTypeVoiceChatDecr   = 3 // 通话使用扣除
	ChatStatUpdateTypeOrderUserStop   = 4 // 订单过期减少
)

// 延时检查用户的通话 分钟
const DeferCheckVoiceChatMinute = 1

// 延迟检查文字聊天是否结束 分钟
const DeferCheckTextChatMinute = 1

// 用户聊天状态
const (
	UserChatState1 int64 = 1 // 新用户
	UserChatState2 int64 = 2 // 服务中用户
	UserChatState3 int64 = 3 // 已结束服务用户

)
