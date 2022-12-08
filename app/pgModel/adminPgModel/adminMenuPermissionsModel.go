package adminPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	"strings"
)

var _ AdminMenuPermissionsModel = (*customAdminMenuPermissionsModel)(nil)

var (
	adminMenuPermissionsUidRowsExpectAutoSet = strings.Join(stringx.Remove(adminMenuPermissionsFieldNames, "create_time", "update_time", "create_t", "update_at", "id"), ",")
)

type (
	// AdminMenuPermissionsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAdminMenuPermissionsModel.
	AdminMenuPermissionsModel interface {
		adminMenuPermissionsModel
		Find(ctx context.Context, uid int64, menu1Id, menu2Id, pageNo, pageSize int64) ([]*AdminMenuPermissions, error)
		FindCount(ctx context.Context, uid int64, menu1Id, menu2Id int64) (int64, error)
		Insert2(ctx context.Context, data *AdminMenuPermissions) (sql.Result, error)
		Delete2(ctx context.Context, uid int64, menu1Id, menu2Id int64) error
	}

	customAdminMenuPermissionsModel struct {
		*defaultAdminMenuPermissionsModel
	}
)

// NewAdminMenuPermissionsModel returns a model for the database table.
func NewAdminMenuPermissionsModel(conn sqlx.SqlConn) AdminMenuPermissionsModel {
	return &customAdminMenuPermissionsModel{
		defaultAdminMenuPermissionsModel: newAdminMenuPermissionsModel(conn),
	}
}

func (m *defaultAdminMenuPermissionsModel) Find(ctx context.Context, uid int64, menu1Id, menu2Id, pageNo, pageSize int64) ([]*AdminMenuPermissions, error) {
	rb := squirrel.Select(adminMenuPermissionsRows).From(m.table)
	rb = rb.Where("uid = $1", uid)
	argNo := 2
	if menu1Id != 0 {
		rb = rb.Where(fmt.Sprintf("menu_1_id = $%d", argNo), menu1Id)
		argNo++
	}
	if menu2Id != 0 {
		rb = rb.Where(fmt.Sprintf("menu_2_id = $%d", argNo), menu2Id)
		argNo++
	}
	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query, values, err := rb.OrderBy("create_time asc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*AdminMenuPermissions, 0)
	//err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultAdminMenuPermissionsModel) Insert2(ctx context.Context, data *AdminMenuPermissions) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5)", m.table, adminMenuPermissionsUidRowsExpectAutoSet)
	return m.conn.ExecCtx(ctx, query, data.Uid, data.Menu1Id, data.Menu2Id, data.MenuValue, data.State)
}

func (m *defaultAdminMenuPermissionsModel) Delete2(ctx context.Context, uid int64, menu1Id, menu2Id int64) error {
	query := fmt.Sprintf("delete from %s where uid = $1 and menu_1_id = $2 and menu_2_id = $3", m.table)
	_, err := m.conn.ExecCtx(ctx, query, uid, menu1Id, menu2Id)
	return err
}

func (m *defaultAdminMenuPermissionsModel) FindCount(ctx context.Context, uid int64, menu1Id, menu2Id int64) (int64, error) {
	rb := squirrel.Select("COUNT(id)").From(m.table)
	rb = rb.Where("uid = $1", uid)
	argNo := 2
	if menu1Id != 0 {
		rb = rb.Where(fmt.Sprintf("menu_1_id = $%d", argNo), menu1Id)
		argNo++
	}
	if menu2Id != 0 {
		rb = rb.Where(fmt.Sprintf("menu_2_id = $%d", argNo), menu2Id)
		argNo++
	}
	query, values, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var cnt int64
	//err = m.QueryRowsNoCacheCtx(ctx, &cnt, query, values...)
	err = m.conn.QueryRowCtx(ctx, &cnt, query, values...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}
