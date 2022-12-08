package listenerPgModel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ListenerStatAverageModel = (*customListenerStatAverageModel)(nil)

type (
	// ListenerStatAverageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerStatAverageModel.
	ListenerStatAverageModel interface {
		listenerStatAverageModel
	}

	customListenerStatAverageModel struct {
		*defaultListenerStatAverageModel
	}
)

// NewListenerStatAverageModel returns a model for the database table.
func NewListenerStatAverageModel(conn sqlx.SqlConn) ListenerStatAverageModel {
	return &customListenerStatAverageModel{
		defaultListenerStatAverageModel: newListenerStatAverageModel(conn),
	}
}
