package listenerkey

const Rating5Star = 5 // 满意
const Rating3Star = 3 // 一般
const Rating1Star = 1 // 不满意

func GetRatingText(star int64) string {
	switch star {
	case Rating5Star:
		return "满意"
	case Rating3Star:
		return "良好"
	case Rating1Star:
		return "不满意"
	default:
		return "未知"
	}
}

// 展示最高几个评价标签
const ShowTopCommentTagCnt = 4

const DefaultComment = "用户未填写评价内容"

const DefaultFeedback = "XXX未填写反馈内容"
