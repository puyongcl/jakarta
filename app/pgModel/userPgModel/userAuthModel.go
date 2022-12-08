package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/userkey"
	"time"
)

var _ UserAuthModel = (*customUserAuthModel)(nil)

type (
	// UserAuthModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAuthModel.
	UserAuthModel interface {
		userAuthModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		FindByUserType(ctx context.Context, checkStatus int64, pageNo, pageSize int64) ([]*UserAuth, error)
		UpdateUserType(ctx context.Context, uid, userType int64) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *UserAuth) (sql.Result, error)
		UpdateBanAccount(ctx context.Context, uid int64, freeTime time.Time, banReason string) error
		FindUserByCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, userType int64, channel string, pageNo, pageSize int64) ([]*UserAuth, error)
		FindUserByCreateTimeCount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, userType int64, channel string) (int64, error)
		DeleteUserAccount(ctx context.Context, uid int64) (*UserAuth, error)
		CountNewUser(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		FindChannelList(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time) ([]string, error)
		FindFirstUid(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time) (int64, error)
		FindLastUid(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time) (int64, error)
		FindUid(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string, pageNo, pageSize int64) ([]int64, error)
		FindOneByAuthKeyAuthType2(ctx context.Context, authKey string, authType string) (*UserAuth, error)
	}

	customUserAuthModel struct {
		*defaultUserAuthModel
	}
)

// NewUserAuthModel returns a model for the database table.
func NewUserAuthModel(conn sqlx.SqlConn, c cache.CacheConf) UserAuthModel {
	return &customUserAuthModel{
		defaultUserAuthModel: newUserAuthModel(conn, c),
	}
}

// export logic
func (m *defaultUserAuthModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultUserAuthModel) FindByUserType(ctx context.Context, userType int64, pageNo, pageSize int64) ([]*UserAuth, error) {
	query := fmt.Sprintf("select %s from %s where user_type = $1 order by create_time asc limit $2 offset $3", userAuthRows, m.table)
	var resp []*UserAuth
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userType, pageSize, (pageNo-1)*pageSize)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) UpdateUserType(ctx context.Context, uid, userType int64) error {
	data, err := m.FindOne(ctx, uid)
	if err != nil {
		return err
	}

	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, data.Uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set user_type = $1 where uid = $2", m.table)
		return conn.ExecCtx(ctx, query, userType, uid)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return err
}

func (m *defaultUserAuthModel) InsertTrans(ctx context.Context, session sqlx.Session, data *UserAuth) (sql.Result, error) {
	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, userAuthRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Uid, data.AuthKey, data.AuthType, data.Password, data.AccountState, data.UserType, data.FreeTime, data.BanReason, data.Channel)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return ret, err
}

// 封禁不改变account_state 动态根据封禁时间判断状态
func (m *defaultUserAuthModel) UpdateBanAccount(ctx context.Context, uid int64, freeTime time.Time, banReason string) error {
	data, err := m.FindOne(ctx, uid)
	if err != nil {
		return err
	}

	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, data.Uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set free_time = $1, ban_reason = $2 where uid = $3", m.table)
		return conn.ExecCtx(ctx, query, freeTime, banReason, uid)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return err
}

