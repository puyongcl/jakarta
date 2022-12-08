package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
)

var _ UserProfileModel = (*customUserProfileModel)(nil)

type (
	// UserProfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserProfileModel.
	UserProfileModel interface {
		userProfileModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *UserProfile) (sql.Result, error)
		Update2(ctx context.Context, session sqlx.Session, data *UserProfile) error
	}

	customUserProfileModel struct {
		*defaultUserProfileModel
	}
)

// NewUserProfileModel returns a model for the database table.
func NewUserProfileModel(conn sqlx.SqlConn, c cache.CacheConf) UserProfileModel {
	return &customUserProfileModel{
		defaultUserProfileModel: newUserProfileModel(conn, c),
	}
}

func (m *defaultUserProfileModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultUserProfileModel) InsertTrans(ctx context.Context, session sqlx.Session, data *UserProfile) (sql.Result, error) {
	jakartaUserProfileUidKey := fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, data.Uid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8)", m.table, userProfileRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Uid, data.Nickname, data.Avatar, data.Gender, data.Introduction, data.Constellation, data.Birthday, data.PhoneNumber)
	}, jakartaUserProfileUidKey)
	return ret, err
}

func (m *defaultUserProfileModel) Update2(ctx context.Context, session sqlx.Session, newData *UserProfile) error {
	data, err := m.FindOne(ctx, newData.Uid)
	if err != nil {
		return err
	}
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.Nickname != data.Nickname {
		rb = rb.Set("nickname", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Nickname))
		argNo++
	}
	if newData.Avatar != data.Avatar {
		rb = rb.Set("avatar", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Avatar))
		argNo++
	}
	var b1, b2 string
	if newData.Birthday.Valid {
		b1 = newData.Birthday.Time.Format(db.DateFormat)
	}
	if data.Birthday.Valid {
		b2 = data.Birthday.Time.Format(db.DateFormat)
	}
	if b1 != b2 {
		rb = rb.Set("birthday", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Birthday))
		argNo++
	}
	if newData.Gender != data.Gender {
		rb = rb.Set("gender", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Gender))
		argNo++
	}
	if newData.Introduction != data.Introduction {
		rb = rb.Set("introduction", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Introduction))
		argNo++
	}
	if newData.Constellation != data.Constellation {
		rb = rb.Set("constellation", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Constellation))
		argNo++
	}
	if newData.PhoneNumber != data.PhoneNumber {
		rb = rb.Set("phone_number", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.PhoneNumber))
		argNo++
	}
	if argNo == 1 {
		return nil
	}

	query, args, err := rb.Where(fmt.Sprintf("uid = $%d", argNo), data.Uid).ToSql()
	if err != nil {
		return err
	}

	jakartaUserProfileUidKey := fmt.Sprintf("%s%v", cacheJakartaUserProfileUidPrefix, newData.Uid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return session.ExecCtx(ctx, query, args...)
	}, jakartaUserProfileUidKey)
	return err
}
