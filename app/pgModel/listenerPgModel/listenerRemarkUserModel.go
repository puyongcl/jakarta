package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ListenerRemarkUserModel = (*customListenerRemarkUserModel)(nil)

type (
	// ListenerRemarkUserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerRemarkUserModel.
	ListenerRemarkUserModel interface {
		listenerRemarkUserModel
		InsertOrUpdateUserRemark(ctx context.Context, id string, uid, listenerUid int64, remark string, userDesc string) error
	}

	customListenerRemarkUserModel struct {
		*defaultListenerRemarkUserModel
	}
)

// NewListenerRemarkUserModel returns a model for the database table.
func NewListenerRemarkUserModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerRemarkUserModel {
	return &customListenerRemarkUserModel{
		defaultListenerRemarkUserModel: newListenerRemarkUserModel(conn, c),
	}
}

func (m *defaultListenerRemarkUserModel) InsertOrUpdateUserRemark(ctx context.Context, id string, uid, listenerUid int64, remark string, userDesc string) error {
	data, err := m.FindOne(ctx, id)
	if err != nil && err != ErrNotFound {
		return err
	}
	if data != nil {
		if data.Remark == remark && data.UserDesc == userDesc {
			return nil
		}
		jakartaListenerRemarkUserIdKey := fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, id)
		_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
			query := fmt.Sprintf("update %s set remark = $1, user_desc = $2 where id = $3", m.table)
			return conn.ExecCtx(ctx, query, remark, userDesc, id)
		}, jakartaListenerRemarkUserIdKey)
		return err
	}

	jakartaListenerRemarkUserIdKey := fmt.Sprintf("%s%v", cacheJakartaListenerRemarkUserIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5)", m.table, listenerRemarkUserRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, id, listenerUid, uid, remark, userDesc)
	}, jakartaListenerRemarkUserIdKey)

	return err
}
