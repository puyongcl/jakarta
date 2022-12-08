// Code generated by goctl. DO NOT EDIT!

package paymentPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	thirdPaymentFlowFieldNames          = builder.RawFieldNames(&ThirdPaymentFlow{}, true)
	thirdPaymentFlowRows                = strings.Join(thirdPaymentFlowFieldNames, ",")
	thirdPaymentFlowRowsExpectAutoSet   = strings.Join(stringx.Remove(thirdPaymentFlowFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	thirdPaymentFlowRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(thirdPaymentFlowFieldNames, "flow_no", "create_time", "update_time", "create_at", "update_at"))
)

type (
	thirdPaymentFlowModel interface {
		Insert(ctx context.Context, data *ThirdPaymentFlow) (sql.Result, error)
		FindOne(ctx context.Context, flowNo string) (*ThirdPaymentFlow, error)
		Update(ctx context.Context, data *ThirdPaymentFlow) error
		Delete(ctx context.Context, flowNo string) error
	}

	defaultThirdPaymentFlowModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ThirdPaymentFlow struct {
		CreateTime      time.Time    `db:"create_time"`
		UpdateTime      time.Time    `db:"update_time"`
		FlowNo          string       `db:"flow_no"`
		Uid             int64        `db:"uid"`
		PayMode         string       `db:"pay_mode"`         // 支付方式 1:微信支付
		TradeType       string       `db:"trade_type"`       // 第三方支付类型
		TradeState      string       `db:"trade_state"`      // 第三方交易状态
		PayAmount       int64        `db:"pay_amount"`       // 支付总金额(分)
		TransactionId   string       `db:"transaction_id"`   // 第三方支付单号
		TradeStateDesc  string       `db:"trade_state_desc"` // 支付状态描述
		OrderId         string       `db:"order_id"`         // 业务订单号
		OrderType       int64        `db:"order_type"`       // 订单类型
		PayStatus       int64        `db:"pay_status"`
		PayTime         sql.NullTime `db:"pay_time"` // 支付成功时间
		BankType        string       `db:"bank_type"`
		ActualPayAmount int64        `db:"actual_pay_amount"` // 实际支付金额
		Description     string       `db:"description"`       // 商品描述
	}
)

func newThirdPaymentFlowModel(conn sqlx.SqlConn) *defaultThirdPaymentFlowModel {
	return &defaultThirdPaymentFlowModel{
		conn:  conn,
		table: `"jakarta"."third_payment_flow"`,
	}
}

func (m *defaultThirdPaymentFlowModel) Delete(ctx context.Context, flowNo string) error {
	query := fmt.Sprintf("delete from %s where flow_no = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, flowNo)
	return err
}

func (m *defaultThirdPaymentFlowModel) FindOne(ctx context.Context, flowNo string) (*ThirdPaymentFlow, error) {
	query := fmt.Sprintf("select %s from %s where flow_no = $1 limit 1", thirdPaymentFlowRows, m.table)
	var resp ThirdPaymentFlow
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

func (m *defaultThirdPaymentFlowModel) Insert(ctx context.Context, data *ThirdPaymentFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)", m.table, thirdPaymentFlowRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.Uid, data.PayMode, data.TradeType, data.TradeState, data.PayAmount, data.TransactionId, data.TradeStateDesc, data.OrderId, data.OrderType, data.PayStatus, data.PayTime, data.BankType, data.ActualPayAmount, data.Description)
	return ret, err
}

func (m *defaultThirdPaymentFlowModel) Update(ctx context.Context, data *ThirdPaymentFlow) error {
	query := fmt.Sprintf("update %s set %s where flow_no = $1", m.table, thirdPaymentFlowRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.Uid, data.PayMode, data.TradeType, data.TradeState, data.PayAmount, data.TransactionId, data.TradeStateDesc, data.OrderId, data.OrderType, data.PayStatus, data.PayTime, data.BankType, data.ActualPayAmount, data.Description)
	return err
}

func (m *defaultThirdPaymentFlowModel) tableName() string {
	return m.table
}