package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ ListenerUserRecommendStatModel = (*customListenerUserRecommendStatModel)(nil)

type (
	// ListenerUserRecommendStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerUserRecommendStatModel.
	ListenerUserRecommendStatModel interface {
		listenerUserRecommendStatModel
		InsertOrUpdateRecommend(ctx context.Context, data *ListenerUserRecommendStat) (sql.Result, error)
		FindCountRecommendUserCntRangeCreateTime(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error)
	}

	customListenerUserRecommendStatModel struct {
		*defaultListenerUserRecommendStatModel
	}
)

// NewListenerUserRecommendStatModel returns a model for the database table.
func NewListenerUserRecommendStatModel(conn sqlx.SqlConn) ListenerUserRecommendStatModel {
	return &customListenerUserRecommendStatModel{
		defaultListenerUserRecommendStatModel: newListenerUserRecommendStatModel(conn),
	}
}

const insertOrUpdateRecommendSQL = "insert into %s (%s) values ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET recommend_time = $6, recommend_cnt = %s.recommend_cnt + $7"

func (m *defaultListenerUserRecommendStatModel) InsertOrUpdateRecommend(ctx context.Context, data *ListenerUserRecommendStat) (sql.Result, error) {
	query := fmt.Sprintf(insertOrUpdateRecommendSQL, m.table, listenerUserRecommendStatRowsExpectAutoSet, m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.RecommendTime, data.RecommendCnt, data.RecommendTime, data.RecommendCnt)
	return ret, err
}

const queryFindRecommendUserCntByCreateTimeSQL = "select count(id) from %s where listener_uid = $1 and update_time between $2 and $3"

func (m *defaultListenerUserRecommendStatModel) FindCountRecommendUserCntRangeCreateTime(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error) {
	query := fmt.Sprintf(queryFindRecommendUserCntByCreateTimeSQL, m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, listenerUid, start, end)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
