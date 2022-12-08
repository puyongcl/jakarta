package userPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/userkey"
	"strings"
	"time"
)

type JoinUserModel struct {
	conn sqlx.SqlConn
}

func NewJoinUserModel(conn sqlx.SqlConn) *JoinUserModel {
	return &JoinUserModel{
		conn: conn,
	}
}

type UserListDetail struct {
	Uid             int64        `db:"uid"`
	CreateTime      time.Time    `db:"create_time"`        // 注册时间
	Channel         string       `db:"channel"`            // 获客渠道
	AuthKey         string       `db:"auth_key"`           // 登陆方式
	AuthType        string       `db:"auth_type"`          // 登陆方式
	FreeTime        sql.NullTime `db:"free_time"`          // 解封时间
	BanReason       string       `db:"ban_reason"`         // 封禁原因
	CostAmountSum   int64        `db:"cost_amount_sum"`    // 支付成功总额
	RefundAmountSum int64        `db:"refund_amount_sum"`  // 退款总额
	PaidOrderCnt    int64        `db:"paid_order_cnt"`     // 支付成功订单数量
	RefundOrderCnt  int64        `db:"refund_order_cnt"`   // 退款订单数
	NoCondRefundCnt int64        `db:"no_cond_refund_cnt"` // 无条件退款机会
}

var (
	userListDetailFieldNames = builder.RawFieldNames(&UserListDetail{}, true)
	userListDetailRows       = strings.Join(userListDetailFieldNames, ",")
)

func (m *JoinUserModel) FindUserList(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, isPaidUser bool, userType int64, channel string, pageNo, pageSize int64) ([]*UserListDetail, error) {
	rb := squirrel.Select(`uid, user_auth.create_time, channel, auth_key, auth_type, free_time, ban_reason, cost_amount_sum, refund_amount_sum, paid_order_cnt, refund_order_cnt, no_cond_refund_cnt`).From(`"jakarta"."user_auth"`)
	argNo := 1
	if isPaidUser {
		rb = rb.RightJoin(`"jakarta"."user_stat" using(uid)`)
	}
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."create_time" > $%d`, argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."create_time" < $%d`, argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."create_time" between $%d and $%d`, argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	if userType != 0 {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."user_type" = $%d`, argNo), userType)
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."user_type" = ALL($%d)`, argNo), pq.Int64Array{userkey.UserTypeListener, userkey.UserTypeNormalUser})
		argNo++
	}
	if channel != "" {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."channel" = $%d`, argNo), channel)
		argNo++
	}
	if isPaidUser {
		rb = rb.Where(fmt.Sprintf(`"user_stat"."cost_amount_sum" > 0`))
	}

	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	query, args, err := rb.OrderBy(`"user_auth"."create_time" DESC`).Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*UserListDetail, 0)
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *JoinUserModel) FindUserListCount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, isPaidUser bool, userType int64, channel string) (int64, error) {
	rb := squirrel.Select("count(uid)").From(`"jakarta"."user_auth"`)
	argNo := 1
	if isPaidUser {
		rb = rb.RightJoin(`"jakarta"."user_stat" using(uid)`)
	}
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."create_time" > $%d`, argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."create_time" < $%d`, argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."create_time" between $%d and $%d`, argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	if userType != 0 {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."user_type" = $%d`, argNo), userType)
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."user_type" = ALL($%d)`, argNo), pq.Int64Array{userkey.UserTypeListener, userkey.UserTypeNormalUser})
		argNo++
	}
	if channel != "" {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."channel" = $%d`, argNo), channel)
		argNo++
	}
	if isPaidUser {
		rb = rb.Where(fmt.Sprintf(`"user_stat"."cost_amount_sum" > 0`))
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var cnt int64
	err = m.conn.QueryRowCtx(ctx, &cnt, query, args...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}

// 统计活跃用户
func (m *JoinUserModel) StatLoginUserSum(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string, userType int64) (int64, error) {
	rb := squirrel.Select("count(uid)").From(`"jakarta"."user_auth"`).RightJoin(`"jakarta"."user_login_state" using(uid)`)

	argNo := 1
	if createTimeStart != nil && createTimeEnd == nil {
		rb = rb.Where(fmt.Sprintf(`"user_login_state"."login_time" > $%d`, argNo), createTimeStart)
		argNo++
	}
	if createTimeStart == nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf(`"user_login_state"."login_time" < $%d`, argNo), createTimeEnd)
		argNo++
	}
	if createTimeStart != nil && createTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf(`"user_login_state"."login_time" between $%d and $%d`, argNo, argNo+1), createTimeStart, createTimeEnd)
		argNo += 2
	}
	if userType != 0 {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."user_type" = $%d`, argNo), userType)
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."user_type" = ALL($%d)`, argNo), pq.Int64Array{userkey.UserTypeListener, userkey.UserTypeNormalUser})
		argNo++
	}
	if channel != "" {
		rb = rb.Where(fmt.Sprintf(`"user_auth"."channel" = $%d`, argNo), channel)
		argNo++
	}

	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var cnt int64
	err = m.conn.QueryRowCtx(ctx, &cnt, query, args...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}
