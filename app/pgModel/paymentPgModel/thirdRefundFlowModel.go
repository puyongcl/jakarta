package paymentPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ThirdRefundFlowModel = (*customThirdRefundFlowModel)(nil)

type (
	// ThirdRefundFlowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customThirdRefundFlowModel.
	ThirdRefundFlowModel interface {
		thirdRefundFlowModel
		FindCount(ctx context.Context, orderId string, tradeState []int64) (int64, error)
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *ThirdRefundFlow) (sql.Result, error)
		UpdatePart(ctx context.Context, newData *ThirdRefundFlow) error
	}

	customThirdRefundFlowModel struct {
		*defaultThirdRefundFlowModel
	}
)

// NewThirdRefundFlowModel returns a model for the database table.
func NewThirdRefundFlowModel(conn sqlx.SqlConn) ThirdRefundFlowModel {
	return &customThirdRefundFlowModel{
		defaultThirdRefundFlowModel: newThirdRefundFlowModel(conn),
	}
}

const findCountRefundFlowSQL = "select count(flow_no) from %s where order_id = $1 and refund_status = any ($2)"

func (m *defaultThirdRefundFlowModel) FindCount(ctx context.Context, orderId string, tradeState []int64) (int64, error) {
	query := fmt.Sprintf(findCountRefundFlowSQL, m.table)
	var cnt int64
	err := m.conn.QueryRowCtx(ctx, &cnt, query, orderId, pq.Int64Array(tradeState))
	return cnt, err
}

func (m *defaultThirdRefundFlowModel) Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultThirdRefundFlowModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ThirdRefundFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, thirdRefundFlowRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.FlowNo, data.PayFlowNo, data.TransactionId, data.OrderId, data.Reason, data.PayAmount, data.RefundAmount, data.ActualRefundAmount, data.RefundStatus, data.RefundTime, data.Uid, data.ReceivedAccount, data.WxStatus, data.Remark)
	return ret, err
}

func (m *defaultThirdRefundFlowModel) UpdatePart(ctx context.Context, newData *ThirdRefundFlow) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.RefundStatus != 0 {
		rb = rb.Set("refund_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RefundStatus))
		argNo++
	}
	if newData.TransactionId != "" {
		rb = rb.Set("transaction_id", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TransactionId))
		argNo++
	}
	if newData.ActualRefundAmount != 0 {
		rb = rb.Set("actual_refund_amount", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ActualRefundAmount))
		argNo++
	}
	if newData.WxStatus != "" {
		rb = rb.Set("wx_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.WxStatus))
		argNo++
	}
	if newData.ReceivedAccount != "" {
		rb = rb.Set("received_account", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ReceivedAccount))
		argNo++
	}
	if newData.RefundTime.Valid {
		rb = rb.Set("refund_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RefundTime.Time))
		argNo++
	}

	query, args, err := rb.Where(fmt.Sprintf("flow_no = $%d", argNo), newData.FlowNo).ToSql()
	if err != nil {
		return err
	}

	_, err = m.conn.ExecCtx(ctx, query, args...)
	return err
}
