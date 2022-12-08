package notify

// 自定义消息类型
const (
	// 订单通知
	DefineNotifyMsgTypeOrderMsg1 = 1 // 用户下单文字XX成功 对XXX 跳转用户聊天
	DefineNotifyMsgTypeOrderMsg2 = 2 // 用户下单语音XX成功 对XXX 跳转用户聊天
	DefineNotifyMsgTypeOrderMsg3 = 3 // 语音XX订单有效期剩余1天时的通知 对用户 跳转XXX聊天
	DefineNotifyMsgTypeOrderMsg4 = 4 // 语音XX订单有效期到期结束时的通知 对用户 跳转XXX聊天
	DefineNotifyMsgTypeOrderMsg5 = 5 // XXX收到用户的评价 对XXX 跳转到订单详情
	DefineNotifyMsgTypeOrderMsg6 = 6 // 用户收到XXX的反馈鼓励通知 对用户 跳转到订单详情

	//
	DefineNotifyMsgTypeOrderRefundMsg7  = 7  // 用户发起退款通知 对用户 跳转到退款单详情
	DefineNotifyMsgTypeOrderRefundMsg8  = 8  // 用户发起退款通知 对XXX 跳转到退款单详情
	DefineNotifyMsgTypeOrderRefundMsg9  = 9  // XXX收到用户申请退款通知 选择 拒绝 通过
	DefineNotifyMsgTypeOrderRefundMsg10 = 10 // 用户收到XXX退款拒绝结果
	DefineNotifyMsgTypeOrderRefundMsg11 = 11 // 用户收到XXX退款通过结果
	DefineNotifyMsgTypeOrderRefundMsg12 = 12 // 用户收到客服介入申请退款 不通过
	DefineNotifyMsgTypeOrderRefundMsg13 = 13 // 用户收到客服介入申请退款 通过
	DefineNotifyMsgTypeOrderRefundMsg14 = 14 // XXX收到 客服介入申请退款 通过
	DefineNotifyMsgTypeOrderRefundMsg15 = 15 // 用户收到 退款到账
	DefineNotifyMsgTypeOrderMsg16       = 16 // 订单已完成 XXX收益到账

	// 浏览通知 只对XXX
	DefineNotifyMsgTypeViewMsg17 = 17 // 用户浏览XXX资料或者进入聊天

	// 系统通知
	DefineNotifyMsgTypeSystemMsg18 = 18 // 用户须知 对新用户
	DefineNotifyMsgTypeSystemMsg19 = 19 // XXX审核不通过
	DefineNotifyMsgTypeSystemMsg20 = 20 // XXX审核通过
	DefineNotifyMsgTypeSystemMsg21 = 21 // XXX资料修改审核通过
	DefineNotifyMsgTypeSystemMsg22 = 22 // XXX资料修改审核不通过
	DefineNotifyMsgTypeSystemMsg23 = 23 // XXX可服务通知 对要求提醒的用户
	DefineNotifyMsgTypeSystemMsg24 = 24 // XXX提现成功通知
	DefineNotifyMsgTypeSystemMsg25 = 25 // XXX提现失败通知
	DefineNotifyMsgTypeSystemMsg26 = 26 // XXX指南 跳转查看指南
	DefineNotifyMsgTypeSystemMsg27 = 27 // 未及时回复消息 自动转为休息中
	DefineNotifyMsgTypeSystemMsg28 = 28 // 对用户 XXX回复XX

	// 聊天消息
	DefineNotifyMsgTypeChatMsg1 = 10001 // 对用户 文字XX订单下单成功，现在可以XX了
	DefineNotifyMsgTypeChatMsg2 = 10002 // 对用户 语音XX订单下单成功，现在可以XX了
	DefineNotifyMsgTypeChatMsg3 = 10003 // 对XXX 用户xx对你下了x分钟文字XX订单，注意及时服务
	DefineNotifyMsgTypeChatMsg4 = 10004 // 对XXX 用户xx对你下了x分钟语音XX订单，注意及时服务
	DefineNotifyMsgTypeChatMsg5 = 10005 // 聊天时长 00：21
	DefineNotifyMsgTypeChatMsg6 = 10006 // 通话接通

	DefineNotifyMsgTypeChatMsg20 = 10020 // 语音聊天时间快用完提醒
	DefineNotifyMsgTypeChatMsg21 = 10021 // 文字聊天时间快用完提醒
	DefineNotifyMsgTypeChatMsg22 = 10022 // 文字服务结束 对用户
	DefineNotifyMsgTypeChatMsg23 = 10023 // 文字服务结束 对XXX
	DefineNotifyMsgTypeChatMsg24 = 10024 // 评价
	DefineNotifyMsgTypeChatMsg25 = 10025 // 反馈
	DefineNotifyMsgTypeChatMsg26 = 10026 // 感谢评价 对用户
	DefineNotifyMsgTypeChatMsg27 = 10027 // 感谢反馈 对XXX
	DefineNotifyMsgTypeChatMsg28 = 10028 // 回复订单评价
	DefineNotifyMsgTypeChatMsg29 = 10029 // 语音服务结束 对用户
	DefineNotifyMsgTypeChatMsg30 = 10030 // 语音服务结束 对XXX
	DefineNotifyMsgTypeChatMsg31 = 10031 // 氛围提示语 仅对用户

	// 小程序服务通知
	DefineNotifyMsgTypeMiniProgramMsg1 = 20001 // 对XXX或用户 留言通知 用户离线时 收到im消息 一次性通知
	DefineNotifyMsgTypeMiniProgramMsg2 = 20002 // 对用户 XXX回复XX 收到的第一条回复 一次性通知
	DefineNotifyMsgTypeMiniProgramMsg3 = 20003 // 对用户 反馈鼓励消息

	// 服务号通知 （暂时只有对XXX）
	DefineNotifyMsgTypeFwhMsg1 = 30001 // 浏览通知 访客接待提醒
	DefineNotifyMsgTypeFwhMsg2 = 30002 // 新订单
	DefineNotifyMsgTypeFwhMsg3 = 30003 // 资料审核通过
	DefineNotifyMsgTypeFwhMsg4 = 30004 // 用户发布XX 通知XXX回复
	DefineNotifyMsgTypeFwhMsg5 = 30005 // XXX被自动转为休息中 通知XXX

	// 发送普通消息
	TextMsgTypeIntroMsg1   = 40001 // 个人介绍
	TextMsgTypeHelloMsg2   = 40002 // 问候语
	TextMsgTypeAdviserMsg3 = 40003 // 用户回复XX顾问消息发送给XXX
)
