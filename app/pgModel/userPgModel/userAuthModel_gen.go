// Code generated by goctl. DO NOT EDIT!

package userPgModel

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
	userAuthFieldNames          = builder.RawFieldNames(&UserAuth{}, true)
	userAuthRows                = strings.Join(userAuthFieldNames, ",")
	userAuthRowsExpectAutoSet   = strings.Join(stringx.Remove(userAuthFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	userAuthRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(userAuthFieldNames, "uid", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaUserAuthUidPrefix             = "cache:jakarta:userAuth:uid:"
	cacheJakartaUserAuthAuthKeyAuthTypePrefix = "cache:jakarta:userAuth:authKey:authType:"
)

type (
	userAuthModel interface {
		Insert(ctx context.Context, data *UserAuth) (sql.Result, error)
		FindOne(ctx context.Context, uid int64) (*UserAuth, error)
		FindOneByAuthKeyAuthType(ctx context.Context, authKey string, authType string) (*UserAuth, error)
		Update(ctx context.Context, data *UserAuth) error
		Delete(ctx context.Context, uid int64) error
	}

	defaultUserAuthModel struct {
		sqlc.CachedConn
		table string
	}

	UserAuth struct {
		CreateTime   time.Time    `db:"create_time"`
		UpdateTime   time.Time    `db:"update_time"`
		Uid          int64        `db:"uid"`
		AuthKey      string       `db:"auth_key"`  // 平台唯一id
		AuthType     string       `db:"auth_type"` // 登陆类型
		Password     string       `db:"password"`
		AccountState int64        `db:"account_state"` // 2 正常 6 封禁 8 注销
		UserType     int64        `db:"user_type"`     // 2 普通用户 4 XXX 6 管理员
		FreeTime     sql.NullTime `db:"free_time"`     // 解封时间
		BanReason    string       `db:"ban_reason"`    // 封禁原因
		Channel      string       `db:"channel"`       // 获客渠道
	}
)

func newUserAuthModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserAuthModel {
	return &defaultUserAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."user_auth"`,
	}
}

func (m *defaultUserAuthModel) Delete(ctx context.Context, uid int64) error {
	data, err := m.FindOne(ctx, uid)
	if err != nil {
		return err
	}

	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where uid = $1", m.table)
		return conn.ExecCtx(ctx, query, uid)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return err
}

func (m *defaultUserAuthModel) FindOne(ctx context.Context, uid int64) (*UserAuth, error) {
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, uid)
	var resp UserAuth
	err := m.QueryRowCtx(ctx, &resp, jakartaUserAuthUidKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where uid = $1 limit 1", userAuthRows, m.table)
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

func (m *defaultUserAuthModel) FindOneByAuthKeyAuthType(ctx context.Context, authKey string, authType string) (*UserAuth, error) {
	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, authKey, authType)
	var resp UserAuth
	err := m.QueryRowIndexCtx(ctx, &resp, jakartaUserAuthAuthKeyAuthTypeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where auth_key = $1 and auth_type = $2 limit 1", userAuthRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, authKey, authType); err != nil {
			return nil, err
		}
		return resp.Uid, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) Insert(ctx context.Context, data *UserAuth) (sql.Result, error) {
	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, userAuthRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Uid, data.AuthKey, data.AuthType, data.Password, data.AccountState, data.UserType, data.FreeTime, data.BanReason, data.Channel)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return ret, err
}

func (m *defaultUserAuthModel) Update(ctx context.Context, newData *UserAuth) error {
	data, err := m.FindOne(ctx, newData.Uid)
	if err != nil {
		return err
	}

	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, data.Uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where uid = $1", m.table, userAuthRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.Uid, newData.AuthKey, newData.AuthType, newData.Password, newData.AccountState, newData.UserType, newData.FreeTime, newData.BanReason, newData.Channel)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return err
}

func (m *defaultUserAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, primary)
}

func (m *defaultUserAuthModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where uid = $1 limit 1", userAuthRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserAuthModel) tableName() string {
	return m.table
}