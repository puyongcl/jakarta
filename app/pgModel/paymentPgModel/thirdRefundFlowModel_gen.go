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
	thirdRefundFlowFieldNames          = builder.RawFieldNames(&ThirdRefundFlow{}, true)
	thirdRefundFlowRows                = strings.Join(thirdRefundFlowFieldNames, ",")
	thirdRefundFlowRowsExpectAutoSet   = strings.Join(stringx.Remove(thirdRefundFlowFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	thirdRefundFlowRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(thirdRefundFlowFieldNames, "flow_no", "create_time", "update_time", "create_at", "update_at"))
)

type (
	thirdRefundFlowModel interface {
		Insert(ctx context.Context, data *ThirdRefundFlow) (sql.Result, error)
		FindOne(ctx context.Context, flowNo string) (*ThirdRefundFlow, error)
		Update(ctx context.Context, data *ThirdRefundFlow) error
		Delete(ctx context.Context, flowNo string) error
	}

	defaultThirdRefundFlowModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ThirdRefundFlow struct {
		CreateTime         time.Time    `db:"create_time"`
		UpdateTime         time.Time    `db:"update_time"`
		FlowNo             string       `db:"flow_no"`
		PayFlowNo          string       `db:"pay_flow_no"`    // 原支付流水号
		TransactionId      string       `db:"transaction_id"` // 微信支付交易订单号
		OrderId            string       `db:"order_id"`
		Reason             string       `db:"reason"`
		PayAmount          int64        `db:"pay_amount"`
		RefundAmount       int64        `db:"refund_amount"`
		ActualRefundAmount int64        `db:"actual_refund_amount"`
		RefundStatus       int64        `db:"refund_status"`
		RefundTime         sql.NullTime `db:"refund_time"`
		Uid                int64        `db:"uid"`
		ReceivedAccount    string       `db:"received_account"` // 用户接收退款账户
		WxStatus           string       `db:"wx_status"`        // 微信支付返回的状态
		Remark             string       `db:"remark"`
	}
)

func newThirdRefundFlowModel(conn sqlx.SqlConn) *defaultThirdRefundFlowModel {
	return &defaultThirdRefundFlowModel{
		conn:  conn,
		table: `"jakarta"."third_refund_flow"`,
	}
}

func (m *defaultThirdRefundFlowModel) Delete(ctx context.Context, flowNo string) error {
	query := fmt.Sprintf("delete from %s where flow_no = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, flowNo)
	return err
}

func (m *defaultThirdRefundFlowModel) FindOne(ctx context.Context, flowNo string) (*ThirdRefundFlow, error) {
	query := fmt.Sprintf("select %s from %s where flow_no = $1 limit 1", thirdRefundFlowRows, m.table)
	var resp ThirdRefundFlow
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

func (m *defaultThirdRefundFlowModel) Insert(ctx context.Context, data *ThirdRefundFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)", m.table, thirdRefundFlowRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.PayFlowNo, data.TransactionId, data.OrderId, data.Reason, data.PayAmount, data.RefundAmount, data.ActualRefundAmount, data.RefundStatus, data.RefundTime, data.Uid, data.ReceivedAccount, data.WxStatus, data.Remark)
	return ret, err
}

func (m *defaultThirdRefundFlowModel) Update(ctx context.Context, data *ThirdRefundFlow) error {
	query := fmt.Sprintf("update %s set %s where flow_no = $1", m.table, thirdRefundFlowRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.PayFlowNo, data.TransactionId, data.OrderId, data.Reason, data.PayAmount, data.RefundAmount, data.ActualRefundAmount, data.RefundStatus, data.RefundTime, data.Uid, data.ReceivedAccount, data.WxStatus, data.Remark)
	return err
}

func (m *defaultThirdRefundFlowModel) tableName() string {
	return m.table
}