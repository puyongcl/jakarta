// Code generated by goctl. DO NOT EDIT!

package chatPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userVoiceChatStateFieldNames          = builder.RawFieldNames(&UserVoiceChatState{}, true)
	userVoiceChatStateRows                = strings.Join(userVoiceChatStateFieldNames, ",")
	userVoiceChatStateRowsExpectAutoSet   = strings.Join(stringx.Remove(userVoiceChatStateFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	userVoiceChatStateRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(userVoiceChatStateFieldNames, "uid", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaUserVoiceChatStateUidPrefix = "cache:jakarta:userVoiceChatState:uid:"
)

type (
	userVoiceChatStateModel interface {
		Insert(ctx context.Context, data *UserVoiceChatState) (sql.Result, error)
		FindOne(ctx context.Context, uid int64) (*UserVoiceChatState, error)
		Update(ctx context.Context, data *UserVoiceChatState) error
		Delete(ctx context.Context, uid int64) error
	}

	defaultUserVoiceChatStateModel struct {
		sqlc.CachedConn
		table string
	}

	UserVoiceChatState struct {
		CreateTime  time.Time    `db:"create_time"`
		UpdateTime  time.Time    `db:"update_time"`
		Uid         int64        `db:"uid"`
		ListenerUid int64        `db:"listener_uid"`
		State       int64        `db:"state"`
		StartTime   sql.NullTime `db:"start_time"`
		EndTime     sql.NullTime `db:"end_time"`
		SettleTime  sql.NullTime `db:"settle_time"`
	}
)

func newUserVoiceChatStateModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserVoiceChatStateModel {
	return &defaultUserVoiceChatStateModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."user_voice_chat_state"`,
	}
}

func (m *defaultUserVoiceChatStateModel) Delete(ctx context.Context, uid int64) error {
	jakartaUserVoiceChatStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, uid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where uid = $1", m.table)
		return conn.ExecCtx(ctx, query, uid)
	}, jakartaUserVoiceChatStateUidKey)
	return err
}

func (m *defaultUserVoiceChatStateModel) FindOne(ctx context.Context, uid int64) (*UserVoiceChatState, error) {
	jakartaUserVoiceChatStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, uid)
	var resp UserVoiceChatState
	err := m.QueryRowCtx(ctx, &resp, jakartaUserVoiceChatStateUidKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where uid = $1 limit 1", userVoiceChatStateRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, uid)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserVoiceChatStateModel) Insert(ctx context.Context, data *UserVoiceChatState) (sql.Result, error) {
	jakartaUserVoiceChatStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6)", m.table, userVoiceChatStateRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Uid, data.ListenerUid, data.State, data.StartTime, data.EndTime, data.SettleTime)
	}, jakartaUserVoiceChatStateUidKey)
	return ret, err
}

func (m *defaultUserVoiceChatStateModel) Update(ctx context.Context, data *UserVoiceChatState) error {
	jakartaUserVoiceChatStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, data.Uid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where uid = $1", m.table, userVoiceChatStateRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Uid, data.ListenerUid, data.State, data.StartTime, data.EndTime, data.SettleTime)
	}, jakartaUserVoiceChatStateUidKey)
	return err
}

func (m *defaultUserVoiceChatStateModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, primary)
}

func (m *defaultUserVoiceChatStateModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where uid = $1 limit 1", userVoiceChatStateRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserVoiceChatStateModel) tableName() string {
	return m.table
}
