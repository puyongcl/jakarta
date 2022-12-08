package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/userkey"
)

var _ UserLoginStateModel = (*customUserLoginStateModel)(nil)

type (
	// UserLoginStateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserLoginStateModel.
	UserLoginStateModel interface {
		userLoginStateModel
		InsertTrans(ctx context.Context, session sqlx.Session, data *UserLoginState) (sql.Result, error)
		CountUserOnline(ctx context.Context) (int64, error)
		UpdateLoginState(ctx context.Context, data *UserLoginState) (int64, error)
	}

	customUserLoginStateModel struct {
		*defaultUserLoginStateModel
	}
)

// NewUserLoginStateModel returns a model for the database table.
func NewUserLoginStateModel(conn sqlx.SqlConn, c cache.CacheConf) UserLoginStateModel {
	return &customUserLoginStateModel{
		defaultUserLoginStateModel: newUserLoginStateModel(conn, c),
	}
}

func (m *defaultUserLoginStateModel) InsertTrans(ctx context.Context, session sqlx.Session, data *UserLoginState) (sql.Result, error) {
	jakartaUserLoginStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserLoginStateUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7)", m.table, userLoginStateRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Uid, data.LoginTime, data.OfflineTime, data.LoginState, data.LoginCntSum, data.LoginCntToday, data.ImEventTime)
	}, jakartaUserLoginStateUidKey)
	return ret, err
}

func (m *defaultUserLoginStateModel) CountUserOnline(ctx context.Context) (int64, error) {
	rb := squirrel.Select("count(uid)").From(m.table).Where("login_state=$1", userkey.Login)
	query, args, err := rb.ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultUserLoginStateModel) UpdateLoginState(ctx context.Context, data *UserLoginState) (int64, error) {
	jakartaUserLoginStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserLoginStateUidPrefix, data.Uid)
	rs, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if data.LoginState == userkey.Login {
			query := fmt.Sprintf("update %s set login_time=$1,login_state=$2,login_cnt_sum=$3,login_cnt_today=$4,im_event_time=$5 where uid=$6 and im_event_time<$7", m.table)
			return conn.ExecCtx(ctx, query, data.LoginTime, data.LoginState, data.LoginCntSum, data.LoginCntToday, data.ImEventTime, data.Uid, data.ImEventTime)
		}
		query := fmt.Sprintf("update %s set offline_time=$1, login_state=$2,im_event_time=$3 where uid=$4 and im_event_time<$5", m.table)
		return conn.ExecCtx(ctx, query, data.OfflineTime, data.LoginState, data.ImEventTime, data.Uid, data.ImEventTime)
	}, jakartaUserLoginStateUidKey)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}
