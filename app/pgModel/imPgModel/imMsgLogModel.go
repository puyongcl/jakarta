package imPgModel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ImMsgLogModel = (*customImMsgLogModel)(nil)

type (
	// ImMsgLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customImMsgLogModel.
	ImMsgLogModel interface {
		imMsgLogModel
	}

	customImMsgLogModel struct {
		*defaultImMsgLogModel
	}
)

// NewImMsgLogModel returns a model for the database table.
func NewImMsgLogModel(conn sqlx.SqlConn) ImMsgLogModel {
	return &customImMsgLogModel{
		defaultImMsgLogModel: newImMsgLogModel(conn),
	}
}
