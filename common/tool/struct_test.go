package tool

import (
	"fmt"
	"testing"
)

func TestGetStructFieldName(t *testing.T) {
	// XXX首页数据看板
	type listenerDashboard struct {
		ListenerUid                        int64 `json:"listenerUid"`
		TodayOrderCnt                      int64 `json:"todayOrderCnt"`                      // 今日接单数
		TodayOrderCntRank                  int64 `json:"todayOrderCntRank"`                  // 今日接单数排名
		TodayOrderAmount                   int64 `json:"todayOrderAmount"`                   // 今日接单金额
		TodayOrderAmountRank               int64 `json:"todayOrderAmountRank"`               // 今日接单金额排名
		TodayRecommendUserCnt              int64 `json:"todayRecommendUserCnt"`              // 今日推荐用户数
		TodayRecommendUserCntRank          int64 `json:"todayRecommendUserCntRank"`          // 今日推荐用户数排名
		TodayEnterChatUserCnt              int64 `json:"todayEnterChatUserCnt"`              // 今日进入聊天页面用户数
		TodayEnterChatUserCntRank          int64 `json:"todayEnterChatUserCntRank"`          // 今日进入聊天界面用户数排名
		TodayViewUserCnt                   int64 `json:"todayViewUserCnt"`                   // 今日访问资料页面用户数
		TodayViewUserCntRank               int64 `json:"todayViewUserCntRank"`               // 今日资料页面用户数排名
		Last30DaysPaidUserRate             int64 `json:"last30DaysPaidUserRate"`             // 过去30天下单率（下单人数占进入聊天页面的人数比例）
		Last7DaysRepeatPaidUserRate        int64 `json:"last7DaysRepeatPaidUserRate"`        // 过去7天复购率 （复购人数占下单人数比例）
		Last30DaysRepeatPaidUserRate       int64 `json:"last30DaysRepeatPaidUserRate"`       // 过去30天复购率
		Last7DaysAveragePaidAmountPerUser  int64 `json:"last7DaysAveragePaidAmountPerUser"`  // 过去7天人均消费
		Last30DaysAveragePaidAmountPerUser int64 `json:"last30DaysAveragePaidAmountPerUser"` // 过去7天人均消费
		Last7DaysAveragePaidAmountPerDay   int64 `json:"last7DaysAveragePaidAmountPerDay"`   // 过去7天日均消费
		Last30DaysAveragePaidAmountPerDay  int64 `json:"last30DaysAveragePaidAmountPerDay"`  // 过去30天日均消费
		OneStarRatingOrderCnt              int64 `json:"oneStarRatingOrderCnt"`              // 累计不满意评价订单数
		RefundOrderCnt                     int64 `json:"refundOrderCnt"`                     // 累计退款订单数
	}
	fmt.Println(GetStructFieldName(listenerDashboard{}))
}
