package wxpay

import (
	"jakarta/common/key/paykey"
)

// 交易状态，枚举值：
//SUCCESS：支付成功
//REFUND：转入退款
//NOTPAY：未支付
//CLOSED：已关闭
//REVOKED：已撤销（付款码支付）
//USERPAYING：用户支付中（付款码支付）
//PAYERROR：支付失败(其他原因，如银行返回失败)
const (
	SUCCESS    = "SUCCESS"    //支付成功
	REFUND     = "REFUND"     //转入退款
	NOTPAY     = "NOTPAY"     //未支付
	CLOSED     = "CLOSED"     //已关闭
	REVOKED    = "REVOKED"    //已撤销（付款码支付）
	USERPAYING = "USERPAYING" //用户支付中（付款码支付）
	PAYERROR   = "PAYERROR"   //支付失败(其他原因，如银行返回失败)
	ABNORMAL   = "ABNORMAL"   // 退款异常，退款
	PROCESSING = "PROCESSING" // 退款处理中
)

func GetPayStatusByWXPayTradeState(wxPayTradeState string) int64 {
	switch wxPayTradeState {
	case SUCCESS: // 支付成功
		return paykey.ThirdPaymentPayTradeStateSuccess
	case USERPAYING: // 支付中
		return paykey.ThirdPaymentPayTradeStateUserPaying
	case NOTPAY: // 未支付
		return paykey.ThirdPaymentPayTradeStateNotPay
	case CLOSED: // 已关闭
		return paykey.ThirdPaymentPayTradeStateClosed
	case REVOKED: // 已撤销
		return paykey.ThirdPaymentPayTradeStateRevoked
	case PAYERROR:
		return paykey.ThirdPaymentPayTradeStatePayError
	case REFUND: // 已退款
		return paykey.ThirdPaymentPayTradeStateRefundSuccess
	default:
		return 0
	}
}

func GetRefundStatusByWXPayTradeState(wxPayTradeState string) int64 {
	switch wxPayTradeState {
	case SUCCESS: //退款成功
		return paykey.ThirdPaymentPayTradeStateRefundSuccess
	case CLOSED: //退款关闭
		return paykey.ThirdPaymentPayTradeStateRefundClosed
	case PROCESSING:
		return paykey.ThirdPaymentPayTradeStateRefundProcessing
	case ABNORMAL:
		return paykey.ThirdPaymentPayTradeStateRefundFail
	default:
		return paykey.ThirdPaymentPayTradeStateFail
	}
}
