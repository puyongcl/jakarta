package orderkey

// 聊天服务订单状态
const (
	ChatOrderStateWaitPay1                    int64 = 1  // 待支付
	ChatOrderStateCancel2                     int64 = 2  // 已取消 订单结束
	ChatOrderStatePaySuccess3                 int64 = 3  // 已支付 语音订单待使用 文字订单自动变为使用中
	ChatOrderStateUsing4                      int64 = 4  // 服务中 已开始服务
	ChatOrderStateApplyRefund5                int64 = 5  // 申请退款
	ChatOrderStateAutoAgreeRefund6            int64 = 6  // 退款中 符合条件自动同意退款
	ChatOrderStateListenerAgreeRefund7        int64 = 7  // 退款中 XXX通过同意退款
	ChatOrderStateFinishRefund8               int64 = 8  // 退款完成 订单结束
	ChatOrderStateListenerRefuseRefund9       int64 = 9  // XXX拒绝退款
	ChatOrderStateSecondApplyRefund10         int64 = 10 // 再次申请退款 客服可以介入
	ChatOrderStateAdminAgreeRefund11          int64 = 11 // 退款中 客服介入同意退款
	ChatOrderStateUserStopService12           int64 = 12 // 服务结束 用户主动结束服务 可以退款
	ChatOrderStateUseOutWaitUserConfirm13     int64 = 13 // 服务结束 使用完 待用户确认 可以退款
	ChatOrderState5StartRatingAndFinish14     int64 = 14 // 订单完成 用户评价 满意并确认 订单结束
	ChatOrderStateNot5StarWaitConfirm15       int64 = 15 // 评价完成 用户不满意或者一般 可以退款 可以自动确认完成
	ChatOrderStateConfirmFinish16             int64 = 16 // 订单完成 用户手动确认 订单结束
	ChatOrderStateExpire17                    int64 = 17 // 服务过期 未使用或者未使用完 待用户确认 可以退款
	ChatOrderStateAutoCommentFinish18         int64 = 18 // 订单完成 订单由系统自动评价并确认完成 订单结束
	ChatOrderStateAutoConfirmFinish19         int64 = 19 // 订单完成 系统自动确认完成（不满意或一般评价）订单结束
	ChatOrderStateRefundRefuseAutoFinish20    int64 = 20 // 订单完成 系统自动确认完成（退款拒绝订单）订单结束
	ChatOrderStateSettle21                    int64 = 21 // 订单已经结算 订单完成 之后 自动进行
	ChatOrderStatePayFail22                   int64 = 22 // 支付失败
	ChatOrderStateRefundPayFail23             int64 = 23 // 退款失败
	ChatOrderStateTextOrderPaySuccess24       int64 = 24 // 文字订单 支付成功即时生效 可随时结束和评价
	ChatOrderStateAutoStartRefund25           int64 = 25 // 退款中 自动提交退款 请求支付方退款
	ChatOrderStateAdminStartRefund26          int64 = 26 // 退款中 管理员手动提交退款 请求支付方退款
	ChatOrderStateAutoAgreeNotProcessRefund27 int64 = 27 // 退款中 XXX超时未处理 自动同意退款
)

// 可以发起退款的状态
var CanApplyRefundOrderState = []int64{ChatOrderStateUserStopService12, ChatOrderStateUseOutWaitUserConfirm13,
	ChatOrderStateNot5StarWaitConfirm15, ChatOrderStateExpire17}

// 可以更新退款到账结果
var CanRefundPassOrderState = []int64{ChatOrderStateAutoStartRefund25, ChatOrderStateAdminStartRefund26}

// 可以结束服务
var CanStopOrderState = []int64{ChatOrderStatePaySuccess3, ChatOrderStateUsing4, ChatOrderStateTextOrderPaySuccess24}

// 可以开始向支付方发起退款
var CanStartRefundOrderState = []int64{ChatOrderStateListenerAgreeRefund7, ChatOrderStateAdminAgreeRefund11,
	ChatOrderStateAutoAgreeRefund6, ChatOrderStateAutoAgreeNotProcessRefund27}

// 可以评价
var CanCommentOrderState = []int64{ChatOrderStateUserStopService12, ChatOrderStateUseOutWaitUserConfirm13,
	ChatOrderStateExpire17, ChatOrderStateTextOrderPaySuccess24}

// 可以自动好评
var CanAutoGoodCommentAndFinishOrderState = []int64{ChatOrderStateUserStopService12, ChatOrderStateUseOutWaitUserConfirm13,
	ChatOrderStateExpire17}

