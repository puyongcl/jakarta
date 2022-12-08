package adminPgModel

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdminLogModel = (*customAdminLogModel)(nil)

type (
	// AdminLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminLogModel.
	AdminLogModel interface {
		adminLogModel
	}

	customAdminLogModel struct {
		*defaultAdminLogModel
	}
)

// NewAdminLogModel returns a model for the database table.
func NewAdminLogModel(conn sqlx.SqlConn) AdminLogModel {
	return &customAdminLogModel{
		defaultAdminLogModel: newAdminLogModel(conn),
	}
}
