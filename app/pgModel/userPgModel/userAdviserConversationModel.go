package userPgModel

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserAdviserConversationModel = (*customUserAdviserConversationModel)(nil)

type (
	// UserAdviserConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserAdviserConversationModel.
	UserAdviserConversationModel interface {
		userAdviserConversationModel
	}

	customUserAdviserConversationModel struct {
		*defaultUserAdviserConversationModel
	}
)

// NewUserAdviserConversationModel returns a model for the database table.
func NewUserAdviserConversationModel(conn sqlx.SqlConn, c cache.CacheConf) UserAdviserConversationModel {
	return &customUserAdviserConversationModel{
		defaultUserAdviserConversationModel: newUserAdviserConversationModel(conn, c),
	}
}
