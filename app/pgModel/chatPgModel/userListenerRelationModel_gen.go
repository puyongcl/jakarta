// Code generated by goctl. DO NOT EDIT!

package chatPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userListenerRelationFieldNames          = builder.RawFieldNames(&UserListenerRelation{}, true)
	userListenerRelationRows                = strings.Join(userListenerRelationFieldNames, ",")
	userListenerRelationRowsExpectAutoSet   = strings.Join(stringx.Remove(userListenerRelationFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	userListenerRelationRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(userListenerRelationFieldNames, "id", "create_time", "update_time", "create_at", "update_at"))
)

type (
	userListenerRelationModel interface {
		Insert(ctx context.Context, data *UserListenerRelation) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*UserListenerRelation, error)
		Update(ctx context.Context, data *UserListenerRelation) error
		Delete(ctx context.Context, id string) error
	}

	defaultUserListenerRelationModel struct {
		conn  sqlx.SqlConn
		table string
	}

	UserListenerRelation struct {
		CreateTime  time.Time `db:"create_time"`
		UpdateTime  time.Time `db:"update_time"`
		Id          string    `db:"id"`
		Uid         int64     `db:"uid"`
		ListenerUid int64     `db:"listener_uid"`
		TotalScore  int64     `db:"total_score"` // 交互产生的累计分值
	}
)

func newUserListenerRelationModel(conn sqlx.SqlConn) *defaultUserListenerRelationModel {
	return &defaultUserListenerRelationModel{
		conn:  conn,
		table: `"jakarta"."user_listener_relation"`,
	}
}

func (m *defaultUserListenerRelationModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultUserListenerRelationModel) FindOne(ctx context.Context, id string) (*UserListenerRelation, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", userListenerRelationRows, m.table)
	var resp UserListenerRelation
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserListenerRelationModel) Insert(ctx context.Context, data *UserListenerRelation) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4)", m.table, userListenerRelationRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.TotalScore)
	return ret, err
}

func (m *defaultUserListenerRelationModel) Update(ctx context.Context, data *UserListenerRelation) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, userListenerRelationRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.TotalScore)
	return err
}

func (m *defaultUserListenerRelationModel) tableName() string {
	return m.table
}