package chatPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/chatkey"
	"time"
)

var _ ListenerVoiceChatStateModel = (*customListenerVoiceChatStateModel)(nil)

type (
	// ListenerVoiceChatStateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerVoiceChatStateModel.
	ListenerVoiceChatStateModel interface {
		listenerVoiceChatStateModel
		UpdateState(ctx context.Context, session sqlx.Session, uid, listenerUid, state int64) error
	}

	customListenerVoiceChatStateModel struct {
		*defaultListenerVoiceChatStateModel
	}
)

// NewListenerVoiceChatStateModel returns a model for the database table.
func NewListenerVoiceChatStateModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerVoiceChatStateModel {
	return &customListenerVoiceChatStateModel{
		defaultListenerVoiceChatStateModel: newListenerVoiceChatStateModel(conn, c),
	}
}

func (m *defaultListenerVoiceChatStateModel) UpdateState(ctx context.Context, session sqlx.Session, uid, listenerUid, state int64) error {
	rb := squirrel.Update(m.table)
	switch state {
	case chatkey.VoiceChatStateStart:
		rb = rb.Set("uid", squirrel.Expr("$1", uid)).Set("state", squirrel.Expr("$2", state)).Set("start_time", squirrel.Expr("$3", time.Now())).Where("listener_uid = $4", listenerUid)
	case chatkey.VoiceChatStateStop:
		rb = rb.Set("state", squirrel.Expr("$1", state)).Set("end_time", squirrel.Expr("$2", time.Now())).Where("listener_uid = $3 and uid = $4", listenerUid, uid)
	case chatkey.VoiceChatStateSettle:
		rb = rb.Set("state", squirrel.Expr("$1", state)).Set("settle_time", squirrel.Expr("$2", time.Now())).Where("listener_uid = $3 and uid = $4", listenerUid, uid)
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}

	jakartaListenerVoiceChatStateListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerVoiceChatStateListenerUidPrefix, listenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerVoiceChatStateListenerUidKey)
	return err
}
