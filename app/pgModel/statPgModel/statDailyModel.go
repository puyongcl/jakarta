package statPgModel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StatDailyModel = (*customStatDailyModel)(nil)

type (
	// StatDailyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStatDailyModel.
	StatDailyModel interface {
		statDailyModel
	}

	customStatDailyModel struct {
		*defaultStatDailyModel
	}
)

// NewStatDailyModel returns a model for the database table.
func NewStatDailyModel(conn sqlx.SqlConn) StatDailyModel {
	return &customStatDailyModel{
		defaultStatDailyModel: newStatDailyModel(conn),
	}
}
