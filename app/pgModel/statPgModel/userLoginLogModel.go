package statPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ UserLoginLogModel = (*customUserLoginLogModel)(nil)

type (
	// UserLoginLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLoginLogModel.
	UserLoginLogModel interface {
		userLoginLogModel
		InsertOrUpdate(ctx context.Context, id string, uid, userType int64, channel string) (sql.Result, error)
		Count(ctx context.Context, start, end *time.Time, channel string, userType int64) (int64, error)
		FindUid(ctx context.Context, start, end *time.Time, channel string, userType int64, pageNo, pageSize int64) ([]int64, error)
	}

	customUserLoginLogModel struct {
		*defaultUserLoginLogModel
	}
)

// NewUserLoginLogModel returns a model for the database table.
func NewUserLoginLogModel(conn sqlx.SqlConn) UserLoginLogModel {
	return &customUserLoginLogModel{
		defaultUserLoginLogModel: newUserLoginLogModel(conn),
	}
}

func (m *defaultUserLoginLogModel) InsertOrUpdate(ctx context.Context, id string, uid, userType int64, channel string) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE SET cnt = %s.cnt + 1", m.table, userLoginLogRowsExpectAutoSet, m.table)
	return m.conn.ExecCtx(ctx, query, id, time.Now(), 1, uid, userType, channel)
}

func (m *defaultUserLoginLogModel) Count(ctx context.Context, start, end *time.Time, channel string, userType int64) (int64, error) {
	rb := squirrel.Select("count(id)").From(m.table).Where("create_time between $1 and $2", start, end).Where("user_type = $3", userType)
	if channel != "" {
		rb = rb.Where("channel = $4", channel)
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var cnt int64
	err = m.conn.QueryRowCtx(ctx, &cnt, query, args...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}

func (m *defaultUserLoginLogModel) FindUid(ctx context.Context, start, end *time.Time, channel string, userType int64, pageNo, pageSize int64) ([]int64, error) {
	rb := squirrel.Select("uid").From(m.table).Where("create_time between $1 and $2", start, end).Where("user_type = $3", userType)
	if channel != "" {
		rb = rb.Where("channel = $4", channel)
	}
	query, args, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var uids []int64
	err = m.conn.QueryRowCtx(ctx, &uids, query, args...)
	switch err {
	case nil:
		return uids, nil
	default:
		return nil, err
	}
}
