package orderPgModel

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatOrderStatusLogModel = (*customChatOrderStatusLogModel)(nil)

type (
	// ChatOrderStatusLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatOrderStatusLogModel.
	ChatOrderStatusLogModel interface {
		chatOrderStatusLogModel
		Find(ctx context.Context, orderId string, state []int64, pageNo, pageSize int64) ([]*ChatOrderStatusLog, error)
	}

	customChatOrderStatusLogModel struct {
		*defaultChatOrderStatusLogModel
	}
)

// NewChatOrderStatusLogModel returns a model for the database table.
func NewChatOrderStatusLogModel(conn sqlx.SqlConn) ChatOrderStatusLogModel {
	return &customChatOrderStatusLogModel{
		defaultChatOrderStatusLogModel: newChatOrderStatusLogModel(conn),
	}
}

func (m *defaultChatOrderStatusLogModel) Find(ctx context.Context, orderId string, state []int64, pageNo, pageSize int64) ([]*ChatOrderStatusLog, error) {
	rb := squirrel.Select(chatOrderStatusLogRows).From(m.table).Where("order_id = $1", orderId)
	if len(state) > 0 {
		rb = rb.Where("state = ANY($2)", pq.Int64Array(state)).Where("action_result = 0")
	}
	rb = rb.OrderBy("create_time ASC")

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query, values, err := rb.Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatOrderStatusLog, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
