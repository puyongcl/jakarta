package rediskey

// RedisKeyUidIndex user uid index
const RedisKeyUidIndex = "user:uid:idx"

// CacheTIMUserSignKey /** 用户登陆的IM user sign
const CacheTIMUserSignKey = "user:tim_sign:%d"

// 拉黑列表 set
const RedisKeyUserBlacklist = "user:blacklist:%d"

// 用户每次登陆发送的消息 用一条删除一条
const RedisKeyUserLoginHelloMsg = "user:hello_msg:%d"

// 用户登陆发送消息过期时间 没有发完的消息数据自动过期
const RedisKeyUserLoginHelloMsgExpire = 21 * 24 * 60 * 60
