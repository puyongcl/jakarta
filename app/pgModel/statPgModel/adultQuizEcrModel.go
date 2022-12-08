package statPgModel

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AdultQuizEcrModel = (*customAdultQuizEcrModel)(nil)

type (
	// AdultQuizEcrModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdultQuizEcrModel.
	AdultQuizEcrModel interface {
		adultQuizEcrModel
	}

	customAdultQuizEcrModel struct {
		*defaultAdultQuizEcrModel
	}
)

// NewAdultQuizEcrModel returns a model for the database table.
func NewAdultQuizEcrModel(conn sqlx.SqlConn) AdultQuizEcrModel {
	return &customAdultQuizEcrModel{
		defaultAdultQuizEcrModel: newAdultQuizEcrModel(conn),
	}
}
