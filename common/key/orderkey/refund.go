package orderkey

// 退款原因标签
const RefundReasonTagCnt = 6

var RefundReasonTag = map[int]string{
	1: "不需要了/问题已解决",
	2: "下错订单/话题不匹配",
	3: "接单中/没人接听/太久不回复",
	4: "通讯问题",
	5: "问题没有得到解决/XXX态度不好",
	6: "其他",
}
