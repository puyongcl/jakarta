package chatPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatBalanceLogModel = (*customChatBalanceLogModel)(nil)

type (
	// ChatBalanceLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatBalanceLogModel.
	ChatBalanceLogModel interface {
		chatBalanceLogModel
		InsertTrans(ctx context.Context, session sqlx.Session, data *ChatBalanceLog) (sql.Result, error)
	}

	customChatBalanceLogModel struct {
		*defaultChatBalanceLogModel
	}
)

// NewChatBalanceLogModel returns a model for the database table.
func NewChatBalanceLogModel(conn sqlx.SqlConn) ChatBalanceLogModel {
	return &customChatBalanceLogModel{
		defaultChatBalanceLogModel: newChatBalanceLogModel(conn),
	}
}

func (m *defaultChatBalanceLogModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ChatBalanceLog) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7)", m.table, chatBalanceLogRowsExpectAutoSet)
	return session.ExecCtx(ctx, query, data.Id, data.EventType, data.EventId, data.Value, data.Uid, data.ListenerUid, data.ChatBalanceId)
}
