package adminPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ Contract1021Model = (*customContract1021Model)(nil)

type (
	// Contract1021Model is an interface to be customized, add more methods here,
	// and implement the added methods in customContract1021Model.
	Contract1021Model interface {
		contract1021Model
		Find(ctx context.Context, pageNo, pageSize int64) ([]*Contract1021, error)
		Count(ctx context.Context) (int64, error)
		Sign(ctx context.Context, id, signName string) error
		QuerySignTime(ctx context.Context, id string) (*time.Time, error)
	}

	customContract1021Model struct {
		*defaultContract1021Model
	}
)

// NewContract1021Model returns a model for the database table.
func NewContract1021Model(conn sqlx.SqlConn) Contract1021Model {
	return &customContract1021Model{
		defaultContract1021Model: newContract1021Model(conn),
	}
}

func (m *defaultContract1021Model) Find(ctx context.Context, pageNo, pageSize int64) ([]*Contract1021, error) {
	q, a, err := squirrel.Select("*").From(m.table).OrderBy("create_time desc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var resp []*Contract1021
	err = m.conn.QueryRowsCtx(ctx, &resp, q, a...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultContract1021Model) Count(ctx context.Context) (int64, error) {
	q, a, err := squirrel.Select("count(contract_id)").From(m.table).Where("contract_id != $1", "").ToSql()
	if err != nil {
		return 0, err
	}
	var cnt int64
	err = m.conn.QueryRowCtx(ctx, &cnt, q, a...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}

const signContract1021SQL = `update %s set sign_name = $1, sign_time = $2 where contract_id = $3`

func (m *defaultContract1021Model) Sign(ctx context.Context, id, signName string) error {
	q := fmt.Sprintf(signContract1021SQL, m.table)
	_, err := m.conn.ExecCtx(ctx, q, signName, time.Now(), id)
	return err
}

type QueryContract1021SignTime struct {
	SignTime sql.NullTime `db:"sign_time"`
}

const queryContract1021SignTimeSQL = `select sign_time from %s where contract_id = $1`

func (m *defaultContract1021Model) QuerySignTime(ctx context.Context, id string) (*time.Time, error) {
	q := fmt.Sprintf(queryContract1021SignTimeSQL, m.table)
	var rs QueryContract1021SignTime
	err := m.conn.QueryRowCtx(ctx, &rs, q, id)
	switch err {
	case nil:
		if !rs.SignTime.Valid {
			return nil, nil
		}
		return &rs.SignTime.Time, nil
	default:
		return nil, err
	}
}
