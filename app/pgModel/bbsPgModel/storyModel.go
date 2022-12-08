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
)

var _ StoryModel = (*customStoryModel)(nil)

type (
	// StoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStoryModel.
	StoryModel interface {
		storyModel
		UpdateStoryState(ctx context.Context, id string, state int64) error
		AddStoryReplyCnt(ctx context.Context, id string) error
		AddStoryViewCnt(ctx context.Context, id string) error
		Find(ctx context.Context, uid int64, storyType int64, pageNo, pageSize int64) ([]*Story, error)
		FindRec(ctx context.Context, storyType, spec int64, pageNo, pageSize int64, blk []int64) ([]*Story, error)
		CountByTextMd5(ctx context.Context, textMd5 string) (int64, error)
	}

	customStoryModel struct {
		*defaultStoryModel
	}
)

// NewStoryModel returns a model for the database table.
func NewStoryModel(conn sqlx.SqlConn, c cache.CacheConf) StoryModel {
	return &customStoryModel{
		defaultStoryModel: newStoryModel(conn, c),
	}
}

func (m *defaultStoryModel) UpdateStoryState(ctx context.Context, id string, state int64) error {
	jakartaStoryIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set state = $1 where id = $2", m.table)
		return conn.ExecCtx(ctx, query, state, id)
	}, jakartaStoryIdKey)
	return err
}

func (m *defaultStoryModel) AddStoryReplyCnt(ctx context.Context, id string) error {
	jakartaStoryIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set reply_cnt = reply_cnt + 1 where id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, jakartaStoryIdKey)
	return err
}

func (m *defaultStoryModel) AddStoryViewCnt(ctx context.Context, id string) error {
	jakartaStoryIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set view_cnt = view_cnt + 1 where id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, jakartaStoryIdKey)
	return err
}

func (m *defaultStoryModel) Find(ctx context.Context, uid int64, storyType int64, pageNo, pageSize int64) ([]*Story, error) {
	rb := squirrel.Select(storyRows).From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if storyType != 0 {
		rb = rb.Where(fmt.Sprintf("story_type = $%d", argNo), storyType)
		argNo++
	}

	rb = rb.Where(fmt.Sprintf("state != $%d", argNo), bbskey.StoryStateDeleted)
	argNo++

	q, args, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*Story, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, q, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultStoryModel) FindRec(ctx context.Context, storyType, spec int64, pageNo, pageSize int64, blk []int64) ([]*Story, error) {
	rb := squirrel.Select(storyRows).From(m.table)
	argNo := 1
	if spec != 0 {
		rb = rb.Where(fmt.Sprintf("spec = $%d", argNo), spec)
		argNo++
	}
	if storyType != 0 {
		rb = rb.Where(fmt.Sprintf("story_type = $%d", argNo), storyType)
		argNo++
	}

	if len(blk) > 0 {
		rb = rb.Where(fmt.Sprintf("uid != ALL($%d)", argNo), pq.Int64Array(blk))
		argNo++
	}
	rb = rb.Where(fmt.Sprintf("state != $%d", argNo), bbskey.StoryStateDeleted)
	argNo++

	q, args, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*Story, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, q, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultStoryModel) CountByTextMd5(ctx context.Context, textMd5 string) (int64, error) {
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
