package listenerPgModel

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ListenerContractModel = (*customListenerContractModel)(nil)

type (
	// ListenerContractModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerContractModel.
	ListenerContractModel interface {
		listenerContractModel
		UpdateListenerContract(ctx context.Context, id string, contractFile string, state int64, remark string) error
	}

	customListenerContractModel struct {
		*defaultListenerContractModel
	}
)

// NewListenerContractModel returns a model for the database table.
func NewListenerContractModel(conn sqlx.SqlConn) ListenerContractModel {
	return &customListenerContractModel{
		defaultListenerContractModel: newListenerContractModel(conn),
	}
}

func (m *defaultListenerContractModel) UpdateListenerContract(ctx context.Context, id string, contractFile string, state int64, remark string) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if contractFile != "" {
		rb = rb.Set("contract_file", squirrel.Expr(fmt.Sprintf("$%d", argNo), contractFile))
		argNo++
	}
	if state != 0 {
		rb = rb.Set("state", squirrel.Expr(fmt.Sprintf("$%d", argNo), state))
		argNo++
	}
	if remark != "" {
		rb = rb.Set("remark", squirrel.Expr(fmt.Sprintf("$%d", argNo), remark))
		argNo++
	}
	query, args, err := rb.Where(fmt.Sprintf("id = $%d", argNo), id).ToSql()
	if err != nil {
		return err
	}
	_, err = m.conn.ExecCtx(ctx, query, args...)
	return err
}
