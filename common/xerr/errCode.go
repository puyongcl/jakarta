package xerr

//成功返回
const OK uint32 = 200

/**(前3位代表业务,后三位代表具体功能)**/

//全局错误码
const (
	ServerCommonError         uint32 = 1001
	RequestParamError         uint32 = 1002
	TokenExpireError          uint32 = 1003
	TokenGenerateError        uint32 = 1004
	DbError                   uint32 = 1005
	DbUpdateAffectedZeroError uint32 = 1006
	ThirdPartRequestError     uint32 = 1007
	RedisLockFail             uint32 = 1008
	PaymentFail               uint32 = 1009
	NotAllowEmptyParam        uint32 = 1100
)

//用户模块
const (
	UserNotExist          uint32 = 200001
	WXminiAuthFail        uint32 = 200002
	GetUserSignFail       uint32 = 200003
	UserTypeNotAllowLogin uint32 = 200004
	UserRegInfoError      uint32 = 200005
	UserUidNotMatch       uint32 = 200006
)

// XXX
const (
	ListenerErrorMoveCash           uint32 = 300001 // 提现相关错误
	ListenerErrorNotSetBankCard     uint32 = 300002 // 未设置结算方式
	ListenerErrorProfile            uint32 = 300003 // 资料不符合规则
	ListenerErrorStartMoveCashError uint32 = 300004 // 请求第三方提现错误
	ListenerErrorEditWords          uint32 = 300005 // 编辑XXX常用语错误
)

// 订单
const (
	OrderError                      uint32 = 400001
	OrderErrorAlreadyComment        uint32 = 400002 // 已经评价
	OrderErrorAlreadyFeedback       uint32 = 400003 // 已经反馈
	OrderErrorFeedbackNotAllowEmpty uint32 = 400004 // 反馈不能为空
	OrderErrorNotAllowStopOrder     uint32 = 400005 // 当前不能结束订单
)

// 聊天
const (
	ChatError             uint32 = 500001
	ChatErrorListenerBusy uint32 = 500002 // 对方占线或通话中
)

// 管理后台
const (
	AdminGenContractError = 600001 // 生成合同错误
)

// 支付
const (
	PaymentErrorFlowAlreadyExist = 700001 // 流水已经存在
)

// 发布XX
const (
	BbsErrorAlreadyReply     = 800001 // 已经回复
	BbsErrorFreq             = 800002 // 操作频繁
	BbsErrorStoryNotFound    = 800003 // XX或者回复查不到
	BbsErrorStoryVerifyError = 800004 // XX校验错误
)
