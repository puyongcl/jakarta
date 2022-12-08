package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserWechatInfoModel = (*customUserWechatInfoModel)(nil)

type (
	// UserWechatInfoModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserWechatInfoModel.
	UserWechatInfoModel interface {
		userWechatInfoModel
		InsertOrUpdateMPTrans(ctx context.Context, session sqlx.Session, data *UserWechatInfo) (sql.Result, error)
		InsertOrUpdateFwhTrans(ctx context.Context, data *UserWechatInfo) (sql.Result, error)
		UpdateFwh(ctx context.Context, uid int64, unionId string, fwhOpenId string, fwhState int64) error
		UpdateFwhUnsubscribe(ctx context.Context, fwhOpenId string, fwhState int64) error
		FindOneByFwhOpenId(ctx context.Context, fwhOpenId string) (*UserWechatInfo, error)
	}

	customUserWechatInfoModel struct {
		*defaultUserWechatInfoModel
	}
)

// NewUserWechatInfoModel returns a model for the database table.
func NewUserWechatInfoModel(conn sqlx.SqlConn, c cache.CacheConf) UserWechatInfoModel {
	return &customUserWechatInfoModel{
		defaultUserWechatInfoModel: newUserWechatInfoModel(conn, c),
	}
}

func (m *defaultUserWechatInfoModel) InsertOrUpdateMPTrans(ctx context.Context, session sqlx.Session, data *UserWechatInfo) (sql.Result, error) {
	jakartaUserWechatInfoUnionIdKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUnionIdPrefix, data.UnionId)
	jakartaUserWechatInfoUidKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5) ON CONFLICT (union_id) DO UPDATE SET mp_openid = $6, uid = $7", m.table, userWechatInfoRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.Uid, data.MpOpenid, data.FwhOpenid, data.UnionId, data.FwhState, data.MpOpenid, data.Uid)
		}
		return conn.ExecCtx(ctx, query, data.Uid, data.MpOpenid, data.FwhOpenid, data.UnionId, data.FwhState, data.MpOpenid, data.Uid)
	}, jakartaUserWechatInfoUnionIdKey, jakartaUserWechatInfoUidKey)
	return ret, err
}

func (m *defaultUserWechatInfoModel) InsertOrUpdateFwhTrans(ctx context.Context, data *UserWechatInfo) (sql.Result, error) {
	jakartaUserWechatInfoUnionIdKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUnionIdPrefix, data.UnionId)
	jakartaUserWechatInfoUidKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5) ON CONFLICT (union_id) DO UPDATE SET fwh_openid = $6, fwh_state = $7", m.table, userWechatInfoRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Uid, data.MpOpenid, data.FwhOpenid, data.UnionId, data.FwhState, data.FwhOpenid, data.FwhState)
	}, jakartaUserWechatInfoUnionIdKey, jakartaUserWechatInfoUidKey)
	return ret, err
}

func (m *defaultUserWechatInfoModel) UpdateFwh(ctx context.Context, uid int64, unionId string, fwhOpenId string, fwhState int64) error {
	var err error
	if uid == 0 {
		var data *UserWechatInfo
		data, err = m.FindOne(ctx, unionId)
		if err != nil {
			return err
		}
		uid = data.Uid
	}
	rb := squirrel.Update(m.table)
	argNo := 1
	if fwhOpenId != "" {
		rb = rb.Set("fwh_openid", squirrel.Expr(fmt.Sprintf("$%d", argNo), fwhOpenId))
		argNo++
	}
	if fwhState != 0 {
		rb = rb.Set("fwh_state", squirrel.Expr(fmt.Sprintf("$%d", argNo), fwhState))
		argNo++
	}
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaUserWechatInfoUidKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUidPrefix, uid)
	jakartaUserWechatInfoUnionIdKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUnionIdPrefix, unionId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaUserWechatInfoUidKey, jakartaUserWechatInfoUnionIdKey)
	return err
}

func (m *defaultUserWechatInfoModel) UpdateFwhUnsubscribe(ctx context.Context, fwhOpenId string, fwhState int64) error {
	var data *UserWechatInfo
	var err error
	data, err = m.FindOneByFwhOpenId(ctx, fwhOpenId)
	if err != nil && err != ErrNotFound {
		return err
	}

	if data == nil {
		return nil
	}

	query := fmt.Sprintf("update %s set fwh_state = $1 where union_id = $2", m.table)
	jakartaUserWechatInfoUidKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUidPrefix, data.Uid)
	jakartaUserWechatInfoUnionIdKey := fmt.Sprintf("%s%v", cacheJakartaUserWechatInfoUnionIdPrefix, data.UnionId)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, fwhState, data.UnionId)
	}, jakartaUserWechatInfoUidKey, jakartaUserWechatInfoUnionIdKey)
	return err
}

func (m *defaultUserWechatInfoModel) FindOneByFwhOpenId(ctx context.Context, fwhOpenId string) (*UserWechatInfo, error) {
	var resp UserWechatInfo
	query := fmt.Sprintf("select %s from %s where fwh_openid = $1 limit 1", userWechatInfoRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, fwhOpenId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
