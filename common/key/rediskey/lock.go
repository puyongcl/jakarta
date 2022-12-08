package rediskey

// XXX修改资料加锁
const RedisLockEditListenerProfile = "listener:edit:listener:%d"

// 拨打电话加锁
const RedisLockStartCallVoiceChat = "chat:call:listener:%d"

// 上传腾讯云合同文件加锁
const RedisLockUploadCosContract = "admin:listener:upload:contract:%d"

// 用户关键操作锁
const RedisLockUser = "bbs:user:request:%d"

// 生成1021合同操作锁
const RedisLockGenContract1021 = "admin:contract:gen:%s"

// 用户注册锁
const RedisLockRegUser = "bbs:user:reg:%s"
