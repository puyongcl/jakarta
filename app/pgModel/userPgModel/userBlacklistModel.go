package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserBlacklistModel = (*customUserBlacklistModel)(nil)

type (
	// UserBlacklistModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserBlacklistModel.
	UserBlacklistModel interface {
		userBlacklistModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		InsertOrUpdate(ctx context.Context, session sqlx.Session, data *UserBlacklist) (sql.Result, error)
		Find(ctx context.Context, uid int64, targetUid int64, state int64, pageNo, pageSize int64) ([]*UserBlacklist, error)
	}

	customUserBlacklistModel struct {
		*defaultUserBlacklistModel
	}
)

// NewUserBlacklistModel returns a model for the database table.
func NewUserBlacklistModel(conn sqlx.SqlConn) UserBlacklistModel {
	return &customUserBlacklistModel{
		defaultUserBlacklistModel: newUserBlacklistModel(conn),
	}
}

func (m *defaultUserBlacklistModel) InsertOrUpdate(ctx context.Context, session sqlx.Session, data *UserBlacklist) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO UPDATE SET state = $7, target_avatar = $8, target_nick_name = $9", m.table, userBlacklistRowsExpectAutoSet)
	return session.ExecCtx(ctx, query, data.Uid, data.Id, data.TargetUid, data.State, data.TargetAvatar, data.TargetNickName, data.State, data.TargetAvatar, data.TargetNickName)
}

func (m *defaultUserBlacklistModel) Find(ctx context.Context, uid int64, targetUid int64, state int64, pageNo, pageSize int64) ([]*UserBlacklist, error) {
	if uid != 0 { // 查询拉黑了谁
		query, args, err := squirrel.Select(userBlacklistRows).From(m.table).Where("uid = $1", uid).Where("state = $2", state).OrderBy("create_time desc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
		if err != nil {
			return nil, err
		}
		resp := make([]*UserBlacklist, 0)
		err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
		switch err {
		case nil:
			return resp, nil
		default:
			return nil, err
		}

	}
	// 查询被谁拉黑
	query, args, err := squirrel.Select("uid").From(m.table).Where("target_uid = $1", targetUid).Where("state = $2", state).OrderBy("create_time desc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	resp := make([]*UserBlacklist, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUserBlacklistModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}
