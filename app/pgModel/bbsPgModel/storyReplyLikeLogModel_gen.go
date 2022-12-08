// Code generated by goctl. DO NOT EDIT!

package bbsPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	storyReplyLikeLogFieldNames          = builder.RawFieldNames(&StoryReplyLikeLog{}, true)
	storyReplyLikeLogRows                = strings.Join(storyReplyLikeLogFieldNames, ",")
	storyReplyLikeLogRowsExpectAutoSet   = strings.Join(stringx.Remove(storyReplyLikeLogFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	storyReplyLikeLogRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(storyReplyLikeLogFieldNames, "id", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaStoryReplyLikeLogIdPrefix = "cache:jakarta:storyReplyLikeLog:id:"
)

type (
	storyReplyLikeLogModel interface {
		Insert(ctx context.Context, data *StoryReplyLikeLog) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*StoryReplyLikeLog, error)
		Update(ctx context.Context, data *StoryReplyLikeLog) error
		Delete(ctx context.Context, id string) error
	}

	defaultStoryReplyLikeLogModel struct {
		sqlc.CachedConn
		table string
	}

	StoryReplyLikeLog struct {
		CreateTime   time.Time `db:"create_time"`
		UpdateTime   time.Time `db:"update_time"`
		Id           string    `db:"id"`
		LikeCnt      int64     `db:"like_cnt"`
		Uid          int64     `db:"uid"`
		StoryReplyId string    `db:"story_reply_id"`
	}
)

func newStoryReplyLikeLogModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultStoryReplyLikeLogModel {
	return &defaultStoryReplyLikeLogModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."story_reply_like_log"`,
	}
}

func (m *defaultStoryReplyLikeLogModel) Delete(ctx context.Context, id string) error {
	jakartaStoryReplyLikeLogIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyLikeLogIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = $1", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, jakartaStoryReplyLikeLogIdKey)
	return err
}

func (m *defaultStoryReplyLikeLogModel) FindOne(ctx context.Context, id string) (*StoryReplyLikeLog, error) {
	jakartaStoryReplyLikeLogIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyLikeLogIdPrefix, id)
	var resp StoryReplyLikeLog
	err := m.QueryRowCtx(ctx, &resp, jakartaStoryReplyLikeLogIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = $1 limit 1", storyReplyLikeLogRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultStoryReplyLikeLogModel) Insert(ctx context.Context, data *StoryReplyLikeLog) (sql.Result, error) {
	jakartaStoryReplyLikeLogIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyLikeLogIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4)", m.table, storyReplyLikeLogRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Id, data.LikeCnt, data.Uid, data.StoryReplyId)
	}, jakartaStoryReplyLikeLogIdKey)
	return ret, err
}

func (m *defaultStoryReplyLikeLogModel) Update(ctx context.Context, data *StoryReplyLikeLog) error {
	jakartaStoryReplyLikeLogIdKey := fmt.Sprintf("%s%v", cacheJakartaStoryReplyLikeLogIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = $1", m.table, storyReplyLikeLogRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.Id, data.LikeCnt, data.Uid, data.StoryReplyId)
	}, jakartaStoryReplyLikeLogIdKey)
	return err
}

func (m *defaultStoryReplyLikeLogModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaStoryReplyLikeLogIdPrefix, primary)
}

func (m *defaultStoryReplyLikeLogModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", storyReplyLikeLogRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultStoryReplyLikeLogModel) tableName() string {
	return m.table
}
