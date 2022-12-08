package chatPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/chatkey"
	"jakarta/common/xerr"
	"time"
)

var _ UserVoiceChatStateModel = (*customUserVoiceChatStateModel)(nil)

type (
	// UserVoiceChatStateModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserVoiceChatStateModel.
	UserVoiceChatStateModel interface {
		userVoiceChatStateModel
		UpdateState(ctx context.Context, session sqlx.Session, uid, listenerUid, state int64) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *UserVoiceChatState) (sql.Result, error)
	}

	customUserVoiceChatStateModel struct {
		*defaultUserVoiceChatStateModel
	}
)

// NewUserVoiceChatStateModel returns a model for the database table.
func NewUserVoiceChatStateModel(conn sqlx.SqlConn, c cache.CacheConf) UserVoiceChatStateModel {
	return &customUserVoiceChatStateModel{
		defaultUserVoiceChatStateModel: newUserVoiceChatStateModel(conn, c),
	}
}

func (m *defaultUserVoiceChatStateModel) UpdateState(ctx context.Context, session sqlx.Session, uid, listenerUid, state int64) error {
	rb := squirrel.Update(m.table)
	switch state {
	case chatkey.VoiceChatStateStart:
		rb = rb.Set("listener_uid", squirrel.Expr("$1", listenerUid)).Set("state", squirrel.Expr("$2", state)).Set("start_time", squirrel.Expr("$3", time.Now())).Where("uid = $4", uid)
	case chatkey.VoiceChatStateStop:
		rb = rb.Set("state", squirrel.Expr("$1", state)).Set("end_time", squirrel.Expr("$2", time.Now())).Where("listener_uid = $3 and uid = $4", listenerUid, uid)
	case chatkey.VoiceChatStateSettle:
		rb = rb.Set("state", squirrel.Expr("$1", state)).Set("settle_time", squirrel.Expr("$2", time.Now())).Where("listener_uid = $3 and uid = $4", listenerUid, uid)
	default:
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("状态错误 %d", state))
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}

	jakartaUserVoiceChatStateListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaUserVoiceChatStateListenerUidKey)
	return err
}

func (m *defaultUserVoiceChatStateModel) InsertTrans(ctx context.Context, session sqlx.Session, data *UserVoiceChatState) (sql.Result, error) {
	jakartaUserVoiceChatStateUidKey := fmt.Sprintf("%s%v", cacheJakartaUserVoiceChatStateUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6)", m.table, userVoiceChatStateRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Uid, data.ListenerUid, data.State, data.StartTime, data.EndTime, data.SettleTime)
	}, jakartaUserVoiceChatStateUidKey)
	return ret, err
}
