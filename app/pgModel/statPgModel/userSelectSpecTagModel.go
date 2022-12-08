package statPgModel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ UserSelectSpecTagModel = (*customUserSelectSpecTagModel)(nil)

type (
	// UserSelectSpecTagModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserSelectSpecTagModel.
	UserSelectSpecTagModel interface {
		userSelectSpecTagModel
	}

	customUserSelectSpecTagModel struct {
		*defaultUserSelectSpecTagModel
	}
)

// NewUserSelectSpecTagModel returns a model for the database table.
func NewUserSelectSpecTagModel(conn sqlx.SqlConn) UserSelectSpecTagModel {
	return &customUserSelectSpecTagModel{
		defaultUserSelectSpecTagModel: newUserSelectSpecTagModel(conn),
	}
}
