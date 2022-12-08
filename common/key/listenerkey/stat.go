package listenerkey

// 自动刷新XXX推荐、访问个人资料、进入聊天页面的用户数 频率 分钟
const AutoUpdateListenerUserStatIntervalMinute = 2

// 每次自动刷新几天内有更新过的XXX
const AutoUpdateListenerUserStatRangeDay = 7

// 自动清理多少天以上的历史统计数据
const AutoClearListenerTodayStatHistoryData = 60

// 自动清理多少天以上的历史新用户推荐XXX数据
const AutoClearNewUserRecommendListenerData = 60
