package chatPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ VoiceChatLogModel = (*customVoiceChatLogModel)(nil)

type (
	// VoiceChatLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customVoiceChatLogModel.
	VoiceChatLogModel interface {
		voiceChatLogModel
		Insert2(ctx context.Context, session sqlx.Session, data *VoiceChatLog) error
		UpdateChatLog(ctx context.Context, session sqlx.Session, id string, action, state int64, endTime *time.Time) error
		Find(ctx context.Context, uid, listenerUid int64, state int64) ([]*VoiceChatLog, error)
		UpdateChatLogState(ctx context.Context, session sqlx.Session, id string, state int64) error
	}

	customVoiceChatLogModel struct {
		*defaultVoiceChatLogModel
	}
)

// NewVoiceChatLogModel returns a model for the database table.
func NewVoiceChatLogModel(conn sqlx.SqlConn, c cache.CacheConf) VoiceChatLogModel {
	return &customVoiceChatLogModel{
		defaultVoiceChatLogModel: newVoiceChatLogModel(conn),
	}
}

func (m *defaultVoiceChatLogModel) UpdateChatLog(ctx context.Context, session sqlx.Session, id string, action, state int64, endTime *time.Time) error {
	query := fmt.Sprintf("update %s set end_time=$1, stop_action = $2, state = $3 where id = $4", m.table)
	if session != nil {
		_, err := session.ExecCtx(ctx, query, endTime, action, state, id)
		return err
	}
	_, err := m.conn.ExecCtx(ctx, query, endTime, action, state, id)
	return err
}

func (m *defaultVoiceChatLogModel) Find(ctx context.Context, uid, listenerUid int64, state int64) ([]*VoiceChatLog, error) {
	query := fmt.Sprintf("select %s from %s where uid = $1 and listener_uid = $2 and state = $3", voiceChatLogRows, m.table)
	resp := make([]*VoiceChatLog, 0)
	err := m.conn.QueryRowsCtx(ctx, &resp, query, uid, listenerUid, state)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultVoiceChatLogModel) UpdateChatLogState(ctx context.Context, session sqlx.Session, id string, state int64) error {
	query := fmt.Sprintf("update %s set state = $1 where id = $2", m.table)
	if session != nil {
		_, err := session.ExecCtx(ctx, query, state, id)
		return err
	}
	_, err := m.conn.ExecCtx(ctx, query, state, id)
	return err
}

func (m *defaultVoiceChatLogModel) Insert2(ctx context.Context, session sqlx.Session, data *VoiceChatLog) error {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)", m.table, voiceChatLogRowsExpectAutoSet)
	_, err := session.ExecCtx(ctx, query, data.Id, data.ListenerUid, data.Uid, data.EndTime, data.StartTime, data.StartAction, data.StopAction, data.State)
	return err
}
