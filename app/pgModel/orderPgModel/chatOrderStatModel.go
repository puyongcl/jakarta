package orderPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatOrderStatModel = (*customChatOrderStatModel)(nil)

type (
	// ChatOrderStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatOrderStatModel.
	ChatOrderStatModel interface {
		chatOrderStatModel
		Update2(ctx context.Context, id string, addMin int64) error
		Find(ctx context.Context, pageNo, pageSize int64) ([]*ChatOrderStat, error)
		ResetBuyMinute(ctx context.Context, id string, min int64) error
	}

	customChatOrderStatModel struct {
		*defaultChatOrderStatModel
	}
)

// NewChatOrderStatModel returns a model for the database table.
func NewChatOrderStatModel(conn sqlx.SqlConn, c cache.CacheConf) ChatOrderStatModel {
	return &customChatOrderStatModel{
		defaultChatOrderStatModel: newChatOrderStatModel(conn, c),
	}
}

func (m *defaultChatOrderStatModel) Update2(ctx context.Context, id string, addMin int64) error {
	jakartaChatOrderStatIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderStatIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set confirm_minute_sum = confirm_minute_sum + $1 where id = $2", m.table)
		return conn.ExecCtx(ctx, query, addMin, id)
	}, jakartaChatOrderStatIdKey)
	return err
}

func (m *defaultChatOrderStatModel) Find(ctx context.Context, pageNo, pageSize int64) ([]*ChatOrderStat, error) {
	rb := squirrel.Select(chatOrderStatRows).From(m.table).OrderBy("create_time desc")

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query, values, err := rb.Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatOrderStat, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatOrderStatModel) ResetBuyMinute(ctx context.Context, id string, min int64) error {
	jakartaChatOrderStatIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderStatIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set confirm_minute_sum = $1 where id = $2", m.table)
		return conn.ExecCtx(ctx, query, min, id)
	}, jakartaChatOrderStatIdKey)
	return err
}
