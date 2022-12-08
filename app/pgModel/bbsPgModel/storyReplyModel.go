package bbsPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/bbskey"
	"time"
)

var _ StoryReplyModel = (*customStoryReplyModel)(nil)

type (
	// StoryReplyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStoryReplyModel.
	StoryReplyModel interface {
		storyReplyModel
		AddStoryReplyLikeCnt(ctx context.Context, id string) error
		UpdateStoryReplyState(ctx context.Context, id string, state int64) error
		UpdateStoryNewReply(ctx context.Context, id string, text, voice string) error
		Find(ctx context.Context, listenerUid int64, storyId string, pageNo, pageSize int64, blk []int64) ([]*StoryReply, error)
		CountByTextMd5(ctx context.Context, textMd5 string) (int64, error)
	}

	customStoryReplyModel struct {
		*defaultStoryReplyModel
	}
)

// NewStoryReplyModel returns a model for the database table.
func NewStoryReplyModel(conn sqlx.SqlConn, c cache.CacheConf) StoryReplyModel {
	return &customStoryReplyModel{
		defaultStoryReplyModel: newStoryReplyModel(conn, c),
	}
}

func (m *defaultStoryReplyModel) AddStoryReplyLikeCnt(ctx context.Context, id string) error {
	jakartaStoryReplyIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set like_cnt = like_cnt + 1 where id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, jakartaStoryReplyIdKey)
	return err
}

func (m *defaultStoryReplyModel) UpdateStoryReplyState(ctx context.Context, id string, state int64) error {
	jakartaStoryReplyIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set state = $1 where id = $2", m.table)
		return conn.ExecCtx(ctx, query, state, id)
	}, jakartaStoryReplyIdKey)
	return err
}

func (m *defaultStoryReplyModel) Find(ctx context.Context, listenerUid int64, storyId string, pageNo, pageSize int64, blk []int64) ([]*StoryReply, error) {
	rb := squirrel.Select(storyReplyRows).From(m.table)
	argNo := 1
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}

	if storyId != "" {
		rb = rb.Where(fmt.Sprintf("story_id = $%d", argNo), storyId)
		argNo++
	}
	if len(blk) > 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid != ALL($%d)", argNo), pq.Int64Array(blk))
		argNo++
	}
	rb = rb.Where(fmt.Sprintf("state != $%d", argNo), bbskey.StoryStateDeleted)
	argNo++

	q, args, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*StoryReply, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, q, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultStoryReplyModel) UpdateStoryNewReply(ctx context.Context, id string, text, voice string) error {
	jakartaStoryReplyIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set reply_text=$1, reply_voice=$2, like_cnt=0, state=$3, create_time=$4 where id=$5", m.table)
		return conn.ExecCtx(ctx, query, text, voice, bbskey.StoryStatusNotCheck, time.Now(), id)
	}, jakartaStoryReplyIdKey)
	return err
}

func (m *defaultStoryReplyModel) CountByTextMd5(ctx context.Context, textMd5 string) (int64, error) {
	var resp int64
	query := fmt.Sprintf("select count(id) from %s where text_md5 = $1", m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, textMd5)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
