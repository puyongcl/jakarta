package rediskey

// 今日接单量 key_prefix:20220804 zset XXX 接单量
const RedisKeyListenerTodayOrderCnt = "listener:stat:today:order_cnt:%s"

// 今日接单金额 key_prefix:20220804 zset XXX 接单金额
const RedisKeyListenerTodayOrderAmount = "listener:stat:today:order_amount:%s"

// 今日曝光推荐用户数 key_prefix:20220804 zset XXX 用户数
const RedisKeyListenerTodayRecommendUserCnt = "listener:stat:today:recommend_user_cnt:%s"

// 今日访问XXX个人资料用户数 key_prefix:20220804 zset XXX 用户数
const RedisKeyListenerTodayViewUserCnt = "listener:stat:today:view_user_cnt:%s"

// 今日访问XXX聊天页面用户数 key_prefix:20220804 zset XXX 用户数
const RedisKeyListenerTodayEnterChatUserCnt = "listener:stat:today:enter_chat_user_cnt:%s"

// 需要清理30天以上的今日数据
var NeedDelTodayRedisData = []string{RedisKeyListenerTodayOrderCnt, RedisKeyListenerTodayOrderAmount, RedisKeyListenerTodayRecommendUserCnt, RedisKeyListenerTodayViewUserCnt, RedisKeyListenerTodayEnterChatUserCnt}

// XXX首页统计数据看板 hash
const RedisKeyListenerDashboard = "listener:stat:dashboard:%d"

// 昨日XXX平均统计数据
const RedisKeyListenerLastDayAverage = "listener:last_day:stat:average"
