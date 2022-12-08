package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/xerr"
	"time"
)

var _ ListenerWalletModel = (*customListenerWalletModel)(nil)

type (
	// ListenerWalletModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerWalletModel.
	ListenerWalletModel interface {
		listenerWalletModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		AddOrderAmount(ctx context.Context, listenerUid int64, amount int64) error
		AddConfirmAmount(ctx context.Context, listenerUid int64, amount int64, currentAddMount int64) error
		RefundAmount(ctx context.Context, listenerUid int64, amount int64) error
		ApplyCashAmount(ctx context.Context, listenerUid int64, amount int64) error
		AlreadyCashAmount(ctx context.Context, listenerUid int64, amount int64) error
		CashFail(ctx context.Context, listenerUid int64, amount int64) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerWallet) (sql.Result, error)
		FixConfirmAmount(ctx context.Context, listenerUid int64, amount int64, outTime *time.Time) error
		// TODO 对外获取钱包请使用此方法，因为需要判断是否重置本月统计数据
		FindOne2(ctx context.Context, listenerUid int64) (*ListenerWallet, error)
	}

	customListenerWalletModel struct {
		*defaultListenerWalletModel
	}
)

// NewListenerWalletModel returns a model for the database table.
func NewListenerWalletModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerWalletModel {
	return &customListenerWalletModel{
		defaultListenerWalletModel: newListenerWalletModel(conn, c),
	}
}

func (m *defaultListenerWalletModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultListenerWalletModel) AddOrderAmount(ctx context.Context, listenerUid int64, amount int64) error {
	data, err := m.FindOne(ctx, listenerUid)
	if err != nil {
		return err
	}
	//判断是否重置本月统计
	var query string
	now := time.Now()
	args := make([]interface{}, 0)
	if data.ResetStatTime.Month() != now.Month() { // 需要重置
		query = fmt.Sprintf("update %s set current_month_order_amount = $1, current_month_amount = 0, reset_stat_time = $2 where listener_uid = $3", m.table)
		args = append(args, amount, now, listenerUid)
	} else {
		query = fmt.Sprintf("update %s set current_month_order_amount = current_month_order_amount + $1 where listener_uid = $2", m.table)
		args = append(args, amount, listenerUid)
	}
	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerWalletListenerUidKey)
	return err
}

func (m *defaultListenerWalletModel) AddConfirmAmount(ctx context.Context, listenerUid int64, amount int64, currentAddMount int64) error {
	data, err := m.FindOne(ctx, listenerUid)
	if err != nil {
		return err
	}
	//判断是否重置本月统计
	var query string
	now := time.Now()
	args := make([]interface{}, 0)
	if data.ResetStatTime.Month() != now.Month() { // 需要重置
		query = fmt.Sprintf("update %s set amount = amount + $1, current_month_amount = 0, current_month_order_amount = 0, reset_stat_time = $2, last_confirm_time = $3 where listener_uid = $4", m.table)
		args = append(args, amount, now, now, listenerUid)
	} else {
		query = fmt.Sprintf("update %s set amount = amount + $1, current_month_amount = current_month_amount + $2, last_confirm_time = $3 where listener_uid = $4", m.table)
		args = append(args, amount, currentAddMount, now, listenerUid)
	}

	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerWalletListenerUidKey)
	return err
}

func (m *defaultListenerWalletModel) FixConfirmAmount(ctx context.Context, listenerUid int64, amount int64, outTime *time.Time) error {
	var query string
	args := make([]interface{}, 0)
	if outTime.Month() != time.Now().Month() { // 本月不加
		query = fmt.Sprintf("update %s set amount = amount + $1 where listener_uid = $2", m.table)
		args = append(args, amount, listenerUid)
	} else {
		query = fmt.Sprintf("update %s set amount = amount + $1, current_month_amount = current_month_amount + $2 where listener_uid = $3", m.table)
		args = append(args, amount, amount, listenerUid)
	}

	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerWalletListenerUidKey)
	return err
}

func (m *defaultListenerWalletModel) RefundAmount(ctx context.Context, listenerUid int64, amount int64) error {
	data, err := m.FindOne(ctx, listenerUid)
	if err != nil {
		return err
	}
	//判断是否重置本月统计
	var query string
	now := time.Now()
	args := make([]interface{}, 0)
	if data.ResetStatTime.Month() != now.Month() { // 需要重置
		query = fmt.Sprintf("update %s set refund_sum_amount = refund_sum_amount + $1, current_month_order_amount = 0, current_month_amount = 0, reset_stat_time = $2 where listener_uid = $3", m.table)
		args = append(args, amount, now, listenerUid)
	} else {
		query = fmt.Sprintf("update %s set refund_sum_amount = refund_sum_amount + $1 where listener_uid = $2", m.table)
		args = append(args, amount, listenerUid)
	}
	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerWalletListenerUidKey)
	return err
}

func (m *defaultListenerWalletModel) AlreadyCashAmount(ctx context.Context, listenerUid int64, amount int64) error {
	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set apply_cash_amount = apply_cash_amount - $1, cash_sum_amount = cash_sum_amount + $2 where listener_uid = $3", m.table)
		return conn.ExecCtx(ctx, query, amount, amount, listenerUid)
	}, jakartaListenerWalletListenerUidKey)
	return err
}

func (m *defaultListenerWalletModel) ApplyCashAmount(ctx context.Context, listenerUid int64, amount int64) error {
	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	rs, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set amount = amount - $1, apply_cash_amount = apply_cash_amount + $2 where listener_uid = $3 and amount >= $4", m.table)
		return conn.ExecCtx(ctx, query, amount, amount, listenerUid, amount)
	}, jakartaListenerWalletListenerUidKey)
	if err != nil {
		return err
	}
	ra, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if ra <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "提现失败")
	}
	return err
}

func (m *defaultListenerWalletModel) CashFail(ctx context.Context, listenerUid int64, amount int64) error {
	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set apply_cash_amount = apply_cash_amount - $1, amount = amount + $2 where listener_uid = $3", m.table)
		return conn.ExecCtx(ctx, query, amount, amount, listenerUid)
	}, jakartaListenerWalletListenerUidKey)
	return err
}

func (m *defaultListenerWalletModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerWallet) (sql.Result, error) {
	jakartaListenerWalletListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerWalletListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)", m.table, listenerWalletRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.ListenerUid, data.Amount, data.RefundSumAmount, data.CashSumAmount, data.ApplyCashAmount, data.LastConfirmTime, data.CurrentMonthOrderAmount, data.CurrentMonthAmount, data.ResetStatTime)
	}, jakartaListenerWalletListenerUidKey)
	return ret, err
}

func (m *defaultListenerWalletModel) FindOne2(ctx context.Context, listenerUid int64) (*ListenerWallet, error) {
	data, err := m.FindOne(ctx, listenerUid)
	if err != nil {
		return nil, err
	}
	// 判断是否需要重置本月统计 此处不能修改数据库
	if data.ResetStatTime.Month() != time.Now().Month() { // 需要重置
		data.CurrentMonthOrderAmount = 0
		data.CurrentMonthAmount = 0
	}
	return data, nil
}
