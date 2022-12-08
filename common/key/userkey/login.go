package userkey

// 登陆方式
const (
	UserAuthTypePasswd = "passwd" //用户名密码登陆
	UserAuthTypeWXMini = "wxMini" //微信小程序登陆
	UserAuthTypeSMS    = "sms"    //手机短信登陆
	UserAuthTypeWXApp  = "wxApp"  //APP内微信登陆
)

// 判断登陆方式参数是否合法
var AuthType = []string{UserAuthTypePasswd, UserAuthTypeWXMini}

// 账户状态
const (
	AccountStateNormal = 2 // 2 正常
	AccountStateBan    = 4 // 6 封禁
	AccountStateCancel = 8 // 8 注销
)

// 删除的账号 在authkey 前 加上 DEL#时间戳# 以便识别账户之间的关系
const DeleteAccountAuthKeyPrefix = "DEL#%d#"

// 账户角色
const (
	UserTypeNormalUser      = 2 // 2 普通用户
	UserTypeListener        = 4 // 4 XXX
	UserTypeAdmin           = 6 // 6 管理员
	UserTypeCustomerService = 7 // 7 客服
	UserTypeNotify          = 8 // 8 通知消息发送（客服）
)

var AllowLoginUserType = []int64{UserTypeNormalUser, UserTypeListener, UserTypeAdmin}
