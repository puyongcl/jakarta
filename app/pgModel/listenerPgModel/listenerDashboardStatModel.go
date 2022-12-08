package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ ListenerDashboardStatModel = (*customListenerDashboardStatModel)(nil)

type (
	// ListenerDashboardStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerDashboardStatModel.
	ListenerDashboardStatModel interface {
		listenerDashboardStatModel
		InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerDashboardStat) (sql.Result, error)
		UpdateLastDaysStat(ctx context.Context, listenerUid int64, data *ListenerDashboardStatLastDays) error
		CalculateAverage(ctx context.Context, field string) (float64, error)
		UpdateSuggestion(ctx context.Context, listenerUid int64, sg []int64) error
	}

	customListenerDashboardStatModel struct {
		*defaultListenerDashboardStatModel
	}

	ListenerDashboardStatLastDays struct {
		Last30DaysPaidUserCnt              int64     `db:"last_30_days_paid_user_cnt"`
		Last30DaysEnterChatUserCnt         int64     `db:"last_30_days_enter_chat_user_cnt"`
		Last30DaysRepeatPaidUserCnt        int64     `db:"last_30_days_repeat_paid_user_cnt"`
		Last30DaysAveragePaidAmountPerUser int64     `db:"last_30_days_average_paid_amount_per_user"`
		Last30DaysAveragePaidAmountPerDay  int64     `db:"last_30_days_average_paid_amount_per_day"`
		Last7DaysPaidUserCnt               int64     `db:"last_7_days_paid_user_cnt"`
		Last7DaysRepeatPaidUserCnt         int64     `db:"last_7_days_repeat_paid_user_cnt"`
		Last7DaysAveragePaidAmountPerUser  int64     `db:"last_7_days_average_paid_amount_per_user"`
		Last7DaysAveragePaidAmountPerDay   int64     `db:"last_7_days_average_paid_amount_per_day"`
		LastDayLastUpdateTime              time.Time `db:"last_day_last_update_time"`          // 近几日数据快照时间
		Last30DaysPaidUserRate             int64     `db:"last_30_days_paid_user_rate"`        // 近30天内下单人数占进入聊天页面人数比例(万分)
		Last30DaysRepeatPaidUserRate       int64     `db:"last_30_days_repeat_paid_user_rate"` // 近30天内复购人数占下单人数比例（万分）
		Last7DaysRepeatPaidUserRate        int64     `db:"last_7_days_repeat_paid_user_rate"`  // 近7天内复购人数占下单人数比例
		Last7DaysPaidUserRate              int64     `db:"last_7_days_paid_user_rate"`         // 近7天内下单人数占进入聊天界面用户数比例（万分）
		Last7DaysEnterChatUserCnt          int64     `db:"last_7_days_enter_chat_user_cnt"`    // 最近7天进入聊天界面的人数
		Last30DaysPaidAmountSum            int64     `db:"last_30_days_paid_amount_sum"`       // 近30天下单总金额
		Last7DaysPaidAmountSum             int64     `db:"last_7_days_paid_amount_sum"`        // 近7天下单总金额
	}
)

var listenerDashboardStatLastDaysStatRowsWithPlaceHolder = builder.PostgreSqlJoin(builder.RawFieldNames(&ListenerDashboardStatLastDays{}, true))

// NewListenerDashboardStatModel returns a model for the database table.
func NewListenerDashboardStatModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerDashboardStatModel {
	return &customListenerDashboardStatModel{
		defaultListenerDashboardStatModel: newListenerDashboardStatModel(conn, c),
	}
}

func (m *defaultListenerDashboardStatModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerDashboardStat) (sql.Result, error) {
	jakartaListenerDashboardStatListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerDashboardStatListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)", m.table, listenerDashboardStatRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.ListenerUid, data.YesterdayOrderCnt, data.YesterdayOrderCntRank, data.YesterdayOrderAmount, data.YesterdayOrderAmountRank, data.YesterdayRecommendUserCnt, data.YesterdayRecommendUserCntRank, data.YesterdayEnterChatUserCnt, data.YesterdayEnterChatUserCntRank, data.YesterdayViewUserCnt, data.YesterdayViewUserCntRank, data.YesterdayLastUpdateTime, data.Last30DaysPaidUserCnt, data.Last30DaysEnterChatUserCnt, data.Last30DaysRepeatPaidUserCnt, data.Last30DaysAveragePaidAmountPerUser, data.Last30DaysAveragePaidAmountPerDay, data.Last7DaysPaidUserCnt, data.Last7DaysRepeatPaidUserCnt, data.Last7DaysAveragePaidAmountPerUser, data.Last7DaysAveragePaidAmountPerDay, data.LastDayLastUpdateTime, data.Last30DaysPaidUserRate, data.Last30DaysRepeatPaidUserRate, data.Last7DaysRepeatPaidUserRate, data.Last7DaysPaidUserRate, data.Last7DaysEnterChatUserCnt, data.Last30DaysPaidAmountSum, data.Last7DaysPaidAmountSum, data.Suggestion)
	}, jakartaListenerDashboardStatListenerUidKey)
	return ret, err
}

func (m *defaultListenerDashboardStatModel) UpdateLastDaysStat(ctx context.Context, listenerUid int64, data *ListenerDashboardStatLastDays) error {
	jakartaListenerDashboardStatListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerDashboardStatListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where listener_uid = $1", m.table, listenerDashboardStatLastDaysStatRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, listenerUid, data.Last30DaysPaidUserCnt, data.Last30DaysEnterChatUserCnt, data.Last30DaysRepeatPaidUserCnt, data.Last30DaysAveragePaidAmountPerUser, data.Last30DaysAveragePaidAmountPerDay, data.Last7DaysPaidUserCnt, data.Last7DaysRepeatPaidUserCnt, data.Last7DaysAveragePaidAmountPerUser, data.Last7DaysAveragePaidAmountPerDay, data.LastDayLastUpdateTime, data.Last30DaysPaidUserRate, data.Last30DaysRepeatPaidUserRate, data.Last7DaysRepeatPaidUserRate, data.Last7DaysPaidUserRate, data.Last7DaysEnterChatUserCnt, data.Last30DaysPaidAmountSum, data.Last7DaysPaidAmountSum)
	}, jakartaListenerDashboardStatListenerUidKey)
	return err
}

const calculateFieldAverageSQL = "select coalesce(avg(%s), 0) from %s where update_time > $1 and %s != 0"

func (m *defaultListenerDashboardStatModel) CalculateAverage(ctx context.Context, field string) (float64, error) {
	query := fmt.Sprintf(calculateFieldAverageSQL, field, m.table, field)
	t := time.Now().Add(-60 * time.Minute)
	var cnt float64
	err := m.QueryRowNoCacheCtx(ctx, &cnt, query, t)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (m *defaultListenerDashboardStatModel) UpdateSuggestion(ctx context.Context, listenerUid int64, sg []int64) error {
	jakartaListenerDashboardStatListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerDashboardStatListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set suggestion = $1 where listener_uid = $2", m.table)
		return conn.ExecCtx(ctx, query, pq.Int64Array(sg), listenerUid)
	}, jakartaListenerDashboardStatListenerUidKey)
	return err
}
