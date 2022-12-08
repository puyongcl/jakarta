package orderPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ChatOrderPricingPlanModel = (*customChatOrderPricingPlanModel)(nil)

type (
	// ChatOrderPricingPlanModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatOrderPricingPlanModel.
	ChatOrderPricingPlanModel interface {
		chatOrderPricingPlanModel
		FindPriceConfig(ctx context.Context) (*ChatOrderPricingPlan, error)
	}

	customChatOrderPricingPlanModel struct {
		*defaultChatOrderPricingPlanModel
	}
)

// NewChatOrderPricingPlanModel returns a model for the database table.
func NewChatOrderPricingPlanModel(conn sqlx.SqlConn) ChatOrderPricingPlanModel {
	return &customChatOrderPricingPlanModel{
		defaultChatOrderPricingPlanModel: newChatOrderPricingPlanModel(conn),
	}
}

func (m *defaultChatOrderPricingPlanModel) FindPriceConfig(ctx context.Context) (*ChatOrderPricingPlan, error) {
	query := fmt.Sprintf("select %s from %s where state = 2 limit 1", chatOrderPricingPlanRows, m.table)
	var resp ChatOrderPricingPlan
	err := m.conn.QueryRowCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
