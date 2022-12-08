package chatPgModel

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserListenerRelationModel = (*customUserListenerRelationModel)(nil)

type (
	// UserListenerRelationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserListenerRelationModel.
	UserListenerRelationModel interface {
		userListenerRelationModel
		AddScore(ctx context.Context, addScore int64, id string) error
		FindTopScoreList(ctx context.Context, uid, pageNo, pageSize int64) ([]*UserListenerRelation, error)
	}

	customUserListenerRelationModel struct {
		*defaultUserListenerRelationModel
	}
)

// NewUserListenerRelationModel returns a model for the database table.
func NewUserListenerRelationModel(conn sqlx.SqlConn) UserListenerRelationModel {
	return &customUserListenerRelationModel{
		defaultUserListenerRelationModel: newUserListenerRelationModel(conn),
	}
}

const addUserListenerRelationScoreSQL = `update %s set total_score = total_score + $1 where id = $2`

func (m *defaultUserListenerRelationModel) AddScore(ctx context.Context, addScore int64, id string) error {
	q := fmt.Sprintf(addUserListenerRelationScoreSQL, m.table)

	_, err := m.conn.ExecCtx(ctx, q, addScore, id)
	return err
}

func (m *defaultUserListenerRelationModel) FindTopScoreList(ctx context.Context, uid, pageNo, pageSize int64) ([]*UserListenerRelation, error) {
	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	q, a, err := squirrel.Select(userListenerRelationRows).From(m.table).Where("uid=$1", uid).OrderBy("total_score desc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*UserListenerRelation, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, q, a...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
