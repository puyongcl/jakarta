package paykey

//支付方式
const ThirdPaymentPayModelWechatPay = "WECHAT_PAY" //微信支付

//支付状态
const (
	ThirdPaymentPayTradeStateWait                     int64 = 1  //待支付
	ThirdPaymentPayTradeStateUserPaying               int64 = 2  // 用户支付中
	ThirdPaymentPayTradeStateFail                     int64 = 3  //支付失败
	ThirdPaymentPayTradeStateNotPay                   int64 = 4  // 未支付
	ThirdPaymentPayTradeStateClosed                   int64 = 5  // 已关闭
	ThirdPaymentPayTradeStateRevoked                  int64 = 6  // 已撤销
	ThirdPaymentPayTradeStatePayError                 int64 = 7  // 支付错误
	ThirdPaymentPayTradeStateSuccess                  int64 = 8  //支付成功
	ThirdPaymentPayTradeStateSuccessAmountError       int64 = 9  //支付成功金额不正确
	ThirdPaymentPayTradeStateStartRefund              int64 = 10 //发起退款
	ThirdPaymentPayTradeStateRefundProcessing         int64 = 11 //退款处理中
	ThirdPaymentPayTradeStateRefundSuccess            int64 = 12 //退款成功
	ThirdPaymentPayTradeStateRefundSuccessAmountError int64 = 13 //退款成功部分退款
	ThirdPaymentPayTradeStateRefundClosed             int64 = 14 //退款关闭
	ThirdPaymentPayTradeStateRefundFail               int64 = 15 //退款失败
)

// 可以更新支付成功失败的状态
var CanUpdatePayResultState = []int64{ThirdPaymentPayTradeStateWait, ThirdPaymentPayTradeStateUserPaying}

// 可以更新退款的状态
var CanUpdateRefundState = []int64{ThirdPaymentPayTradeStateSuccess, ThirdPaymentPayTradeStateSuccessAmountError}

// 不能再次支付的状态
var NotNeedPayState = []int64{ThirdPaymentPayTradeStateWait, ThirdPaymentPayTradeStateUserPaying, ThirdPaymentPayTradeStateSuccess, ThirdPaymentPayTradeStateRefundSuccess}

// 不能再次退款的状态
var NotNeedRefundState = []int64{ThirdPaymentPayTradeStateStartRefund, ThirdPaymentPayTradeStateRefundSuccess}

// 可以更新退款结果的
var CanUpdateRefundResult = []int64{ThirdPaymentPayTradeStateStartRefund, ThirdPaymentPayTradeStateRefundProcessing}

// 支付过期时间 分钟
const PaymentExpireTimeMinute = 3
