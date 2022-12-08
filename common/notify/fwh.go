package notify

// 服务号订阅模版消息
const (
	// 访客接待提醒
	//{{first.DATA}}
	//访客昵称：{{keyword1.DATA}}
	//来访时间：{{keyword2.DATA}}
	//{{remark.DATA}}
	// 有访客需要您前往接待
	//访客昵称：洛唐小白
	//来访时间：2018年6月28日
	//洛唐提示：贵公司有客户到访，请尽快联系并前往接待

	FwhTemplateMsg1 = "Nw52mSOwIJaUJediiO1Y4csTvf1Lqu_pUCUQ1FTuVmU"

	// 新订单
	// {{first.DATA}}
	//提交时间：{{keyword1.DATA}}
	//订单类型：{{keyword2.DATA}}
	//客户信息：{{keyword3.DATA}}
	//订单信息：{{keyword4.DATA}}
	//{{remark.DATA}}
	// 您收到了一条新的订单
	//提交时间：2018年2月29日
	//订单类型：微信订单
	//客户信息：张三
	//订单信息：一盒新鲜空气
	//截止当前，您有X条订单未处理

	FwhTemplateMsg2 = "gtFlSEzWYXGjegDcs9mR8q5ML4OuXBTWysYAFs8TK08"

	// 资料审核
	// {{first.DATA}}
	//姓名：{{keyword1.DATA}}
	//手机号：{{keyword2.DATA}}.
	//审核时间：{{keyword3.DATA}}
	//{{remark.DATA}}
	// 您好，您的资料已经通过审核。
	//姓名：小明
	//手机号：13866668888
	//审核时间：2014年7月21日 18:36
	//感谢你的使用。

	FwhTemplateMsg3 = "xksUwf_dJL68no-kbUjRidXY3rv4mQCOzSrRvVxBxh0"

	// 用户发布XX 通知XXX回复
	//详细内容
	//{{first.DATA}}
	//问题内容：{{keyword1.DATA}}
	//问题类型：{{keyword2.DATA}}
	//发生时间：{{keyword3.DATA}}
	//{{remark.DATA}}

	FwhTemplateMsg4 = "xOCiQIn8_RbKIRQgbMu___j5NoGqY1LpauTxh2ZwmUo"

	// 通知XXX被转为休息中
	// {{first.DATA}}
	//用户身份：{{keyword1.DATA}}
	//状态：{{keyword2.DATA}}
	//时间：{{keyword3.DATA}}
	//{{remark.DATA}}
	FwhTemplateMsg5 = "ktVuFGDlXQM23Xy8qi-ESrLeYY88-XSHWwm1Px9LVWw"
)
