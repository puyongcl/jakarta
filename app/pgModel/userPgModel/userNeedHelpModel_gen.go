// Code generated by goctl. DO NOT EDIT!

package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userNeedHelpFieldNames          = builder.RawFieldNames(&UserNeedHelp{}, true)
	userNeedHelpRows                = strings.Join(userNeedHelpFieldNames, ",")
	userNeedHelpRowsExpectAutoSet   = strings.Join(stringx.Remove(userNeedHelpFieldNames, "create_time", "update_time", "create_t", "update_at"), ",")
	userNeedHelpRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(userNeedHelpFieldNames, "id", "create_time", "update_time", "create_at", "update_at"))
)

type (
	userNeedHelpModel interface {
		Insert(ctx context.Context, data *UserNeedHelp) (sql.Result, error)
		FindOne(ctx context.Context, id string) (*UserNeedHelp, error)
		Update(ctx context.Context, data *UserNeedHelp) error
		Delete(ctx context.Context, id string) error
	}

	defaultUserNeedHelpModel struct {
		conn  sqlx.SqlConn
		table string
	}

	UserNeedHelp struct {
		CreateTime       time.Time     `db:"create_time"`
		UpdateTime       time.Time     `db:"update_time"`
		Id               string        `db:"id"`
		Uid              int64         `db:"uid"`
		ListenerUid      int64         `db:"listener_uid"`
		ReportContent    string        `db:"report_content"`
		ReportTag        pq.Int64Array `db:"report_tag"`
		Attachment       string        `db:"attachment"`
		State            int64         `db:"state"`
		Remark           string        `db:"remark"`
		Avatar           string        `db:"avatar"`
		NickName         string        `db:"nick_name"`
		ListenerNickName string        `db:"listener_nick_name"`
		ListenerAvatar   string        `db:"listener_avatar"`
	}
)

func newUserNeedHelpModel(conn sqlx.SqlConn) *defaultUserNeedHelpModel {
	return &defaultUserNeedHelpModel{
		conn:  conn,
		table: `"jakarta"."user_need_help"`,
	}
}

func (m *defaultUserNeedHelpModel) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("delete from %s where id = $1", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultUserNeedHelpModel) FindOne(ctx context.Context, id string) (*UserNeedHelp, error) {
	query := fmt.Sprintf("select %s from %s where id = $1 limit 1", userNeedHelpRows, m.table)
	var resp UserNeedHelp
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

func (m *defaultUserNeedHelpModel) Insert(ctx context.Context, data *UserNeedHelp) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", m.table, userNeedHelpRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.ReportContent, data.ReportTag, data.Attachment, data.State, data.Remark, data.Avatar, data.NickName, data.ListenerNickName, data.ListenerAvatar)
	return ret, err
}

func (m *defaultUserNeedHelpModel) Update(ctx context.Context, data *UserNeedHelp) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, userNeedHelpRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Id, data.Uid, data.ListenerUid, data.ReportContent, data.ReportTag, data.Attachment, data.State, data.Remark, data.Avatar, data.NickName, data.ListenerNickName, data.ListenerAvatar)
	return err
}

func (m *defaultUserNeedHelpModel) tableName() string {
	return m.table
}
