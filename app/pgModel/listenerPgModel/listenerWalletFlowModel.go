package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"
)

var _ ListenerWalletFlowModel = (*customListenerWalletFlowModel)(nil)

type (
	// ListenerWalletFlowModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerWalletFlowModel.
	ListenerWalletFlowModel interface {
		listenerWalletFlowModel
		InsertOutLogTrans(ctx context.Context, data *ListenerWalletFlow) (sql.Result, error)
		UpdateRefundIncomeLog(ctx context.Context, outId string, settleType int64, remark string) (sql.Result, error)
		UpdateConfirmIncomeLog(ctx context.Context, outId string, settleType, amount int64, remark string) (sql.Result, error)
		UpdateOutLog(ctx context.Context, flowNo, outId string, settleType int64, remark string) (sql.Result, error)
		InsertIncomeLog(ctx context.Context, data *ListenerWalletFlow) (sql.Result, error)
		Find(ctx context.Context, listenerUid int64, settleType []int64, pageNo, pageSize int64) ([]*ListenerWalletFlow, error)
		FindMoveCashList(ctx context.Context, listenerUid int64, settleType []int64, pageNo, pageSize int64) ([]*ListenerWalletFlow, error)
		CountMoveCashList(ctx context.Context, listenerUid int64, settleType []int64) (int64, error)
		FindOneIncomeFlow(ctx context.Context, outId string, settleType int64) (*ListenerWalletFlow, error)
		ResetWalletFlow(ctx context.Context, flowNo string, amount, orderAmount, settleType int64) (sql.Result, error)
		SumListenerAmount(ctx context.Context, listenerUid int64, startOutTime, endOutTime *time.Time, st []int64) (int64, error)
		SumListenerOrderAmount(ctx context.Context, listenerUid int64, startOutTime, endOutTime *time.Time, st []int64) (int64, error)
	}

	customListenerWalletFlowModel struct {
		*defaultListenerWalletFlowModel
	}
)

// NewListenerWalletFlowModel returns a model for the database table.
func NewListenerWalletFlowModel(conn sqlx.SqlConn) ListenerWalletFlowModel {
	return &customListenerWalletFlowModel{
		defaultListenerWalletFlowModel: newListenerWalletFlowModel(conn),
	}
}

func (m *defaultListenerWalletFlowModel) InsertOutLogTrans(ctx context.Context, data *ListenerWalletFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)", m.table, listenerWalletFlowRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.ListenerUid, data.Amount, data.OutId, data.SettleType, data.Remark, data.OutTime, data.OrderAmount)
	return ret, err
}

func (m *defaultListenerWalletFlowModel) InsertIncomeLog(ctx context.Context, data *ListenerWalletFlow) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)", m.table, listenerWalletFlowRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.FlowNo, data.ListenerUid, data.Amount, data.OutId, data.SettleType, data.Remark, data.OutTime, data.OrderAmount)
	return ret, err
}

func (m *defaultListenerWalletFlowModel) ResetWalletFlow(ctx context.Context, flowNo string, amount, orderAmount, settleType int64) (sql.Result, error) {
	q, a, err := squirrel.Update(m.table).Set("amount", squirrel.Expr("$1", amount)).Set("order_amount", squirrel.Expr("$2", orderAmount)).Set("settle_type", squirrel.Expr("$3", settleType)).Where("flow_no = $4", flowNo).ToSql()
	if err != nil {
		return nil, err
	}
	return m.conn.ExecCtx(ctx, q, a...)
}

func (m *defaultListenerWalletFlowModel) UpdateOutLog(ctx context.Context, flowNo, outId string, settleType int64, remark string) (sql.Result, error) {
	rb := squirrel.Update(m.table).Set("settle_type", squirrel.Expr("$1", settleType)).Set("remark", squirrel.Expr("$2", remark))
	argNo := 3
	if outId != "" {
		rb = rb.Set("out_id", squirrel.Expr(fmt.Sprintf("$%d", argNo), outId))
		argNo++
	}
	query, args, err := rb.Where(fmt.Sprintf("flow_no = $%d", argNo), flowNo).ToSql()
	if err != nil {
		return nil, err
	}
	return m.conn.ExecCtx(ctx, query, args...)
}

