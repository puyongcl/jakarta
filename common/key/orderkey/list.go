package orderkey

// 订单列表类型
const (
	OrderListTypeAll             = 0 // 全部
	OrderListTypeNeedProcess     = 1 // 未完成 需要服务 同意退款 需要填写反馈
	OrderListTypeNeedFeedback    = 2 // 需要填写反馈
	OrderListTypeNeedAgreeRefund = 3 // 需要同意退款
)
