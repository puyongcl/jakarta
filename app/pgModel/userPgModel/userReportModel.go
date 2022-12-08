package userPgModel

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserReportModel = (*customUserReportModel)(nil)

type (
	// UserReportModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserReportModel.
	UserReportModel interface {
		userReportModel
		Find(ctx context.Context, uid int64, targetUid int64, tag, state int64, pageNo, pageSize int64) ([]*UserReport, error)
		UpdateState(ctx context.Context, id string, state int64, remark string) error
		FindCount(ctx context.Context, uid int64, targetUid int64, tag, state int64) (int64, error)
	}

	customUserReportModel struct {
		*defaultUserReportModel
	}
)

// NewUserReportModel returns a model for the database table.
func NewUserReportModel(conn sqlx.SqlConn) UserReportModel {
	return &customUserReportModel{
		defaultUserReportModel: newUserReportModel(conn),
	}
}

func (m *defaultUserReportModel) Find(ctx context.Context, uid int64, targetUid int64, tag, state int64, pageNo, pageSize int64) ([]*UserReport, error) {
	rb := squirrel.Select(userReportRows).From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if targetUid != 0 {
		rb = rb.Where(fmt.Sprintf("target_uid = $%d", argNo), targetUid)
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

	resp := make([]*UserReport, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUserReportModel) UpdateState(ctx context.Context, id string, state int64, remark string) error {
	query := fmt.Sprintf("update %s set state = $1, remark = $2 where id = $3", m.table)
	_, err := m.conn.ExecCtx(ctx, query, state, remark, id)
	return err
}

func (m *defaultUserReportModel) FindCount(ctx context.Context, uid int64, targetUid int64, tag, state int64) (int64, error) {
	rb := squirrel.Select("count(id)").From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if targetUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), targetUid)
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
