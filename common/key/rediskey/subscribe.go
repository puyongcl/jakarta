package rediskey

// 订阅消息成员 订阅对象uid：消息类型：订阅次数
const RedisKeySubscribeNotifyMember = "im:notify:subscribe:%d:%d:%d"

// 一次性订阅消息 消息类型 field:uid value:1
const RedisKeySubscribeOneTimeNotifyMsgMember = "im:notify:subscribe:one_time:%d"
