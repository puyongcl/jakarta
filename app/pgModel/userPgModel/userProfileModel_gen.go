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
	userProfileFieldNames          = builder.RawFieldNames(&UserProfile{}, true)
	userProfileRows                = strings.Join(userProfileFieldNames, ",")
	userProfileRowsExpectAutoSet   = strings.Join(stringx.Remove(userProfileFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	userProfileRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(userProfileFieldNames, "uid", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaUserProfileUidPrefix = "cache:jakarta:userProfile:uid:"
)

type (
	userProfileModel interface {
		Insert(ctx context.Context, data *UserProfile) (sql.Result, error)
		FindOne(ctx context.Context, uid int64) (*UserProfile, error)
		Update(ctx context.Context, data *UserProfile) error
		Delete(ctx context.Context, uid int64) error
	}

	defaultUserProfileModel struct {
		sqlc.CachedConn
		table string
	}

	UserProfile struct {
		CreateTime    time.Time    `db:"create_time"`
		Uid           int64        `db:"uid"`
		Nickname      string       `db:"nickname"`
		Avatar        string       `db:"avatar"`
		Gender        int64        `db:"gender"`
		Introduction  string       `db:"introduction"`
		UpdateTime    time.Time    `db:"update_time"`
		Constellation int64        `db:"constellation"`
		Birthday      sql.NullTime `db:"birthday"`
		PhoneNumber   string       `db:"phone_number"`
	}
)

func newUserProfileModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserProfileModel {
	return &defaultUserProfileModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."user_profile"`,
	}
}

func (m *defaultUserProfileModel) Delete(ctx context.Context, uid int64) error {
	jakartaUserProfileUidKey := fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, uid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where uid = $1", m.table)
		return conn.ExecCtx(ctx, query, uid)
	}, jakartaUserProfileUidKey)
	return err
}

func (m *defaultUserProfileModel) FindOne(ctx context.Context, uid int64) (*UserProfile, error) {
	jakartaUserProfileUidKey := fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, uid)
	var resp UserProfile
	err := m.QueryRowCtx(ctx, &resp, jakartaUserProfileUidKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where uid = $1 limit 1", userProfileRows, m.table)
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

func (m *defaultUserProfileModel) Insert(ctx context.Context, data *UserProfile) (sql.Result, error) {
	jakartaUserProfileUidKey := fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)", m.table, userProfileRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Uid, data.Nickname, data.Avatar, data.Gender, data.Introduction, data.Constellation, data.Birthday, data.PhoneNumber)
	}, jakartaUserProfileUidKey)
	return ret, err
}

func (m *defaultUserProfileModel) Update(ctx context.Context, data *UserProfile) error {
	jakartaUserProfileUidKey := fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, data.Uid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where uid = $1", m.table, userProfileRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Uid, data.Nickname, data.Avatar, data.Gender, data.Introduction, data.Constellation, data.Birthday, data.PhoneNumber)
	}, jakartaUserProfileUidKey)
	return err
}

func (m *defaultUserProfileModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, primary)
}

func (m *defaultUserProfileModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where uid = $1 limit 1", userProfileRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserProfileModel) tableName() string {
	return m.table
}
