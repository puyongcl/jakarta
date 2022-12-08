package bbsPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StoryViewLogModel = (*customStoryViewLogModel)(nil)

type (
	// StoryViewLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStoryViewLogModel.
	StoryViewLogModel interface {
		storyViewLogModel
		AddStoryReplyViewLogLikeCnt(ctx context.Context, id string) error
	}

	customStoryViewLogModel struct {
		*defaultStoryViewLogModel
	}
)

// NewStoryViewLogModel returns a model for the database table.
func NewStoryViewLogModel(conn sqlx.SqlConn) StoryViewLogModel {
	return &customStoryViewLogModel{
		defaultStoryViewLogModel: newStoryViewLogModel(conn),
	}
}

func (m *defaultStoryViewLogModel) AddStoryReplyViewLogLikeCnt(ctx context.Context, id string) error {
	query := fmt.Sprintf("update %s set view_cnt = view_cnt + 1 where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}
