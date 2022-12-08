// Code generated by goctl. DO NOT EDIT!

package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	listenerStatAverageFieldNames          = builder.RawFieldNames(&ListenerStatAverage{}, true)
	listenerStatAverageRows                = strings.Join(listenerStatAverageFieldNames, ",")
	listenerStatAverageRowsExpectAutoSet   = strings.Join(stringx.Remove(listenerStatAverageFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	listenerStatAverageRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(listenerStatAverageFieldNames, "id", "create_time", "update_time", "create_at", "update_at"))
)

type (
	listenerStatAverageModel interface {
		Insert(ctx context.Context, data *ListenerStatAverage) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*ListenerStatAverage, error)
		Update(ctx context.Context, data *ListenerStatAverage) error
		Delete(ctx context.Context, id string) error
	}

	defaultListenerStatAverageModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ListenerStatAverage struct {
		CreateTime                         time.Time `db:"create_time"`
		UpdateTime                         time.Time `db:"update_time"`
		Id                                 string    `db:"id"`
		YesterdayOrderCnt                  int64     `db:"yesterday_order_cnt"`
		YesterdayOrderCntRank              int64     `db:"yesterday_order_cnt_rank"`
		YesterdayOrderAmount               int64     `db:"yesterday_order_amount"`
		YesterdayOrderAmountRank           int64     `db:"yesterday_order_amount_rank"`
		YesterdayRecommendUserCnt          int64     `db:"yesterday_recommend_user_cnt"`
		YesterdayRecommendUserCntRank      int64     `db:"yesterday_recommend_user_cnt_rank"`
		YesterdayEnterChatUserCnt          int64     `db:"yesterday_enter_chat_user_cnt"`
		YesterdayEnterChatUserCntRank      int64     `db:"yesterday_enter_chat_user_cnt_rank"`
		YesterdayViewUserCnt               int64     `db:"yesterday_view_user_cnt"`
		YesterdayViewUserCntRank           int64     `db:"yesterday_view_user_cnt_rank"`
		Last30DaysPaidUserCnt              int64     `db:"last_30_days_paid_user_cnt"`
		Last30DaysEnterChatUserCnt         int64     `db:"last_30_days_enter_chat_user_cnt"`
		Last30DaysRepeatPaidUserCnt        int64     `db:"last_30_days_repeat_paid_user_cnt"`
		Last30DaysAveragePaidAmountPerUser int64     `db:"last_30_days_average_paid_amount_per_user"`
		Last30DaysAveragePaidAmountPerDay  int64     `db:"last_30_days_average_paid_amount_per_day"`
		Last7DaysPaidUserCnt               int64     `db:"last_7_days_paid_user_cnt"`
		Last7DaysRepeatPaidUserCnt         int64     `db:"last_7_days_repeat_paid_user_cnt"`
		Last7DaysAveragePaidAmountPerUser  int64     `db:"last_7_days_average_paid_amount_per_user"`
		Last7DaysAveragePaidAmountPerDay   int64     `db:"last_7_days_average_paid_amount_per_day"`
		Last30DaysPaidUserRate             int64     `db:"last_30_days_paid_user_rate"`        // 近30天内下单人数占进入聊天页面人数比例(万分)
		Last30DaysRepeatPaidUserRate       int64     `db:"last_30_days_repeat_paid_user_rate"` // 近30天内复购人数占下单人数比例（万分）
		Last7DaysRepeatPaidUserRate        int64     `db:"last_7_days_repeat_paid_user_rate"`  // 近7天内复购人数占下单人数比例
		Last7DaysPaidUserRate              int64     `db:"last_7_days_paid_user_rate"`         // 近7天内下单人数占进入聊天界面用户数比例（万分）
		Last7DaysEnterChatUserCnt          int64     `db:"last_7_days_enter_chat_user_cnt"`    // 最近7天进入聊天界面的人数
		Last30DaysPaidAmountSum            int64     `db:"last_30_days_paid_amount_sum"`       // 近30天下单总金额
		Last7DaysPaidAmountSum             int64     `db:"last_7_days_paid_amount_sum"`        // 近7天下单总金额
	}
)

func newListenerStatAverageModel(conn sqlx.SqlConn) *defaultListenerStatAverageModel {
	return &defaultListenerStatAverageModel{
		conn:  conn,
		table: `"jakarta"."listener_stat_average"`,
	}
}

func (m *defaultListenerStatAverageModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultListenerStatAverageModel) FindOne(ctx context.Context, id string) (*ListenerStatAverage, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", listenerStatAverageRows, m.table)
	var resp ListenerStatAverage
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultListenerStatAverageModel) Insert(ctx context.Context, data *ListenerStatAverage) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)", m.table, listenerStatAverageRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.YesterdayOrderCnt, data.YesterdayOrderCntRank, data.YesterdayOrderAmount, data.YesterdayOrderAmountRank, data.YesterdayRecommendUserCnt, data.YesterdayRecommendUserCntRank, data.YesterdayEnterChatUserCnt, data.YesterdayEnterChatUserCntRank, data.YesterdayViewUserCnt, data.YesterdayViewUserCntRank, data.Last30DaysPaidUserCnt, data.Last30DaysEnterChatUserCnt, data.Last30DaysRepeatPaidUserCnt, data.Last30DaysAveragePaidAmountPerUser, data.Last30DaysAveragePaidAmountPerDay, data.Last7DaysPaidUserCnt, data.Last7DaysRepeatPaidUserCnt, data.Last7DaysAveragePaidAmountPerUser, data.Last7DaysAveragePaidAmountPerDay, data.Last30DaysPaidUserRate, data.Last30DaysRepeatPaidUserRate, data.Last7DaysRepeatPaidUserRate, data.Last7DaysPaidUserRate, data.Last7DaysEnterChatUserCnt, data.Last30DaysPaidAmountSum, data.Last7DaysPaidAmountSum)
	return ret, err
}

func (m *defaultListenerStatAverageModel) Update(ctx context.Context, data *ListenerStatAverage) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, listenerStatAverageRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Id, data.YesterdayOrderCnt, data.YesterdayOrderCntRank, data.YesterdayOrderAmount, data.YesterdayOrderAmountRank, data.YesterdayRecommendUserCnt, data.YesterdayRecommendUserCntRank, data.YesterdayEnterChatUserCnt, data.YesterdayEnterChatUserCntRank, data.YesterdayViewUserCnt, data.YesterdayViewUserCntRank, data.Last30DaysPaidUserCnt, data.Last30DaysEnterChatUserCnt, data.Last30DaysRepeatPaidUserCnt, data.Last30DaysAveragePaidAmountPerUser, data.Last30DaysAveragePaidAmountPerDay, data.Last7DaysPaidUserCnt, data.Last7DaysRepeatPaidUserCnt, data.Last7DaysAveragePaidAmountPerUser, data.Last7DaysAveragePaidAmountPerDay, data.Last30DaysPaidUserRate, data.Last30DaysRepeatPaidUserRate, data.Last7DaysRepeatPaidUserRate, data.Last7DaysPaidUserRate, data.Last7DaysEnterChatUserCnt, data.Last30DaysPaidAmountSum, data.Last7DaysPaidAmountSum)
	return err
}

func (m *defaultListenerStatAverageModel) tableName() string {
	return m.table
}
