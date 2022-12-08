package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserStatModel = (*customUserStatModel)(nil)

type (
	// UserStatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserStatModel.
	UserStatModel interface {
		userStatModel
		InsertTrans(ctx context.Context, session sqlx.Session, data *UserStat) (sql.Result, error)
		UpdateUserStat(ctx context.Context, data *UpdateUserStatData) error
		CanNoCondRefund(ctx context.Context, uid int64) (bool, error)
	}

	customUserStatModel struct {
		*defaultUserStatModel
	}
)

type UpdateUserStatData struct {
	Uid                int64 `json:"uid"`
	AddCostAmountSum   int64 `json:"addCostAmountSum"`   // 支付成功总额
	AddRefundAmountSum int64 `json:"addRefundAmountSum"` // 退款总额
	AddPaidOrderCnt    int64 `json:"addPaidOrderCnt"`    // 支付成功订单数量
	AddRefundOrderCnt  int64 `json:"addRefundOrderCnt"`  // 退款订单数
}

// NewUserStatModel returns a model for the database table.
func NewUserStatModel(conn sqlx.SqlConn, c cache.CacheConf) UserStatModel {
	return &customUserStatModel{
		defaultUserStatModel: newUserStatModel(conn, c),
	}
}

func (m *defaultUserStatModel) InsertTrans(ctx context.Context, session sqlx.Session, data *UserStat) (sql.Result, error) {
	jakartaUserStatUidKey := fmt.Sprintf("%s%v", cacheJakartaUserStatUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6)", m.table, userStatRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Uid, data.CostAmountSum, data.RefundAmountSum, data.PaidOrderCnt, data.RefundOrderCnt, data.NoCondRefundCnt)
	}, jakartaUserStatUidKey)
	return ret, err
}

func (m *defaultUserStatModel) UpdateUserStat(ctx context.Context, addData *UpdateUserStatData) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if addData.AddCostAmountSum != 0 {
		rb = rb.Set("cost_amount_sum", squirrel.Expr(fmt.Sprintf("cost_amount_sum + $%d", argNo), addData.AddCostAmountSum))
		argNo++
	}
	if addData.AddRefundAmountSum != 0 {
		rb = rb.Set("refund_amount_sum", squirrel.Expr(fmt.Sprintf("refund_amount_sum + $%d", argNo), addData.AddRefundAmountSum))
		argNo++
	}
	if addData.AddPaidOrderCnt != 0 {
		rb = rb.Set("paid_order_cnt", squirrel.Expr(fmt.Sprintf("paid_order_cnt + $%d", argNo), addData.AddPaidOrderCnt))
		argNo++
	}
	if addData.AddRefundOrderCnt != 0 {
		rb = rb.Set("refund_order_cnt", squirrel.Expr(fmt.Sprintf("refund_order_cnt + $%d", argNo), addData.AddRefundOrderCnt))
		argNo++
	}

	rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), addData.Uid)

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaUserStatUidKey := fmt.Sprintf("%s%v", cacheJakartaUserStatUidPrefix, addData.Uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaUserStatUidKey)
	return err
}

const canNoCondRefundQuery = "update %s set no_cond_refund_cnt = no_cond_refund_cnt - 1 where uid = $1"

func (m *defaultUserStatModel) CanNoCondRefund(ctx context.Context, uid int64) (bool, error) {
	data, err := m.FindOne(ctx, uid)
	if err != nil {
		return false, err
	}
	if data.NoCondRefundCnt <= 0 {
		return false, nil
	}
	query := fmt.Sprintf(canNoCondRefundQuery, m.table)
	jakartaUserStatUidKey := fmt.Sprintf("%s%v", cacheJakartaUserStatUidPrefix, uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, uid)
	}, jakartaUserStatUidKey)
	if err != nil {
		return false, err
	}
	return true, nil
}
