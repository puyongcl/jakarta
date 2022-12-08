package rediskey

// 记录用户锁定了某个用户的通话状态 占线
const RedisKeyRecordCallFromLockTo = "chat:voice_call:lock:target:%d"
