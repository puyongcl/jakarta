package listenerkey

// XXX钱包流水类型
const (
	ListenerSettleTypeConfirm       int64 = 1 // 已经确认的收益
	ListenerSettleTypeOrderAmount   int64 = 2 // 本月订单总金额（尚未分成、可能退款）
	ListenerSettleTypeAlreadyRefund int64 = 3 // 已经退款
	ListenerSettleTypeApplyCash     int64 = 4 // 申请提现
	ListenerSettleTypeStartCash     int64 = 5 // 开始提现
	ListenerSettleTypeCashSuccess   int64 = 6 // 提现成功
	ListenerSettleTypeCashFail      int64 = 7 // 提现失败
	ListenerSettleTypeCashCancel    int64 = 8 // 提现取消
	ListenerSettleTypeFix           int64 = 9 // 维护金额 加 或 减
)

// 提现相关状态
var CashSettleType = []int64{ListenerSettleTypeApplyCash, ListenerSettleTypeStartCash, ListenerSettleTypeCashSuccess, ListenerSettleTypeCashFail}

// 收益相关状态
var IncomeSettleType = []int64{ListenerSettleTypeConfirm}

// 提现日期 每月11-17日
const (
	ListenerMoveCashStatDay int = 11
	ListenerMoveCashEndDay  int = 17
)

// 提现金额最小值
const ListenerMoveCashMinAmount int64 = 10000
