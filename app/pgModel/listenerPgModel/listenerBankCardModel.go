package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ListenerBankCardModel = (*customListenerBankCardModel)(nil)

type (
	// ListenerBankCardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerBankCardModel.
	ListenerBankCardModel interface {
		listenerBankCardModel
		UpdateListenerBankCard(ctx context.Context, newData *ListenerBankCard) error
	}

	customListenerBankCardModel struct {
		*defaultListenerBankCardModel
	}
)

// NewListenerBankCardModel returns a model for the database table.
func NewListenerBankCardModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerBankCardModel {
	return &customListenerBankCardModel{
		defaultListenerBankCardModel: newListenerBankCardModel(conn, c),
	}
}

//
func (m *defaultListenerBankCardModel) UpdateListenerBankCard(ctx context.Context, newData *ListenerBankCard) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.ListenerName != "" {
		rb = rb.Set("listener_name", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ListenerName))
		argNo++
	}
	if newData.PhoneNumber != "" {
		rb = rb.Set("phone_number", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.PhoneNumber))
		argNo++
	}
	if newData.IdNo != "" {
		rb = rb.Set("id_no", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.IdNo))
		argNo++
	}
	if newData.BankCardNo != "" {
		rb = rb.Set("bank_card_no", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.BankCardNo))
		argNo++
	}
	query, args, err := rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), newData.ListenerUid).ToSql()
	if err != nil {
		return err
	}
	jakartaListenerBankCardListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerBankCardListenerUidPrefix, newData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerBankCardListenerUidKey)
	return err
}
