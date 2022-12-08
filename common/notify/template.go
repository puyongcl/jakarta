package notify

// 自定义消息模版
const (
	DefineNotifyMsgTemplateOrderMsgTitle1 = "请及时服务"
	DefineNotifyMsgTemplateOrderMsg1      = "%s对你下了%d分钟文字XX订单，注意及时服务。"
	DefineNotifyMsgTemplateOrderMsgTitle2 = "请及时服务"
	DefineNotifyMsgTemplateOrderMsg2      = "%s对你下了%d分钟语音XX订单，注意及时服务。"
	DefineNotifyMsgTemplateOrderMsgTitle3 = "语音XX服务有效期提醒"
	DefineNotifyMsgTemplateOrderMsg3      = "语音XX服务有效期剩余1天，请尽快使用。"
	DefineNotifyMsgTemplateOrderMsgTitle4 = "语音XX服务有效期结束"
	DefineNotifyMsgTemplateOrderMsg4      = "语音XX服务有效期结束了，如有需要，请重新下单。"
	DefineNotifyMsgTemplateOrderMsgTitle5 = "服务评价"
	DefineNotifyMsgTemplateOrderMsg5      = "%s对你的服务%s，加油，尽你所能去提供更好的服务。"
	DefineNotifyMsgTemplateOrderMsgTitle6 = "反馈鼓励"
	DefineNotifyMsgTemplateOrderMsg6      = "XXX%s给你送来了反馈鼓励，去查看。"

	DefineNotifyMsgTemplateOrderMsgTitle7         = "退款申请已通过"
	DefineNotifyMsgTemplateOrderRefundMsg7        = "符合10分钟内结束服务，3次无条件退款机会，已自动同意你单退款申请。"
	DefineNotifyMsgTemplateOrderMsgTitle8         = "退款申请已自动同意"
	DefineNotifyMsgTemplateOrderRefundMsg8        = "%s对%d月%d日购买对订单（订单号%s）申请退款，退款理由是%s。用户退款申请符合10分钟内结束服务，3次无条件退款机会，已自动同意用户的退款申请。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle9  = "退款申请处理"
	DefineNotifyMsgTemplateOrderRefundMsg9        = "%s对%d月%d日购买对订单（订单号%s）申请退款，退款理由是%s。你可与用户进行协商，若无法达成一致，用户可以申请客服介入。如不处理，1天后将自动同意用户的退款申请。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle10 = "XXX拒绝了你的退款申请"
	DefineNotifyMsgTemplateOrderRefundMsg10       = "XXX拒绝的理由：%s。你可以和XXX沟通处理，也可以申请客服介入，补充更多信息。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle11 = "退款申请已通过"
	DefineNotifyMsgTemplateOrderRefundMsg11       = "XXX同意了你的退款申请，退款金额会在1-2个工作日内到账。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle12 = "退款申请不通过"
	DefineNotifyMsgTemplateOrderRefundMsg12       = "由于不符合10分钟内结束服务，3次无条件退款机会，且服务使用时长较多，无法证明XXX的服务存在问题，退款审核不通过。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle13 = "客服介入退款申请已通过"
	DefineNotifyMsgTemplateOrderRefundMsg13       = "客服同意了你的退款申请，退款会在1-2个工作日内到账。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle14 = "客服介入已通过退款申请"
	DefineNotifyMsgTemplateOrderRefundMsg14       = "%s对%d月%d日购买对订单（订单号%s）申请退款，退款理由是%s。客服根据用户提供的信息，同意了退款申请，如有异议，请联系客服。"
	DefineNotifyMsgTemplateOrderRefundMsgTittle15 = "你的退款已到账"
	DefineNotifyMsgTemplateOrderRefundMsg15       = "你的退款已经返回原支付渠道，请查询确认。"
	DefineNotifyMsgTemplateOrderMsgTittle16       = "钱包收益增加通知"
	DefineNotifyMsgTemplateOrderMsg16             = "用户已确认完成订单（订单号%s），你的钱包收益增加%s元。"

	DefineNotifyMsgTemplateViewMsgTitle17 = "浏览通知"
	DefineNotifyMsgTemplateViewMsg17      = "%s浏览了你，可能有意向你XX，可适度关心用户，增加下单机会。"

	// 系统通知
	DefineNotifyMsgTemplateSystemMsgTitle18 = "用户须知"
	DefineNotifyMsgTemplateSystemMsg18      = "欢迎你来到XXXX，感谢和你的相遇。\nXXXX将帮助你挑选合适的XXX，你可以查看XXX的资质证书、服务时长、用户评价等信息，找到合适你的XXX，获得理解和关怀，听取专业分析和过来人建议指导。\n\n我们将为你提供保障和服务承诺：\n*每日免费5句话（和每个XXX）\n*XX\n*10分钟内结束服务，服务不满意3次内无条件退款\n*全面隐私保护，匿名账号\n*严选优质XXX，严格入驻标准\n\n同时，请你注意，XXXX平台禁止一切淫秽、暴力、辱骂、违法等内容和行为，如你违反，将会被限时冻结账号、封号，感谢你对XXXX的信任与支持。"
	DefineNotifyMsgTemplateSystemMsgTitle19 = "成为XXX审核不通过"
	DefineNotifyMsgTemplateSystemMsg19      = "%s不符合要求：%s，请重新修改后再次提交。"
	DefineNotifyMsgTemplateSystemMsgTitle20 = "恭喜你成为XXX"
	DefineNotifyMsgTemplateSystemMsg20      = "你已通过平台的审核成为XXX，可切换到服务中状态及早开始接单服务，另外敬请严格遵守XXX服务协议和XXX承诺书。"
	DefineNotifyMsgTemplateSystemMsgTitle21 = "审核通过"
	DefineNotifyMsgTemplateSystemMsg21      = "你提交的%s，审核通过，感谢你对XX的支持。"
	DefineNotifyMsgTemplateSystemMsgTitle22 = "审核不通过"
	DefineNotifyMsgTemplateSystemMsg22      = "你提交的%s，没有审核通过%s，请修改后再次提交。"
	DefineNotifyMsgTemplateSystemMsgTitle23 = "XXX可服务通知"
	DefineNotifyMsgTemplateSystemMsg23      = "XXX%s可接单，现在你可以下单了。"
	DefineNotifyMsgTemplateSystemMsgTitle24 = "提现成功"
	DefineNotifyMsgTemplateSystemMsg24      = "你申请的提现成功到账，请查收。"
	DefineNotifyMsgTemplateSystemMsgTitle25 = "提现失败"
	DefineNotifyMsgTemplateSystemMsg25      = "你申请的提现失败，失败原因:%s。"
	DefineNotifyMsgTemplateSystemMsgTitle26 = "建议查看XXX指南"
	DefineNotifyMsgTemplateSystemMsg26      = "建议关注服务号和查看XX指南\n建议搜索关注服务号“XXXX”，获得用户来访下单通知，便于及时服务用户。\n建议查阅XX指南，XX平台为XXX整理了各项学习提升事项，将帮助你更好去服务用户。"
	DefineNotifyMsgTemplateSystemMsgTitle27 = "已被切换到休息中状态"
	DefineNotifyMsgTemplateSystemMsg27      = "由于没有在4分钟内回复用户，可接单状态已被切换到休息中状态。请在休息时段切换到休息中状态。"
	DefineNotifyMsgTemplateSystemMsgTitle28 = "%s回复了你的%sXX"
	DefineNotifyMsgTemplateSystemMsg28      = "" // 回复内容开头

	// 聊天消息
	DefineNotifyMsgTemplateChatMsg1 = "文字XX下单成功，现在可以XX了"
	DefineNotifyMsgTemplateChatMsg2 = "语音XX下单成功，现在可以XX了"
	DefineNotifyMsgTemplateChatMsg3 = "%s对你下了%d分钟文字XX订单，注意及时服务"
	DefineNotifyMsgTemplateChatMsg4 = "%s对你下了%d分钟语音XX订单，注意及时服务"
	DefineNotifyMsgTemplateChatMsg5 = "本次通话时长：%d分钟"
	DefineNotifyMsgTemplateChatMsg6 = "通话已接通"

	DefineNotifyMsgTemplateChatMsg20  = "语音通话还剩1分钟，请抓紧时间沟通"
	DefineNotifyMsgTemplateChatMsg21  = "文字订单还剩1分钟，请抓紧时间沟通"
	DefineNotifyMsgTemplateChatMsg22  = "服务结束，感谢你的下单" // 文字聊天结束 对用户
	DefineNotifyMsgTemplateChatMsg23  = "服务结束，感谢你的服务" // 文字聊天结束 对XXX
	DefineNotifyMsgTemplateChatMsg24  = ""            // 直接发内容
	DefineNotifyMsgTemplateChatMsg25  = ""            // 直接发内容
	DefineNotifyMsgTemplateChatMsg26  = "感谢你的评价"      // 直接发内容
	DefineNotifyMsgTemplateChatMsg27  = "感谢你的反馈鼓励"    // 直接发内容
	DefineNotifyMsgTemplateChatMsg28  = "XXX%s回复了你，去查看"
	DefineNotifyMsgTemplateChatMsg29  = "服务结束，感谢你的下单" // 语音聊天结束 对用户
	DefineNotifyMsgTemplateChatMsg30  = "服务结束，感谢你的服务" // 语音聊天结束 对XXX
	DefineNotifyMsgTemplateChatMsg31A = "%s在%s前收到一个满意评价"
	DefineNotifyMsgTemplateChatMsg31B = "%s在最近%d天帮助了%d人"

	// 小程序通知
	DefineNotifyMsgTypeMiniProgramMsg1VoiceMsgDefaultContent = "你收到了一条语音消息，请打开小程序查看"
	DefineNotifyMsgTypeMiniProgramMsg1ImageMsgDefaultContent = "你收到了一条图片消息，请打开小程序查看"
	DefineNotifyMsgTypeMiniProgramMsg1Path                   = "pages/msg/msg"
	DefineNotifyMsgTypeMiniProgramMsg1Color                  = "#DE8C00"

	DefineNotifyMsgTypeMiniProgramMsg2Path  = "pages/story/story-reply?id=%s"
	DefineNotifyMsgTypeMiniProgramMsg2Color = "#005ADE"

	DefineNotifyMsgTypeMiniProgramMsg3Path  = "pages/msg/msg"
	DefineNotifyMsgTypeMiniProgramMsg3Color = "#2B6AC6"

	// 服务号通知
	DefineNotifyMsgTypeFwhMsg1FirstData  = "刚刚有一个用户查看了你的资料"
	DefineNotifyMsgTypeFwhMsg1RemarkData = "点击查看访客详情"
	DefineNotifyMsgTypeFwhMsg1Path       = "pages/msg/msg"
	DefineNotifyMsgTypeFwhMsg1Color      = "#008BDE"

	DefineNotifyMsgTypeFwhMsg2FirstData     = "收到了一个新订单"
	DefineNotifyMsgTypeFwhMsg2RemarkData    = "\n进入小程序XX用户，避免因服务不及时导致用户退款。立即前往>>"
	DefineNotifyMsgTypeFwhMsg2Path          = "pages/my/index"
	DefineNotifyMsgTypeFwhMsg2OrderTyeText  = "文字XX%d分钟"
	DefineNotifyMsgTypeFwhMsg2OrderTyeVoice = "语音XX%d分钟"
	DefineNotifyMsgTypeFwhMsg2OrderInfo     = "订单金额%s元"
	DefineNotifyMsgTypeFwhMsg2Color         = "#FF8000"

	DefineNotifyMsgTypeFwhMsg3FirstData1 = "恭喜你成为XXX"
	DefineNotifyMsgTypeFwhMsg3FirstData2 = "请修改后再申请"
	DefineNotifyMsgTypeFwhMsg3FirstData3 = "恭喜你，审核通过"
	DefineNotifyMsgTypeFwhMsg3FirstData4 = "审核不通过，请修改后再申请"
	DefineNotifyMsgTypeFwhMsg3RemarkData = "点击查看详情"
	DefineNotifyMsgTypeFwhMsg3Path       = "pages/my/my"
	DefineNotifyMsgTypeFwhMsg3Color      = "#2152DD"

	DefineNotifyMsgTypeFwhMsg4FirstData  = "有一个新的%sXX待回复"
	DefineNotifyMsgTypeFwhMsg4RemarkData = "主动回复可增加接单机会"
	DefineNotifyMsgTypeFwhMsg4Path       = "pages/story/story-detail?id=%s"
	DefineNotifyMsgTypeFwhMsg4Color      = "#00DE4E"

	DefineNotifyMsgTypeFwhMsg5FirstData  = "你没有在%d分钟内回复用户消息，系统自动将你切换到休息中状态"
	DefineNotifyMsgTypeFwhMsg5Key1       = "XXX"
	DefineNotifyMsgTypeFwhMsg5Key2       = "接单状态变为休息中"
	DefineNotifyMsgTypeFwhMsg5RemarkData = "\n点击打开小程序，切换到接单状态"
	DefineNotifyMsgTypeFwhMsg5Path       = "pages/my/index"
	DefineNotifyMsgTypeFwhMsg5Color      = "#FF7000"

	// 服务号 小程序通知 统一备注颜色
	MpNotifyMsgRemarkDefaultColor = "#FF7000"
)