func (m *defaultUserAuthModel) FindUserByCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, userType int64, channel string, pageNo, pageSize int64) ([]*UserAuth, error) {
	rb := squirrel.Select(userAuthRows).From(m.table)
	argNo := 1
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf("create_time > $%d", argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time < $%d", argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time between $%d and $%d", argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	if userType != 0 {
		rb = rb.Where(fmt.Sprintf("user_type = $%d", argNo), userType)
		argNo++
	}
	if channel != "" {
		rb = rb.Where(fmt.Sprintf("channel = $%d", argNo), channel)
		argNo++
	}
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query, args, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*UserAuth, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindUserByCreateTimeCount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, userType int64, channel string) (int64, error) {
	rb := squirrel.Select("count(uid)").From(m.table)
	argNo := 1
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf("create_time > $%d", argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time < $%d", argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time between $%d and $%d", argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	if userType != 0 {
		rb = rb.Where(fmt.Sprintf("user_type = $%d", argNo), userType)
		argNo++
	}
	if channel != "" {
		rb = rb.Where(fmt.Sprintf("channel = $%d", argNo), channel)
		argNo++
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var cnt int64
	err = m.QueryRowNoCacheCtx(ctx, &cnt, query, args...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}

// 删除用户账户
func (m *defaultUserAuthModel) DeleteUserAccount(ctx context.Context, uid int64) (*UserAuth, error) {
	data, err := m.FindOne(ctx, uid)
	if err != nil {
		return data, err
	}

	if data.AccountState == userkey.AccountStateCancel {
		return data, nil
	}

	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, data.AuthKey, data.AuthType)
	jakartaUserAuthUidKey := fmt.Sprintf("%s%v", cacheJakartaUserAuthUidPrefix, data.Uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set account_state = $1 where uid = $2", m.table)
		return conn.ExecCtx(ctx, query, userkey.AccountStateCancel, uid)
	}, jakartaUserAuthAuthKeyAuthTypeKey, jakartaUserAuthUidKey)
	return data, err
}

// 统计新用户
func (m *defaultUserAuthModel) CountNewUser(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	rb := squirrel.Select("count(uid)").From(m.table).Where("user_type = 2")
	argNo := 1
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf("create_time > $%d", argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time < $%d", argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time between $%d and $%d", argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	if channel != "" {
		rb = rb.Where(fmt.Sprintf("channel = $%d", argNo), channel)
		argNo++
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var cnt int64
	err = m.QueryRowNoCacheCtx(ctx, &cnt, query, args...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}

func (m *defaultUserAuthModel) FindChannelList(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time) ([]string, error) {
	rb := squirrel.Select("channel").From(m.table)
	argNo := 1
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf("create_time > $%d", argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time < $%d", argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("create_time between $%d and $%d", argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	query, args, err := rb.GroupBy("channel").ToSql()
	if err != nil {
		return nil, err
	}

	var list []string
	err = m.QueryRowsNoCacheCtx(ctx, &list, query, args...)
	switch err {
	case nil:
		return list, nil
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindFirstUid(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time) (int64, error) {
	q, a, err := squirrel.Select("uid").From(m.table).Where("create_time between $1 and $2", createTimeStart, createTimeEnd).OrderBy("create_time ASC").Limit(1).ToSql()
	if err != nil {
		return 0, err
	}
	var uid int64
	err = m.QueryRowNoCacheCtx(ctx, &uid, q, a...)
	switch err {
	case nil:
		return uid, nil
	default:
		return 0, err
	}
}

func (m *defaultUserAuthModel) FindLastUid(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time) (int64, error) {
	q, a, err := squirrel.Select("uid").From(m.table).Where("create_time between $1 and $2", createTimeStart, createTimeEnd).OrderBy("create_time DESC").Limit(1).ToSql()
	if err != nil {
		return 0, err
	}
	var uid int64
	err = m.QueryRowNoCacheCtx(ctx, &uid, q, a...)
	switch err {
	case nil:
		return uid, nil
	default:
		return 0, err
	}
}

func (m *defaultUserAuthModel) FindUid(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string, pageNo, pageSize int64) ([]int64, error) {
	rb := squirrel.Select("uid").From(m.table).Where("create_time between $1 and $2", createTimeStart, createTimeEnd)
	if channel != "" {
		rb = rb.Where("channel = $3", channel)
	}

	q, a, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	var uid []int64
	err = m.QueryRowsNoCacheCtx(ctx, &uid, q, a...)
	switch err {
	case nil:
		return uid, nil
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByAuthKeyAuthType2(ctx context.Context, authKey string, authType string) (*UserAuth, error) {
	jakartaUserAuthAuthKeyAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheJakartaUserAuthAuthKeyAuthTypePrefix, authKey, authType)
	var resp UserAuth
	err := m.QueryRowIndexCtx(ctx, &resp, jakartaUserAuthAuthKeyAuthTypeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where auth_key = $1 and auth_type = $2 and account_state != $3 limit 1", userAuthRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, authKey, authType, userkey.AccountStateCancel); err != nil {
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
