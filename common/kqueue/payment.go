//KqMessage
package kqueue

//微信支付回调更改支付状态通知
type UpdatePaymentStatusMessage struct {
	FlowNo         string `json:"flowNo"`
	TradeState     string `json:"tradeState"`
	TransactionId  string `json:"transactionId"`
	TradeType      string `json:"tradeType"`
	TradeStateDesc string `json:"tradeStateDesc"`
	BankType       string `json:"bankType"`
	PayTime        int64  `json:"payTime"`
	PayAmount      int64  `json:"payAmount"`
}

// 微信支付退款回调
type UpdateRefundStatusMessage struct {
	FlowNo          string `json:"flowNo"`
	PayFlowNo       string `json:"payFlowNo"`
	TransactionId   string `json:"transactionId"`
	ReceivedAccount string `json:"receivedAccount"`
	SuccessTime     int64  `json:"successTime"`
	Status          string `json:"Status"`
	Amount          int64  `json:"amount"`
}

// 提现回调
type UpdateCashStatusMessage struct {
	Type            int64    `json:"type"`             // 1 支付结果回调 2 任务作废回调
	WorkNumber      string   `json:"work_number"`      // 任务编号
	CompanyId       int64    `json:"company_id"`       //
	Status          string   `json:"status"`           // -1 作废
	TovoidTime      string   `json:"tovoid_time"`      // 作废时间
	CustomerNumbers []string `json:"customer_numbers"` // 此任务包含的所有自定义任务单号
	UserId          string   `json:"user_id"`          // 用户id
	Number          string   `json:"number"`           // 打款流水号
	PayStatus       string   `json:"pay_status"`       // 打款状态 1 未结算 2 待结算 3 结算中 4 已结算 5 结算失败
	CustomNumber    string   `json:"custom_number"`    // 自定义流水号
	Msg             string   `json:"msg"`              // 支付失败原因
	PayTime         string   `json:"pay_time"`         // 支付时间 格式2020-10-10 12:00:00
}

// 提交转账
type CommitMoveCashMessage struct {
	FlowNo string `json:"flowNo"`
	Uid    int64  `json:"uid"`
}
