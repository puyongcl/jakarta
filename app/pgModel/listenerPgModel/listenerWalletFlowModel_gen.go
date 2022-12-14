// Code generated by goctl. DO NOT EDIT!

package listenerPgModel

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
	listenerWalletFlowFieldNames          = builder.RawFieldNames(&ListenerWalletFlow{}, true)
	listenerWalletFlowRows                = strings.Join(listenerWalletFlowFieldNames, ",")
	listenerWalletFlowRowsExpectAutoSet   = strings.Join(stringx.Remove(listenerWalletFlowFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	listenerWalletFlowRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(listenerWalletFlowFieldNames, "flow_no", "create_time", "update_time", "create_at", "update_at"))
)

type (
	listenerWalletFlowModel interface {
		Insert(ctx context.Context, data *ListenerWalletFlow) (sql.Result, error)
		FindOne(ctx context.Context, flowNo string) (*ListenerWalletFlow, error)
		Update(ctx context.Context, data *ListenerWalletFlow) error
		Delete(ctx context.Context, flowNo string) error
	}

	defaultListenerWalletFlowModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ListenerWalletFlow struct {
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		FlowNo      string    `db:"flow_no"`
		ListenerUid int64     `db:"listener_uid"`
		Amount      int64     `db:"amount"`      // 金额（分）
		OutId       string    `db:"out_id"`      // 订单号或提现任务编号，任务编号一对多条flow_no
		SettleType  int64     `db:"settle_type"` // 流水记录操作类型
		Remark      string    `db:"remark"`
		OutTime     time.Time `db:"out_time"` // 下单时间或者开始提现时间
		OrderAmount int64     `db:"order_amount"`
	}
)

func newListenerWalletFlowModel(conn sqlx.SqlConn) *defaultListenerWalletFlowModel {
	return &defaultListenerWalletFlowModel{
		conn:  conn,
		table: `"jakarta"."listener_wallet_flow"`,
	}
}

func (m *defaultListenerWalletFlowModel) Delete(ctx context.Context, flowNo string) error {
	query := fmt.Sprintf("delete from %s where flow_no = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, flowNo)
	return err
}

func (m *defaultListenerWalletFlowModel) FindOne(ctx context.Context, flowNo string) (*ListenerWalletFlow, error) {
	query := fmt.Sprintf("select %s from %s where flow_no = $1 limit 1", listenerWalletFlowRows, m.table)
	var resp ListenerWalletFlow
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

func (m *defaultListenerWalletFlowModel) Insert(ctx context.Context, data *ListenerWalletFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)", m.table, listenerWalletFlowRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.ListenerUid, data.Amount, data.OutId, data.SettleType, data.Remark, data.OutTime, data.OrderAmount)
	return ret, err
}

func (m *defaultListenerWalletFlowModel) Update(ctx context.Context, data *ListenerWalletFlow) error {
	query := fmt.Sprintf("update %s set %s where flow_no = $1", m.table, listenerWalletFlowRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.ListenerUid, data.Amount, data.OutId, data.SettleType, data.Remark, data.OutTime, data.OrderAmount)
	return err
}

func (m *defaultListenerWalletFlowModel) tableName() string {
	return m.table
}
