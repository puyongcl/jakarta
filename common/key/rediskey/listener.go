package rediskey

// XXX编辑资料 审核中的字段 废弃
const RedisKeyListenerProfileFieldCheckStatusChecking = "listener:listener:profile:field:checking:%d"

// 每位XXX每日免费的次数 uid:listenerUid
const RedisKeyUserPerListenerFreeTextMsgCnt = "listener:user:free:text_msg:%d:%d"

// 评价标签统计
const RedisKeyListenerCommentTagStat = "listener:comment_tag:%d"

// 每日推荐XXX给未下单用户 有序集合 uid 推荐最大次数
const RedisKeyListenerRecommendWhenUserLogin = "listener:recommend:new_user:%s"

// 每日推荐新用户聊天 有序集合 uid 推荐最大次数
const RedisKeyListenerRecommendSendHelloMsgWhenUserLogin = "listener:recommend:send_hello_msg:%s"

// 每日推荐XXX回复用户发布XX uid 推荐次数
const RedisKeyListenerRecommendReplyUserStory = "listener:recommend:reply_user_story:%s"

// 需要清理的数据
var NeedDelRedisRecommendData = []string{RedisKeyListenerRecommendWhenUserLogin, RedisKeyListenerRecommendReplyUserStory, RedisKeyListenerRecommendSendHelloMsgWhenUserLogin}
