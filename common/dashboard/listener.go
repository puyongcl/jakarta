package dashboard

import "strconv"

// XXX首页数据看板 数据存在hash 中 field必须和 下方结构体 名字保持一致
type ListenerDashboardRedisHashData struct {
	ListenerUid                        string `json:"listenerUid"`
	TodayOrderCnt                      string `json:"todayOrderCnt"`                      // 今日接单数
	TodayOrderCntRank                  string `json:"todayOrderCntRank"`                  // 今日接单数排名
	TodayOrderAmount                   string `json:"todayOrderAmount"`                   // 今日接单金额
	TodayOrderAmountRank               string `json:"todayOrderAmountRank"`               // 今日接单金额排名
	TodayRecommendUserCnt              string `json:"todayRecommendUserCnt"`              // 今日推荐用户数
	TodayRecommendUserCntRank          string `json:"todayRecommendUserCntRank"`          // 今日推荐用户数排名
	TodayEnterChatUserCnt              string `json:"todayEnterChatUserCnt"`              // 今日进入聊天页面用户数
	TodayEnterChatUserCntRank          string `json:"todayEnterChatUserCntRank"`          // 今日进入聊天界面用户数排名
	TodayViewUserCnt                   string `json:"todayViewUserCnt"`                   // 今日访问资料页面用户数
	TodayViewUserCntRank               string `json:"todayViewUserCntRank"`               // 今日资料页面用户数排名
	TodayStatUpdateTime                string `json:"todayStatUpdateTime"`                // 今日统计数据更新时间
	Last30DaysPaidUserCnt              string `json:"last30DaysPaidUserCnt"`              // 过去30天下单人数（下单人数占进入聊天页面的人数比例）
	Last30DaysEnterChatUserCnt         string `json:"last30DaysEnterChatUserCnt"`         // 过去30天进入聊天页面人数
	Last7DaysEnterChatUserCnt          string `json:"last7DaysEnterChatUserCnt"`          // 过去7天进入聊天页面人数
	Last7DaysRepeatPaidUserCnt         string `json:"last7DaysRepeatPaidUserCnt"`         // 过去7天复购人数 （复购人数占下单人数比例）
	Last7DaysPaidUserCnt               string `json:"last7DaysPaidUserCnt"`               // 过去7天下单人数（
	Last30DaysRepeatPaidUserCnt        string `json:"last30DaysRepeatPaidUserCnt"`        // 过去30天复购人数
	Last7DaysAveragePaidAmountPerUser  string `json:"last7DaysAveragePaidAmountPerUser"`  // 过去7天人均消费
	Last30DaysAveragePaidAmountPerUser string `json:"last30DaysAveragePaidAmountPerUser"` // 过去30天人均消费
	Last7DaysAveragePaidAmountPerDay   string `json:"last7DaysAveragePaidAmountPerDay"`   // 过去7天日均消费
	Last30DaysAveragePaidAmountPerDay  string `json:"last30DaysAveragePaidAmountPerDay"`  // 过去30天日均消费
	Last30DaysPaidAmountSum            string `json:"last30DaysPaidAmountSum"`            // 近30天累计下单金额
	Last7DaysPaidAmountSum             string `json:"last7DaysPaidAmountSum"`             // 近7天累计下单金额
	LastDayStatUpdateTime              string `json:"lastDayStatUpdateTime"`              // 过去几天统计数据更新时间
	OneStarRatingOrderCnt              string `json:"oneStarRatingOrderCnt"`              // 累计不满意评价订单数
	RefundOrderCnt                     string `json:"refundOrderCnt"`                     // 累计退款订单数
}

type ListenerDashboard struct {
	ListenerUid                        int64  `json:"listenerUid"`
	TodayOrderCnt                      int64  `json:"todayOrderCnt"`                      // 今日接单数
	TodayOrderCntRank                  int64  `json:"todayOrderCntRank"`                  // 今日接单数排名
	TodayOrderAmount                   int64  `json:"todayOrderAmount"`                   // 今日接单金额
	TodayOrderAmountRank               int64  `json:"todayOrderAmountRank"`               // 今日接单金额排名
	TodayRecommendUserCnt              int64  `json:"todayRecommendUserCnt"`              // 今日推荐用户数
	TodayRecommendUserCntRank          int64  `json:"todayRecommendUserCntRank"`          // 今日推荐用户数排名
	TodayEnterChatUserCnt              int64  `json:"todayEnterChatUserCnt"`              // 今日进入聊天页面用户数
	TodayEnterChatUserCntRank          int64  `json:"todayEnterChatUserCntRank"`          // 今日进入聊天界面用户数排名
	TodayViewUserCnt                   int64  `json:"todayViewUserCnt"`                   // 今日访问资料页面用户数
	TodayViewUserCntRank               int64  `json:"todayViewUserCntRank"`               // 今日资料页面用户数排名
	TodayStatUpdateTime                string `json:"todayStatUpdateTime"`                // 今日统计数据更新时间
	Last30DaysPaidUserCnt              int64  `json:"last30DaysPaidUserCnt"`              // 过去30天下单人数（下单人数占进入聊天页面的人数比例）
	Last30DaysEnterChatUserCnt         int64  `json:"last30DaysEnterChatUserCnt"`         // 过去30天进入聊天页面人数
	Last7DaysEnterChatUserCnt          int64  `json:"last7DaysEnterChatUserCnt"`          // 过去30天进入聊天页面人数
	Last7DaysRepeatPaidUserCnt         int64  `json:"last7DaysRepeatPaidUserCnt"`         // 过去7天复购人数 （复购人数占下单人数比例）
	Last7DaysPaidUserCnt               int64  `json:"last7DaysPaidUserCnt"`               // 过去7天下单人数（
	Last30DaysRepeatPaidUserCnt        int64  `json:"last30DaysRepeatPaidUserCnt"`        // 过去30天复购人数
	Last7DaysAveragePaidAmountPerUser  int64  `json:"last7DaysAveragePaidAmountPerUser"`  // 过去7天人均消费
	Last30DaysAveragePaidAmountPerUser int64  `json:"last30DaysAveragePaidAmountPerUser"` // 过去30天人均消费
	Last7DaysAveragePaidAmountPerDay   int64  `json:"last7DaysAveragePaidAmountPerDay"`   // 过去7天日均消费
	Last30DaysAveragePaidAmountPerDay  int64  `json:"last30DaysAveragePaidAmountPerDay"`  // 过去30天日均消费
	Last30DaysPaidAmountSum            int64  `json:"last30DaysPaidAmountSum"`            // 近30天累计下单金额
	Last7DaysPaidAmountSum             int64  `json:"last7DaysPaidAmountSum"`             // 近7天累计下单金额
	LastDayStatUpdateTime              string `json:"lastDayStatUpdateTime"`              // 过去几天统计数据更新时间
	OneStarRatingOrderCnt              int64  `json:"oneStarRatingOrderCnt"`              // 累计不满意评价订单数
	RefundOrderCnt                     int64  `json:"refundOrderCnt"`                     // 累计退款订单数
}

