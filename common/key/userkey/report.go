package userkey

var ReportTag = []int64{1, 2, 3, 4, 5, 6}

var ReportTagText = map[int]string{
	1: "接单时不专心",
	2: "恶意骚扰",
	3: "色情/性骚扰",
	4: "涉及政治",
	5: "诈骗",
	6: "其他",
}

// 处理举报用户
const (
	ReportStateCreated         = 1 // 用户上报
	ReportStateBlock2Days      = 2 // 封号2天
	ReportStateBlock7Days      = 3 // 封号7天
	ReportStateBlockForever    = 4 // 封号永久
	ReportStateNotTrue         = 5 // 不符合事实
	ReportStateCancel          = 6 // 取消封号
	ReportStateAdminBan2Days   = 7 // 管理员直接封号2天 无需用户上报
	ReportStateAdminBan7Days   = 8 // 管理员直接封号7天 无需用户上报
	ReportStateAdminBanForever = 9 // 管理员直接封号永久 无需用户上报
)

var AdminBanUserAction = []int64{ReportStateAdminBan2Days, ReportStateAdminBan7Days, ReportStateAdminBanForever}

// XXX上报需要XX援助的用户
const (
	ReportNeedHelpUserStateCreated = 1 // 上报
	ReportNeedHelpUserStateMark    = 2 // 已经处理
	ReportNeedHelpUserStateMarkNot = 3 // 无法处理
)
