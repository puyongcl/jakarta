package userPgModel

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserNeedHelpModel = (*customUserNeedHelpModel)(nil)

type (
	// UserNeedHelpModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserNeedHelpModel.
	UserNeedHelpModel interface {
		userNeedHelpModel
		Find(ctx context.Context, uid int64, targetUid int64, tag, state int64, pageNo, pageSize int64) ([]*UserNeedHelp, error)
		UpdateState(ctx context.Context, id string, state int64, remark string) error
		FindCount(ctx context.Context, uid int64, listenerUid int64, tag, state int64) (int64, error)
	}

	customUserNeedHelpModel struct {
		*defaultUserNeedHelpModel
	}
)

// NewUserNeedHelpModel returns a model for the database table.
func NewUserNeedHelpModel(conn sqlx.SqlConn) UserNeedHelpModel {
	return &customUserNeedHelpModel{
		defaultUserNeedHelpModel: newUserNeedHelpModel(conn),
	}
}

func (m *defaultUserNeedHelpModel) Find(ctx context.Context, uid int64, listenerUid int64, tag, state int64, pageNo, pageSize int64) ([]*UserNeedHelp, error) {
	rb := squirrel.Select(userNeedHelpRows).From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if tag != 0 {
		rb = rb.Where(fmt.Sprintf("$%d = ANY(report_tag)", argNo), tag)
		argNo++
	}
	if state != 0 {
		rb = rb.Where(fmt.Sprintf("state = $%d", argNo), state)
		argNo++
	}

	query, args, err := rb.OrderBy("create_time desc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*UserNeedHelp, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUserNeedHelpModel) UpdateState(ctx context.Context, id string, state int64, remark string) error {
	query := fmt.Sprintf("update %s set state = $1, remark = $2 where id = $3", m.table)
	_, err := m.conn.ExecCtx(ctx, query, state, remark, id)
	return err
}

func (m *defaultUserNeedHelpModel) FindCount(ctx context.Context, uid int64, listenerUid int64, tag, state int64) (int64, error) {
	rb := squirrel.Select("count(id)").From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if tag != 0 {
		rb = rb.Where(fmt.Sprintf("$%d = ANY(report_tag)", argNo), tag)
		argNo++
	}
	if state != 0 {
		rb = rb.Where(fmt.Sprintf("state = $%d", argNo), state)
		argNo++
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
