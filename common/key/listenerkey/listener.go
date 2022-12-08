package listenerkey

// 资质证明类型
const (
	CertTypeProfessional = 2 // 有专业证书
	CertTypeExperienced  = 4 // 有其他平台从业经历
	CertTypeBeginner     = 6 // 新手
)

// 教育水平
const (
	EducationStudent      = 1 // 在校学生
	EducationUnderCollege = 2 // 大专以下
	EducationCollege      = 3 // 大专
	EducationBachelor     = 4 // 本科
	EducationMaster       = 5 // 硕士
	EducationDoctor       = 6 // 博士
)

// 婚姻状况
const (
	MaritalStatusSingle            = 1 // 单身
	MaritalStatusWithPartner       = 2 // 有对象
	MaritalStatusMarried           = 3 // 已婚
	MaritalStatusMarriedWithChild  = 4 // 已婚有小孩
	MaritalStatusDivorced          = 5 // 离异
	MaritalStatusDivorcedWithChild = 6 // 离异有小孩
)

// 年龄
const (
	AgeRange1 = 1 // 60后
	AgeRange2 = 2 // 70后
	AgeRange3 = 3 // 80后
	AgeRange4 = 4 // 90后
)

// 排序字段
const (
	ListenerSortOrderDefault        = 1 // 综合排序
	ListenerSortOrderRatingStar     = 2 // 服务满意率
	ListenerSortOrderRepeatCustomer = 3 // 回头客人数（可以换成回头率）
	ListenerSortOrderChatMinute     = 4 // 服务时长
)

// 性别
const (
	GenderMale   = 1 // 男性
	GenderFemale = 2 // 女性
)

// 更新字符串时 需要清空用NULL
const (
	EmptyString = "NULL"
	EmptyInt    = -1
)

// 需要审核的字段
var ListenerProfileCheckFieldArray = []string{
	"nickName",     // 昵称
	"avatar",       // 头像 大小头像
	"region",       // 地区 省-市
	"gender",       // 性别
	"birthday",     // 出生日期
	"id",           // 身份证和照片
	"specialties",  // 专业领域
	"introduction", // 个人介绍
	"voiceFile",    // 声音文件
	"experience",   // 2段经验
	"cert",         // 5个证明 和 类型
	"autoReply",    // 3个自动回复
	"chatPrice",    // 2个价格
}

var ListenerProfileCheckFieldText = map[string]string{
	"nickName":             "昵称",
	"listenerName":         "姓名",
	"avatar":               "头像",
	"province":             "地区",
	"city":                 "城市",
	"gender":               "性别",
	"birthday":             "生日",
	"id":                   "身份信息",
	"specialties":          "擅长领域",
	"introduction":         "个人介绍",
	"voiceFile":            "声音签名",
	"experience1":          "经历1",
	"experience2":          "经历2",
	"certType":             "资质类型",
	"otherPlatformAccount": "其他平台账号",
	"certFiles1":           "资质附件1",
	"certFiles2":           "资质附件2",
	"certFiles3":           "资质附件3",
	"certFiles4":           "资质附件4",
	"certFiles5":           "资质附件5",
	"autoReplyNew":         "自动回复1",
	"autoReplyProcessing":  "自动回复2",
	"autoReplyFinish":      "自动回复3",
	"textChatPrice":        "文字服务价格",
	"voiceChatPrice":       "通话服务价格",
}
