package fwh

// 关注回复消息
const SubscribeReplyMsg = "欢迎关注XXXX服务号！关注XXXX服务号后，你将及时获得最新通知。 <a href=\"https://www.domain.com\" data-miniprogram-appid=\"%s\" data-miniprogram-path=\"pages/msg/msg\">点击立即XX</a>"

// 服务号事件类型
const (
	// EventSubscribe 订阅
	EventSubscribe string = "subscribe"
	// EventUnsubscribe 取消订阅
	EventUnsubscribe string = "unsubscribe"
	// EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan string = "SCAN"
	// EventLocation 上报地理位置事件
	EventLocation string = "LOCATION"
	// EventClick 点击菜单拉取消息时的事件推送
	EventClick string = "CLICK"
	// EventView 点击菜单跳转链接时的事件推送
	EventView string = "VIEW"
	// EventScancodePush 扫码推事件的事件推送
	EventScancodePush string = "scancode_push"
	// EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg string = "scancode_waitmsg"
	// EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto string = "pic_sysphoto"
	// EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum string = "pic_photo_or_album"
	// EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin string = "pic_weixin"
	// EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect string = "location_select"
	// EventTemplateSendJobFinish 发送模板消息推送通知
	EventTemplateSendJobFinish string = "TEMPLATESENDJOBFINISH"
	// EventMassSendJobFinish 群发消息推送通知
	EventMassSendJobFinish string = "MASSSENDJOBFINISH"
	// EventWxaMediaCheck 异步校验图片/音频是否含有违法违规内容推送事件
	EventWxaMediaCheck string = "wxa_media_check"
	// EventSubscribeMsgPopupEvent 订阅通知事件推送
	EventSubscribeMsgPopupEvent string = "subscribe_msg_popup_event"
	// EventPublishJobFinish 发布任务完成
	EventPublishJobFinish string = "PUBLISHJOBFINISH"
)
