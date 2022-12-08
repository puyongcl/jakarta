package uniqueid

import (
	"fmt"
	"jakarta/common/tool"
	"strings"
	"time"
)

//生成sn单号
const (
	SnPrefixTextChatOrderId    = "ct" //聊天服务文字订单前缀 jakarta_order/chat_order
	SnPrefixVoiceChatOrderId   = "cv" //聊天服务语音订单前缀 jakarta_order/chat_order
	SnPrefixThirdPaymentFlowNo = "cp" //第三方支付流水记录前缀 jakarta_payment/third_payment_flow
	SnPrefixThirdRefundFlowNo  = "cr" //第三方支付退款流水记录前缀 jakarta_payment/third_refund_flow
	SnPrefixThirdCashFlowNo    = "cc" //第三方支付退款流水记录前缀 jakarta_payment/third_cash_flow
	SnPrefixWalletFlowNo       = "cw" //钱包流水 jakarta_payment/listener_wallet_log
)

//生成唯一单号
func GenSn(snPrefix string) string {
	return strings.ToUpper(fmt.Sprintf("%s%s", snPrefix, fmt.Sprintf("%x%x", time.Now().UnixNano(), GenId())))
}

// 有重复
func genSn2(snPrefix string) string {
	return fmt.Sprintf("%s%s%s", snPrefix, time.Now().Format("20060102150405"), tool.Krand(8, tool.KC_RAND_KIND_NUM))
}