func (m *defaultListenerWalletFlowModel) UpdateRefundIncomeLog(ctx context.Context, outId string, settleType int64, remark string) (sql.Result, error) {
	query, args, err := squirrel.Update(m.table).Set("settle_type", squirrel.Expr("$1", settleType)).Set("remark", squirrel.Expr("$2", remark)).Where("out_id = $3", outId).ToSql()
	if err != nil {
		return nil, err
	}
	return m.conn.ExecCtx(ctx, query, args...)
}

func (m *defaultListenerWalletFlowModel) UpdateConfirmIncomeLog(ctx context.Context, outId string, settleType, amount int64, remark string) (sql.Result, error) {
	query, args, err := squirrel.Update(m.table).Set("settle_type", squirrel.Expr("$1", settleType)).Set("remark", squirrel.Expr("$2", remark)).Set("amount", squirrel.Expr("$3", amount)).Where("out_id = $4", outId).ToSql()
	if err != nil {
		return nil, err
	}
	return m.conn.ExecCtx(ctx, query, args...)
}

func (m *defaultListenerWalletFlowModel) Find(ctx context.Context, listenerUid int64, settleType []int64, pageNo, pageSize int64) ([]*ListenerWalletFlow, error) {
	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query, values, err := squirrel.Select(listenerWalletFlowRows).From(m.table).Where("listener_uid = $1", listenerUid).Where("settle_type = ANY($2)", pq.Int64Array(settleType)).OrderBy("out_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ListenerWalletFlow, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerWalletFlowModel) FindMoveCashList(ctx context.Context, listenerUid int64, settleType []int64, pageNo, pageSize int64) ([]*ListenerWalletFlow, error) {
	rb := squirrel.Select(listenerWalletFlowRows).From(m.table).Where("settle_type = ANY($1)", pq.Int64Array(settleType))
	if listenerUid != 0 {
		rb = rb.Where("listener_uid = $2", listenerUid)
	}
	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query, values, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ListenerWalletFlow, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerWalletFlowModel) CountMoveCashList(ctx context.Context, listenerUid int64, settleType []int64) (int64, error) {
	rb := squirrel.Select("count(flow_no)").From(m.table).Where("settle_type = ANY($1)", pq.Int64Array(settleType))
	if listenerUid != 0 {
		rb = rb.Where("listener_uid = $2", listenerUid)
	}

	query, values, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var resp int64
	err = m.conn.QueryRowCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultListenerWalletFlowModel) FindOneIncomeFlow(ctx context.Context, outId string, settleType int64) (*ListenerWalletFlow, error) {
	rb := squirrel.Select(listenerWalletFlowRows).From(m.table).Where("out_id = $1", outId)
	if settleType != 0 {
		rb = rb.Where("settle_type = $2", settleType)
	}

	q, a, err := rb.ToSql()
	if err != nil {
		return nil, err
	}
	var resp ListenerWalletFlow
	err = m.conn.QueryRowCtx(ctx, &resp, q, a...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerWalletFlowModel) SumListenerAmount(ctx context.Context, listenerUid int64, startOutTime, endOutTime *time.Time, st []int64) (int64, error) {
	rb := squirrel.Select("coalesce(sum(amount), 0)").From(m.table).Where("listener_uid = $1", listenerUid).Where("settle_type = ANY($2)", pq.Int64Array(st))
	argNo := 3
	if startOutTime != nil {
		rb = rb.Where(fmt.Sprintf("out_time > $%d", argNo), startOutTime)
		argNo++
	}
	if endOutTime != nil {
		rb = rb.Where(fmt.Sprintf("out_time < $%d", argNo), endOutTime)
		argNo++
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultListenerWalletFlowModel) SumListenerOrderAmount(ctx context.Context, listenerUid int64, startOutTime, endOutTime *time.Time, st []int64) (int64, error) {
	rb := squirrel.Select("coalesce(sum(order_amount), 0)").From(m.table).Where("listener_uid = $1", listenerUid).Where("settle_type = ANY($2)", pq.Int64Array(st))
	argNo := 3
	if startOutTime != nil {
		rb = rb.Where(fmt.Sprintf("out_time > $%d", argNo), startOutTime)
		argNo++
	}
	if endOutTime != nil {
		rb = rb.Where(fmt.Sprintf("out_time < $%d", argNo), endOutTime)
		argNo++
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.conn.QueryRowCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
