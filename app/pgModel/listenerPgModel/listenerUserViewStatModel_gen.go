// Code generated by goctl. DO NOT EDIT!

package listenerPgModel

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
	listenerUserViewStatFieldNames          = builder.RawFieldNames(&ListenerUserViewStat{}, true)
	listenerUserViewStatRows                = strings.Join(listenerUserViewStatFieldNames, ",")
	listenerUserViewStatRowsExpectAutoSet   = strings.Join(stringx.Remove(listenerUserViewStatFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	listenerUserViewStatRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(listenerUserViewStatFieldNames, "id", "create_time", "update_time", "create_at", "update_at"))
)

type (
	listenerUserViewStatModel interface {
		Insert(ctx context.Context, data *ListenerUserViewStat) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*ListenerUserViewStat, error)
		Update(ctx context.Context, data *ListenerUserViewStat) error
		Delete(ctx context.Context, id string) error
	}

	defaultListenerUserViewStatModel struct {
		conn  sqlx.SqlConn
		table string
	}

	ListenerUserViewStat struct {
		CreateTime  time.Time    `db:"create_time"`
		UpdateTime  time.Time    `db:"update_time"`
		Id          string       `db:"id"`
		Uid         int64        `db:"uid"`
		ListenerUid int64        `db:"listener_uid"`
		ViewTime    sql.NullTime `db:"view_time"`
		ViewCnt     int64        `db:"view_cnt"`
	}
)

func newListenerUserViewStatModel(conn sqlx.SqlConn) *defaultListenerUserViewStatModel {
	return &defaultListenerUserViewStatModel{
		conn:  conn,
		table: `"jakarta"."listener_user_view_stat"`,
	}
}

func (m *defaultListenerUserViewStatModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultListenerUserViewStatModel) FindOne(ctx context.Context, id string) (*ListenerUserViewStat, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", listenerUserViewStatRows, m.table)
	var resp ListenerUserViewStat
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

func (m *defaultListenerUserViewStatModel) Insert(ctx context.Context, data *ListenerUserViewStat) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5)", m.table, listenerUserViewStatRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.ViewTime, data.ViewCnt)
	return ret, err
}

func (m *defaultListenerUserViewStatModel) Update(ctx context.Context, data *ListenerUserViewStat) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, listenerUserViewStatRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.ViewTime, data.ViewCnt)
	return err
}

func (m *defaultListenerUserViewStatModel) tableName() string {
	return m.table
}