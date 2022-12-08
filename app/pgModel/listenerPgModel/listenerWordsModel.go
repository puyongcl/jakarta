package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/tool"
)

var _ ListenerWordsModel = (*customListenerWordsModel)(nil)

type (
	// ListenerWordsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerWordsModel.
	ListenerWordsModel interface {
		listenerWordsModel
		InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerWords) (sql.Result, error)
		UpdatePart(ctx context.Context, data *ListenerWords) error
	}

	customListenerWordsModel struct {
		*defaultListenerWordsModel
	}
)

// NewListenerWordsModel returns a model for the database table.
func NewListenerWordsModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerWordsModel {
	return &customListenerWordsModel{
		defaultListenerWordsModel: newListenerWordsModel(conn, c),
	}
}

func (m *defaultListenerWordsModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerWords) (sql.Result, error) {
	jakartaListenerWordsListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, listenerWordsRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.ListenerUid, data.Words1, data.Words2, data.Words3, data.Words4, data.Words5, data.Words6, data.Words7, data.Words8, data.Words9, data.Words10, data.WordsSort)
	}, jakartaListenerWordsListenerUidKey)
	return ret, err
}

func (m *defaultListenerWordsModel) UpdatePart(ctx context.Context, newData *ListenerWords) error {
	data, err := m.FindOne(ctx, newData.ListenerUid)
	if err != nil {
		return err
	}
	jakartaListenerWordsListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWordsListenerUidPrefix, newData.ListenerUid)
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.Words1 != data.Words1 {
		rb = rb.Set("words_1", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words1))
		argNo++
	}
	if newData.Words2 != data.Words2 {
		rb = rb.Set("words_2", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words2))
		argNo++
	}
	if newData.Words3 != data.Words3 {
		rb = rb.Set("words_3", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words3))
		argNo++
	}
	if newData.Words4 != data.Words4 {
		rb = rb.Set("words_4", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words4))
		argNo++
	}
	if newData.Words5 != data.Words5 {
		rb = rb.Set("words_5", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words5))
		argNo++
	}
	if newData.Words6 != data.Words6 {
		rb = rb.Set("words_6", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words6))
		argNo++
	}
	if newData.Words7 != data.Words7 {
		rb = rb.Set("words_7", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words7))
		argNo++
	}
	if newData.Words8 != data.Words8 {
		rb = rb.Set("words_8", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words8))
		argNo++
	}
	if newData.Words9 != data.Words9 {
		rb = rb.Set("words_9", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words9))
		argNo++
	}
	if newData.Words10 != data.Words10 {
		rb = rb.Set("words_10", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Words10))
		argNo++
	}
	if !tool.IsEqualArrayInt64Order(newData.WordsSort, data.WordsSort) {
		rb = rb.Set("words_sort", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.WordsSort))
		argNo++
	}
	query, args, err := rb.Where(fmt.Sprintf("listener_uid=$%d", argNo), newData.ListenerUid).ToSql()
	if err != nil {
		return err
	}
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerWordsListenerUidKey)
	return err
}