// 可以自动确认完成
var CanAutoConfirmFinishOrderState = []int64{ChatOrderStateNot5StarWaitConfirm15, ChatOrderStateListenerRefuseRefund9,
	ChatOrderStateSecondApplyRefund10}

// 可以确认完成
var CanConfirmFinishOrderState = []int64{ChatOrderStateUserStopService12, ChatOrderStateUseOutWaitUserConfirm13,
	ChatOrderStateNot5StarWaitConfirm15, ChatOrderStateExpire17, ChatOrderStateListenerRefuseRefund9,
	ChatOrderStateSecondApplyRefund10}

// 退款拒绝
var RefuseRefundOrderState = []int64{ChatOrderStateListenerRefuseRefund9, ChatOrderStateSecondApplyRefund10}

// 没有使用完和未使用
var CanExpiryOrderSate = []int64{ChatOrderStatePaySuccess3, ChatOrderStateUsing4}

// 退款状态
var ChatOrderRefundState = []int64{ChatOrderStateApplyRefund5, ChatOrderStateFinishRefund8, ChatOrderStateRefundPayFail23,
	ChatOrderStateAutoStartRefund25, ChatOrderStateListenerAgreeRefund7, ChatOrderStateAutoAgreeRefund6,
	ChatOrderStateListenerRefuseRefund9, ChatOrderStateSecondApplyRefund10, ChatOrderStateAdminAgreeRefund11,
	ChatOrderStateAutoStartRefund25, ChatOrderStateAdminStartRefund26, ChatOrderStateAutoAgreeNotProcessRefund27}

// 可以重置免费聊天次数的订单状态
var CanResetTextChatFreeCntState = []int64{ChatOrderStateUserStopService12, ChatOrderStateUseOutWaitUserConfirm13, ChatOrderStatePaySuccess3, ChatOrderStateTextOrderPaySuccess24}

// 系统自动操作的operator uid
const DefaultSystemOperatorUid = 10000

// 需要反馈的状态
var NeedFeedbackOderState = []int64{ChatOrderStateApplyRefund5, ChatOrderStateAutoAgreeRefund6,
	ChatOrderStateListenerAgreeRefund7, ChatOrderStateFinishRefund8, ChatOrderStateListenerRefuseRefund9,
	ChatOrderStateSecondApplyRefund10, ChatOrderStateAdminAgreeRefund11, ChatOrderStateUserStopService12,
	ChatOrderStateUseOutWaitUserConfirm13, ChatOrderState5StartRatingAndFinish14, ChatOrderStateNot5StarWaitConfirm15,
	ChatOrderStateConfirmFinish16, ChatOrderStateAutoCommentFinish18, ChatOrderStateAutoConfirmFinish19,
	ChatOrderStateRefundRefuseAutoFinish20, ChatOrderStateSettle21, ChatOrderStateRefundPayFail23,
	ChatOrderStateTextOrderPaySuccess24, ChatOrderStateAutoStartRefund25, ChatOrderStateAdminStartRefund26,
	ChatOrderStateAutoAgreeNotProcessRefund27}

// 用户已经下过订单的有效状态
var UserNormalChatOrderState = []int64{ChatOrderStatePaySuccess3, ChatOrderStateUsing4, ChatOrderStateApplyRefund5,
	ChatOrderStateAutoAgreeRefund6, ChatOrderStateListenerAgreeRefund7, ChatOrderStateFinishRefund8,
	ChatOrderStateListenerRefuseRefund9, ChatOrderStateSecondApplyRefund10, ChatOrderStateAdminAgreeRefund11,
	ChatOrderStateUserStopService12, ChatOrderStateUseOutWaitUserConfirm13, ChatOrderState5StartRatingAndFinish14,
	ChatOrderStateNot5StarWaitConfirm15, ChatOrderStateConfirmFinish16, ChatOrderStateAutoCommentFinish18,
	ChatOrderStateAutoConfirmFinish19, ChatOrderStateRefundRefuseAutoFinish20, ChatOrderStateSettle21,
	ChatOrderStateRefundPayFail23, ChatOrderStateTextOrderPaySuccess24, ChatOrderStateAutoStartRefund25,
	ChatOrderStateAdminStartRefund26, ChatOrderStateAutoAgreeNotProcessRefund27}

// 不显示 不统计的订单状态
var AbnormalOrderState = []int64{ChatOrderStateWaitPay1, ChatOrderStateCancel2}

// 订单完成的状态
var ConfirmOrderState = []int64{ChatOrderState5StartRatingAndFinish14, ChatOrderStateConfirmFinish16, ChatOrderStateAutoCommentFinish18, ChatOrderStateAutoConfirmFinish19, ChatOrderStateRefundRefuseAutoFinish20, ChatOrderStateSettle21}
