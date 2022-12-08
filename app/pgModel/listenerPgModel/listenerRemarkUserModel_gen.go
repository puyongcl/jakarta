// Code generated by goctl. DO NOT EDIT!

package listenerPgModel

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
	listenerRemarkUserFieldNames          = builder.RawFieldNames(&ListenerRemarkUser{}, true)
	listenerRemarkUserRows                = strings.Join(listenerRemarkUserFieldNames, ",")
	listenerRemarkUserRowsExpectAutoSet   = strings.Join(stringx.Remove(listenerRemarkUserFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	listenerRemarkUserRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(listenerRemarkUserFieldNames, "id", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaListenerRemarkUserIdPrefix = "cache:jakarta:listenerRemarkUser:id:"
)

type (
	listenerRemarkUserModel interface {
		Insert(ctx context.Context, data *ListenerRemarkUser) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*ListenerRemarkUser, error)
		Update(ctx context.Context, data *ListenerRemarkUser) error
		Delete(ctx context.Context, id string) error
	}

	defaultListenerRemarkUserModel struct {
		sqlc.CachedConn
		table string
	}

	ListenerRemarkUser struct {
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		Id          string    `db:"id"`
		ListenerUid int64     `db:"listener_uid"`
		Uid         int64     `db:"uid"`
		Remark      string    `db:"remark"`    // 名称备注
		UserDesc    string    `db:"user_desc"` // 对用户的描述
	}
)

func newListenerRemarkUserModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultListenerRemarkUserModel {
	return &defaultListenerRemarkUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."listener_remark_user"`,
	}
}

func (m *defaultListenerRemarkUserModel) Delete(ctx context.Context, id string) error {
	jakartaListenerRemarkUserIdKey := fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, jakartaListenerRemarkUserIdKey)
	return err
}

func (m *defaultListenerRemarkUserModel) FindOne(ctx context.Context, id string) (*ListenerRemarkUser, error) {
	jakartaListenerRemarkUserIdKey := fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, id)
	var resp ListenerRemarkUser
	err := m.QueryRowCtx(ctx, &resp, jakartaListenerRemarkUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = $1 limit 1", listenerRemarkUserRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
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

func (m *defaultListenerRemarkUserModel) Insert(ctx context.Context, data *ListenerRemarkUser) (sql.Result, error) {
	jakartaListenerRemarkUserIdKey := fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5)", m.table, listenerRemarkUserRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.ListenerUid, data.Uid, data.Remark, data.UserDesc)
	}, jakartaListenerRemarkUserIdKey)
	return ret, err
}

func (m *defaultListenerRemarkUserModel) Update(ctx context.Context, data *ListenerRemarkUser) error {
	jakartaListenerRemarkUserIdKey := fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = $1", m.table, listenerRemarkUserRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Id, data.ListenerUid, data.Uid, data.Remark, data.UserDesc)
	}, jakartaListenerRemarkUserIdKey)
	return err
}

func (m *defaultListenerRemarkUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, primary)
}

func (m *defaultListenerRemarkUserModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", listenerRemarkUserRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultListenerRemarkUserModel) tableName() string {
	return m.table
}