package bbsPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StoryReplyLikeLogModel = (*customStoryReplyLikeLogModel)(nil)

type (
	// StoryReplyLikeLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStoryReplyLikeLogModel.
	StoryReplyLikeLogModel interface {
		storyReplyLikeLogModel
		AddStoryReplyLikeLogLikeCnt(ctx context.Context, id string) error
	}

	customStoryReplyLikeLogModel struct {
		*defaultStoryReplyLikeLogModel
	}
)

// NewStoryReplyLikeLogModel returns a model for the database table.
func NewStoryReplyLikeLogModel(conn sqlx.SqlConn, c cache.CacheConf) StoryReplyLikeLogModel {
	return &customStoryReplyLikeLogModel{
		defaultStoryReplyLikeLogModel: newStoryReplyLikeLogModel(conn, c),
	}
}

func (m *defaultStoryReplyLikeLogModel) AddStoryReplyLikeLogLikeCnt(ctx context.Context, id string) error {
	jakartaStoryReplyLikeLogIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyLikeLogIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set like_cnt = like_cnt + 1 where id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, jakartaStoryReplyLikeLogIdKey)
	return err
}
