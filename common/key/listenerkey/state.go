package listenerkey

// 工作状态
const (
	ListenerWorkStateWorking        int64 = 2 // 手动可接单
	ListenerWorkStateRestingAuto    int64 = 3 // 系统自动转休息中
	ListenerWorkStateRestingManual  int64 = 4 // 手动转休息中
	ListenerWorkStateAccountDeleted int64 = 9 // 账户已经注销
)

// 休息中时间段默认设置
const (
	ListenerRestingStartTime = "07-00"
	ListenerRestingStopTime  = "01-00"
)

// 休息开关
const ListenerRestingSwitchEnable = 1  // 打开
const ListenerRestingSwitchDisable = 2 // 关闭

// 多少分钟内不回复消息 自动转休息中 分钟
const NotReplyAutoSwitchWorkStateIntervalMinute = 4

// 休息中状态
var ListenerRestState = []int64{ListenerWorkStateRestingAuto, ListenerWorkStateRestingManual}
