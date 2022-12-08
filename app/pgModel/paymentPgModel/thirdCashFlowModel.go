package paymentPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ThirdCashFlowModel = (*customThirdCashFlowModel)(nil)

type (
	// ThirdCashFlowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customThirdCashFlowModel.
	ThirdCashFlowModel interface {
		thirdCashFlowModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *ThirdCashFlow) (sql.Result, error)
		UpdatePart(ctx context.Context, newData *ThirdCashFlow) error
		FindOneByWalletFlowNo(ctx context.Context, flowNo string) (*ThirdCashFlow, error)
	}

	customThirdCashFlowModel struct {
		*defaultThirdCashFlowModel
	}
)

// NewThirdCashFlowModel returns a model for the database table.
func NewThirdCashFlowModel(conn sqlx.SqlConn) ThirdCashFlowModel {
	return &customThirdCashFlowModel{
		defaultThirdCashFlowModel: newThirdCashFlowModel(conn),
	}
}

func (m *defaultThirdCashFlowModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultThirdCashFlowModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ThirdCashFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)", m.table, thirdCashFlowRowsExpectAutoSet)
	ret, err := session.ExecCtx(ctx, query, data.FlowNo, data.WorkNumber, data.Amount, data.PhoneNumber, data.Uid, data.Name, data.IdNo, data.BankCardNo, data.TransactionNumber, data.PayStatus, data.PayTime, data.ErrMsg, data.WalletFlowNo)
	return ret, err
}

func (m *defaultThirdCashFlowModel) UpdatePart(ctx context.Context, newData *ThirdCashFlow) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.WorkNumber != "" {
		rb = rb.Set("work_number", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.WorkNumber))
		argNo++
	}
	if newData.TransactionNumber != "" {
		rb = rb.Set("transaction_number", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TransactionNumber))
		argNo++
	}
	if newData.PayStatus != 0 {
		rb = rb.Set("pay_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.PayStatus))
		argNo++
	}
	if newData.ErrMsg != "" {
		rb = rb.Set("err_msg", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ErrMsg))
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

func (m *defaultThirdCashFlowModel) FindOneByWalletFlowNo(ctx context.Context, flowNo string) (*ThirdCashFlow, error) {
	query := fmt.Sprintf("select %s from %s where wallet_flow_no = $1 limit 1", thirdCashFlowRows, m.table)
	var resp ThirdCashFlow
	err := m.conn.QueryRowCtx(ctx, &resp, query, flowNo)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
