package rediskey

// 消息发送频率限制\消息去重 uid+targetUid+msgType value:timestamp
const RedisKeyNotifyMsgLimit = "im:notify:limit:%s:%d"

// 记录每个人每天发送消息数 uid:年月日
const RedisKeyStatSendMsgCntToday = "im:stat:msg_cnt:%d:%s"

// 记录XXX回复消息的最后时间
const RedisKeyMarkReplyMsgTime = "im:mark:msg_reply_time"
