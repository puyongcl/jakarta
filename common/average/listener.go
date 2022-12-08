package average

type ListenerStatAverage struct {
	CreateTime                         string `json:"create_time"`
	Id                                 string `json:"id"`
	YesterdayOrderCnt                  int64  `json:"yesterday_order_cnt"`
	YesterdayOrderCntRank              int64  `json:"yesterday_order_cnt_rank"`
	YesterdayOrderAmount               int64  `json:"yesterday_order_amount"`
	YesterdayOrderAmountRank           int64  `json:"yesterday_order_amount_rank"`
	YesterdayRecommendUserCnt          int64  `json:"yesterday_recommend_user_cnt"`
	YesterdayRecommendUserCntRank      int64  `json:"yesterday_recommend_user_cnt_rank"`
	YesterdayEnterChatUserCnt          int64  `json:"yesterday_enter_chat_user_cnt"`
	YesterdayEnterChatUserCntRank      int64  `json:"yesterday_enter_chat_user_cnt_rank"`
	YesterdayViewUserCnt               int64  `json:"yesterday_view_user_cnt"`
	YesterdayViewUserCntRank           int64  `json:"yesterday_view_user_cnt_rank"`
	Last30DaysPaidUserCnt              int64  `json:"last_30_days_paid_user_cnt"`
	Last30DaysEnterChatUserCnt         int64  `json:"last_30_days_enter_chat_user_cnt"`
	Last30DaysRepeatPaidUserCnt        int64  `json:"last_30_days_repeat_paid_user_cnt"`
	Last30DaysAveragePaidAmountPerUser int64  `json:"last_30_days_average_paid_amount_per_user"`
	Last30DaysAveragePaidAmountPerDay  int64  `json:"last_30_days_average_paid_amount_per_day"`
	Last7DaysPaidUserCnt               int64  `json:"last_7_days_paid_user_cnt"`
	Last7DaysRepeatPaidUserCnt         int64  `json:"last_7_days_repeat_paid_user_cnt"`
	Last7DaysAveragePaidAmountPerUser  int64  `json:"last_7_days_average_paid_amount_per_user"`
	Last7DaysAveragePaidAmountPerDay   int64  `json:"last_7_days_average_paid_amount_per_day"`
	Last30DaysPaidUserRate             int64  `json:"last_30_days_paid_user_rate"`        // 近30天内下单人数占进入聊天页面人数比例(万分)
	Last30DaysRepeatPaidUserRate       int64  `json:"last_30_days_repeat_paid_user_rate"` // 近30天内复购人数占下单人数比例（万分）
	Last7DaysRepeatPaidUserRate        int64  `json:"last_7_days_repeat_paid_user_rate"`  // 近7天内复购人数占下单人数比例
	Last7DaysPaidUserRate              int64  `json:"last_7_days_paid_user_rate"`         // 近7天内下单人数占进入聊天界面用户数比例（万分）
	Last7DaysEnterChatUserCnt          int64  `json:"last_7_days_enter_chat_user_cnt"`    // 最近7天进入聊天界面的人数
	Last30DaysPaidAmountSum            int64  `json:"last_30_days_paid_amount_sum"`       // 近30天下单总金额
	Last7DaysPaidAmountSum             int64  `json:"last_7_days_paid_amount_sum"`        // 近7天下单总金额
}
