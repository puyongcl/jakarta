package hfbfcash

const (
	CashStatusCreated       int64 = 1 // 已请求第三方
	CashStatusNotSettle     int64 = 2 // 未结算
	CashStatusWaitSettle    int64 = 3 // 待结算
	CashStatusStartSettle   int64 = 4 // 结算中
	CashStatusSettleSuccess int64 = 5 // 已结算
	CashStatusSettleFail    int64 = 6 // 结算失败
	CashStatusDrop          int64 = 7 // 作废
)

// 回调类型
const (
	CallbackTypePay  int64 = 1 // 支付结果回调
	CallbackTypeDrop int64 = 2 // 作废回调
)

// 获取打款状态
func GetCashPayStatus(payStatus string) int64 {
	switch payStatus {
	case "1": // 未结算
		return CashStatusNotSettle
	case "2": // 待结算
		return CashStatusWaitSettle
	case "3": // 结算中
		return CashStatusStartSettle
	case "4": // 已结算（打款成功）
		return CashStatusSettleSuccess
	case "5": // 结算失败
		return CashStatusSettleFail

	default:

	}
	return 0
}

// 获取查询任务接口返回的数据状态 不完全对应
func GetDataStatus(st int64) int64 {
	switch st {
	case -1: // 作废
		return CashStatusDrop
	case 1: // 待提交
		return CashStatusNotSettle
	case 2: // 待完成
		return CashStatusWaitSettle
	case 3: // 已完成
		return 0
	default:

	}
	return 0
}
