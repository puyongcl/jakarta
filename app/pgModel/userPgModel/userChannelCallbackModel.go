package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserChannelCallbackModel = (*customUserChannelCallbackModel)(nil)

type (
	// UserChannelCallbackModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserChannelCallbackModel.
	UserChannelCallbackModel interface {
		userChannelCallbackModel
		InsertTrans(ctx context.Context, session sqlx.Session, data *UserChannelCallback) (sql.Result, error)
	}

	customUserChannelCallbackModel struct {
		*defaultUserChannelCallbackModel
	}
)

// NewUserChannelCallbackModel returns a model for the database table.
func NewUserChannelCallbackModel(conn sqlx.SqlConn, c cache.CacheConf) UserChannelCallbackModel {
	return &customUserChannelCallbackModel{
		defaultUserChannelCallbackModel: newUserChannelCallbackModel(conn, c),
	}
}

func (m *defaultUserChannelCallbackModel) InsertTrans(ctx context.Context, session sqlx.Session, data *UserChannelCallback) (sql.Result, error) {
	jakartaUserChannelCallbackUidKey := fmt.Sprintf("%s%v", cacheJakartaUserChannelCallbackUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3)", m.table, userChannelCallbackRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Uid, data.Channel, data.Cb)
	}, jakartaUserChannelCallbackUidKey)
	return ret, err
}
