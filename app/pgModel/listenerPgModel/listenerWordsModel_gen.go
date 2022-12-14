// Code generated by goctl. DO NOT EDIT!

package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	listenerWordsFieldNames          = builder.RawFieldNames(&ListenerWords{}, true)
	listenerWordsRows                = strings.Join(listenerWordsFieldNames, ",")
	listenerWordsRowsExpectAutoSet   = strings.Join(stringx.Remove(listenerWordsFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	listenerWordsRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(listenerWordsFieldNames, "listener_uid", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaListenerWordsListenerUidPrefix = "cache:jakarta:listenerWords:listenerUid:"
)

type (
	listenerWordsModel interface {
		Insert(ctx context.Context, data *ListenerWords) (sql.Result, error)
		FindOne(ctx context.Context, listenerUid int64) (*ListenerWords, error)
		Update(ctx context.Context, data *ListenerWords) error
		Delete(ctx context.Context, listenerUid int64) error
	}

	defaultListenerWordsModel struct {
		sqlc.CachedConn
		table string
	}

	ListenerWords struct {
		CreateTime  time.Time     `db:"create_time"`
		UpdateTime  time.Time     `db:"update_time"`
		ListenerUid int64         `db:"listener_uid"`
		Words1      string        `db:"words_1"`
		Words2      string        `db:"words_2"`
		Words3      string        `db:"words_3"`
		Words4      string        `db:"words_4"`
		Words5      string        `db:"words_5"`
		Words6      string        `db:"words_6"`
		Words7      string        `db:"words_7"`
		Words8      string        `db:"words_8"`
		Words9      string        `db:"words_9"`
		Words10     string        `db:"words_10"`
		WordsSort   pq.Int64Array `db:"words_sort"`
	}
)

func newListenerWordsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultListenerWordsModel {
	return &defaultListenerWordsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."listener_words"`,
	}
}

func (m *defaultListenerWordsModel) Delete(ctx context.Context, listenerUid int64) error {
	jakartaListenerWordsListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where listener_uid = $1", m.table)
		return conn.ExecCtx(ctx, query, listenerUid)
	}, jakartaListenerWordsListenerUidKey)
	return err
}

func (m *defaultListenerWordsModel) FindOne(ctx context.Context, listenerUid int64) (*ListenerWords, error) {
	jakartaListenerWordsListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, listenerUid)
	var resp ListenerWords
	err := m.QueryRowCtx(ctx, &resp, jakartaListenerWordsListenerUidKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where listener_uid = $1 limit 1", listenerWordsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, listenerUid)
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

func (m *defaultListenerWordsModel) Insert(ctx context.Context, data *ListenerWords) (sql.Result, error) {
	jakartaListenerWordsListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, listenerWordsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.ListenerUid, data.Words1, data.Words2, data.Words3, data.Words4, data.Words5, data.Words6, data.Words7, data.Words8, data.Words9, data.Words10, data.WordsSort)
	}, jakartaListenerWordsListenerUidKey)
	return ret, err
}

func (m *defaultListenerWordsModel) Update(ctx context.Context, data *ListenerWords) error {
	jakartaListenerWordsListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, data.ListenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where listener_uid = $1", m.table, listenerWordsRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.ListenerUid, data.Words1, data.Words2, data.Words3, data.Words4, data.Words5, data.Words6, data.Words7, data.Words8, data.Words9, data.Words10, data.WordsSort)
	}, jakartaListenerWordsListenerUidKey)
	return err
}

func (m *defaultListenerWordsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, primary)
}

func (m *defaultListenerWordsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where listener_uid = $1 limit 1", listenerWordsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultListenerWordsModel) tableName() string {
	return m.table
}
