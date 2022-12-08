package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ ListenerUserViewStatModel = (*customListenerUserViewStatModel)(nil)

type (
	// ListenerUserViewStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerUserViewStatModel.
	ListenerUserViewStatModel interface {
		listenerUserViewStatModel
		InsertOrUpdateView(ctx context.Context, data *ListenerUserViewStat) (sql.Result, error)
		FindCountViewUserCntRangeCreateTime(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error)
	}

	customListenerUserViewStatModel struct {
		*defaultListenerUserViewStatModel
	}
)

// NewListenerUserViewStatModel returns a model for the database table.
func NewListenerUserViewStatModel(conn sqlx.SqlConn) ListenerUserViewStatModel {
	return &customListenerUserViewStatModel{
		defaultListenerUserViewStatModel: newListenerUserViewStatModel(conn),
	}
}

const insertOrUpdateViewSQL = "insert into %s (%s) values ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET view_time = $6, view_cnt = %s.view_cnt + $7"

func (m *defaultListenerUserViewStatModel) InsertOrUpdateView(ctx context.Context, data *ListenerUserViewStat) (sql.Result, error) {
	query := fmt.Sprintf(insertOrUpdateViewSQL, m.table, listenerUserViewStatRowsExpectAutoSet, m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.ViewTime, data.ViewCnt, data.ViewTime, data.ViewCnt)
	return ret, err
}

const queryFindViewUserCntByCreateTimeSQL = "select count(id) from %s where listener_uid = $1 and update_time between $2 and $3"

func (m *defaultListenerUserViewStatModel) FindCountViewUserCntRangeCreateTime(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error) {
	query := fmt.Sprintf(queryFindViewUserCntByCreateTimeSQL, m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, listenerUid, start, end)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