func TransferListenerDashboardData(in *ListenerDashboardRedisHashData, out *ListenerDashboard) (err error) {
	var value uint64
	// 字符串转成整数
	value, err = strconv.ParseUint(in.ListenerUid, 10, 64)
	if err != nil {
		return err
	}
	out.ListenerUid = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayOrderCnt, 10, 64)
	if err != nil {
		return err
	}
	out.TodayOrderCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayOrderCntRank, 10, 64)
	if err != nil {
		return err
	}
	out.TodayOrderCntRank = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayOrderAmount, 10, 64)
	if err != nil {
		return err
	}
	out.TodayOrderAmount = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayOrderAmountRank, 10, 64)
	if err != nil {
		return err
	}
	out.TodayOrderAmountRank = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayRecommendUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.TodayRecommendUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayRecommendUserCntRank, 10, 64)
	if err != nil {
		return err
	}
	out.TodayRecommendUserCntRank = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayEnterChatUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.TodayEnterChatUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayEnterChatUserCntRank, 10, 64)
	if err != nil {
		return err
	}
	out.TodayEnterChatUserCntRank = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayViewUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.TodayViewUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.TodayViewUserCntRank, 10, 64)
	if err != nil {
		return err
	}
	out.TodayViewUserCntRank = int64(value)
	//
	out.TodayStatUpdateTime = in.TodayStatUpdateTime
	//
	value, err = strconv.ParseUint(in.Last30DaysPaidUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.Last30DaysPaidUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.Last30DaysEnterChatUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.Last30DaysEnterChatUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.Last7DaysEnterChatUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.Last7DaysEnterChatUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.Last7DaysRepeatPaidUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.Last7DaysRepeatPaidUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.Last7DaysPaidUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.Last7DaysPaidUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.Last30DaysRepeatPaidUserCnt, 10, 64)
	if err != nil {
		return err
	}
	out.Last30DaysRepeatPaidUserCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.Last7DaysAveragePaidAmountPerUser, 10, 64)
	if err != nil {
		return err
	}
	out.Last7DaysAveragePaidAmountPerUser = int64(value)
	//
	value, err = strconv.ParseUint(in.Last30DaysAveragePaidAmountPerUser, 10, 64)
	if err != nil {
		return err
	}
	out.Last30DaysAveragePaidAmountPerUser = int64(value)
	//
	value, err = strconv.ParseUint(in.Last7DaysAveragePaidAmountPerDay, 10, 64)
	if err != nil {
		return err
	}
	out.Last7DaysAveragePaidAmountPerDay = int64(value)
	//
	value, err = strconv.ParseUint(in.Last30DaysAveragePaidAmountPerDay, 10, 64)
	if err != nil {
		return err
	}
	out.Last30DaysAveragePaidAmountPerDay = int64(value)
	//
	value, err = strconv.ParseUint(in.Last30DaysPaidAmountSum, 10, 64)
	if err != nil {
		return err
	}
	out.Last30DaysPaidAmountSum = int64(value)
	//
	value, err = strconv.ParseUint(in.Last7DaysPaidAmountSum, 10, 64)
	if err != nil {
		return err
	}
	out.Last7DaysPaidAmountSum = int64(value)
	//
	out.LastDayStatUpdateTime = in.LastDayStatUpdateTime
	//
	value, err = strconv.ParseUint(in.OneStarRatingOrderCnt, 10, 64)
	if err != nil {
		return err
	}
	out.OneStarRatingOrderCnt = int64(value)
	//
	value, err = strconv.ParseUint(in.RefundOrderCnt, 10, 64)
	if err != nil {
		return err
	}
	out.RefundOrderCnt = int64(value)
	return nil
}

// 统计最近30天的数据
const LastStat30Days = 30

// 统计最近7天数据
const LastStat7Days = 7

// 比例
const DivideNumber = 10000
