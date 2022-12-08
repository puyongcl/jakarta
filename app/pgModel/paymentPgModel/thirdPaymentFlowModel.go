package paymentPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ThirdPaymentFlowModel = (*customThirdPaymentFlowModel)(nil)

type (
	// ThirdPaymentFlowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customThirdPaymentFlowModel.
	ThirdPaymentFlowModel interface {
		thirdPaymentFlowModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *ThirdPaymentFlow) (sql.Result, error)
		FindOneByOrderId(ctx context.Context, orderId string, tradeState []int64) (*ThirdPaymentFlow, error)
		FindCount(ctx context.Context, orderId string, tradeState []int64) (int64, error)
		UpdatePart(ctx context.Context, newData *ThirdPaymentFlow) error
	}

	customThirdPaymentFlowModel struct {
		*defaultThirdPaymentFlowModel
	}
)

// NewThirdPaymentFlowModel returns a model for the database table.
func NewThirdPaymentFlowModel(conn sqlx.SqlConn) ThirdPaymentFlowModel {
	return &customThirdPaymentFlowModel{
		defaultThirdPaymentFlowModel: newThirdPaymentFlowModel(conn),
	}
}

// export logic
func (m *defaultThirdPaymentFlowModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultThirdPaymentFlowModel) FindOneByOrderId(ctx context.Context, orderId string, tradeState []int64) (*ThirdPaymentFlow, error) {
	query := fmt.Sprintf("select %s from %s where order_id = $1 and pay_status = any ($2) limit 1", thirdPaymentFlowRows, m.table)
	var resp ThirdPaymentFlow
	err := m.conn.QueryRowCtx(ctx, &resp, query, orderId, pq.Int64Array(tradeState))
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

const findCountPaymentFlowSQL = "select count(flow_no) from %s where order_id = $1 and pay_status = any ($2)"

func (m *defaultThirdPaymentFlowModel) FindCount(ctx context.Context, orderId string, tradeState []int64) (int64, error) {
	query := fmt.Sprintf(findCountPaymentFlowSQL, m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, orderId, pq.Int64Array(tradeState))
	return cnt, err
}

func (m *defaultThirdPaymentFlowModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ThirdPaymentFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, thirdPaymentFlowRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.FlowNo, data.Uid, data.PayMode, data.TradeType, data.TradeState, data.PayAmount, data.TransactionId, data.TradeStateDesc, data.OrderId, data.OrderType, data.PayStatus, data.PayTime, data.BankType, data.ActualPayAmount, data.Description)
	return ret, err
}

func (m *defaultThirdPaymentFlowModel) UpdatePart(ctx context.Context, newData *ThirdPaymentFlow) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.TradeState != "" {
		rb = rb.Set("trade_state", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TradeState))
		argNo++
	}
	if newData.TransactionId != "" {
		rb = rb.Set("transaction_id", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TransactionId))
		argNo++
	}
	if newData.TradeType != "" {
		rb = rb.Set("trade_type", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TradeType))
		argNo++
	}
	if newData.TradeStateDesc != "" {
		rb = rb.Set("trade_state_desc", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TradeStateDesc))
		argNo++
	}
	if newData.PayStatus != 0 {
		rb = rb.Set("pay_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.PayStatus))
		argNo++
	}
	if newData.BankType != "" {
		rb = rb.Set("bank_type", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.BankType))
		argNo++
	}
	if newData.ActualPayAmount != 0 {
		rb = rb.Set("actual_pay_amount", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ActualPayAmount))
		argNo++
	}
	if newData.PayTime.Valid {
		rb = rb.Set("pay_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.PayTime.Time))
		argNo++
	}

	query, args, err := rb.Where(fmt.Sprintf("flow_no = $%d", argNo), newData.FlowNo).ToSql()
	if err != nil {
		return err
	}

	_, err = m.conn.ExecCtx(ctx, query, args...)
	return err
}
