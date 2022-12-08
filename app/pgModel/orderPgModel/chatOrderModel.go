package orderPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"
	"strings"
	"time"
)

var _ ChatOrderModel = (*customChatOrderModel)(nil)

type (
	// ChatOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatOrderModel.
	ChatOrderModel interface {
		chatOrderModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		FindCount(ctx context.Context, uid int64, listenerUid int64, orderType int64, state []int64, sec int64) (int64, error)
		FindCount2(ctx context.Context, uid int64, listenerUid int64, orderType int64, usedMinCond bool, notState []int64, sec int64) (int64, error)
		UpdateOrderState(ctx context.Context, orderId string, state int64, statesCond []int64) error
		UpdateOrderPaySuccess(ctx context.Context, orderId string, state, buyMinSum int64, startTime, expiryTime *time.Time, statesCond []int64) error
		UpdateOrderUse(ctx context.Context, orderId string, usedChatMinute int64, startTime, endTime *time.Time, state int64, statesCond []int64) error
		Find(ctx context.Context, uid int64, listenerUid int64, orderType int64, state []int64, listType int64, sort string, pageNo, pageSize int64) ([]*ChatOrder, error)
		FindExpireOrder(ctx context.Context, expiryTimeStart, expiryTimeEnd *time.Time, orderType int64, state []int64, pageNo, pageSize int64) ([]*ExpireVoiceChatOrder, error)
		SettleOrder(ctx context.Context, orderId string, shareAmount, amount, state int64, settleTime *time.Time) error
		UpdateOrderOpinionAndState(ctx context.Context, orderId string, newData *ChatOrder, statesCond []int64) error
		UpdateApplyRefund(ctx context.Context, orderId string, newData *ChatOrder, statesCond []int64) error
		FindListenerOrderOpinionList(ctx context.Context, listenerUid int64, star int64, pageNo int64, pageSize int64) ([]*ListenerOrderOpinion, error)
		FindNeedAutoProcessOrderList(ctx context.Context, state []int64, beforeTime *time.Time, pageNo, pageSize int64) ([]*AutoProcessOrder, error)
		FindGoodComment(ctx context.Context, pageNo, pageSize int64) ([]*ListenerGoodComment, error)
		FindLastCommentOrder(ctx context.Context, listenerUid int64, star int64) (*ListenerGoodComment, error)
		CountPaidUserCnt2(ctx context.Context, listenerUid int64, createTimeStart *time.Time, createTimeEnd *time.Time) (int64, error)
		ResetOrderAmount(ctx context.Context, orderId string, sumBuyMin int64, taxAmount int64) error
		ResetOrderAmount2(ctx context.Context, orderId string, sumBuyMin int64, taxAmount, platAmount, listenerAmount int64) error

		FindChatOrderFeedback(ctx context.Context, uid, pageNo, pageSize int64) ([]*ListenerFeedback, error)
		FindAllPaidOrder(ctx context.Context, pageNo, pageSize int64) ([]*ChatOrder, error)

		CountPaidOrderCnt2(ctx context.Context, listenerUid int64) (int64, error)
		FindAllComment(ctx context.Context, listenerUid int64, pageNo, pageSize int64) ([]*ListenerGoodComment, error)

		CountPaidUser2(ctx context.Context, listenerUid int64) (int64, error)
		CountOrder(ctx context.Context, listenerUid int64, state []int64) (int64, error)
		CountRepeatPaidUser2(ctx context.Context, listenerUid int64) (int64, error)
		SumOrderBuyUnitCnt(ctx context.Context, listenerUid int64) (int64, error)
		CountOrderRatingCnt(ctx context.Context, listenerUid int64, star int64) (int64, error)

		SumOrderBuyUnitCntByState(ctx context.Context, uid, listenerUid, state int64) (int64, error)
		SumUserOrderBuyUnitCnt(ctx context.Context, uid, listenerUid int64) (int64, error)
		SumUserOrderBuyUnitCntBefore(ctx context.Context, uid, listenerUid int64, endCreateTime *time.Time, notOrderId string) (int64, error)

		CountPaidUser(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error)
		CountRepeatPaidUser(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error)
		SumListenerPaidAmount(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error)
		CountPaidUserCnt(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		CountPaidOrderCnt(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		SumPaidAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		SumApplyRefundAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		SumRefundSuccessAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		SumConfirmAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		SumListenerAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		SumPlatformAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)

		CountPaidNewUserCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, uidStart, uidEnd int64, channel string) (int64, error)
		CountNewUserPaidOrderCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, uidStart, uidEnd int64, channel string) (int64, error)
		CountNewUserOrderRangeCreateTime(ctx context.Context, startTime, endTime *time.Time, uidStart, uidEnd int64, channel string, state []int64) (int64, error)
		CountRepeatPaidNewUserRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd int64, channel string) (int64, error)
		CountNewUserCommentOrderRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd int64, channel string) (int64, error)
		CountNewUserCommentOrderByStarRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd, star int64, channel string) (int64, error)
		CountNewUserOrderCntByTypeRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd, orderType int64, channel string) (int64, error)
		SumNewUserPaidAmountRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, uidStart, uidEnd int64, channel string) (int64, error)

		CountPaidUserCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		CountUserPaidOrderCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
		CountUserOrderRangeCreateTime(ctx context.Context, startTime, endTime *time.Time, channel string, state []int64) (int64, error)
		CountRepeatPaidUserRangeCreateTime(ctx context.Context, start, end *time.Time, channel string) (int64, error)
		CountUserCommentOrderRangeCreateTime(ctx context.Context, start, end *time.Time, channel string) (int64, error)
		CountUserCommentOrderByStarRangeCreateTime(ctx context.Context, start, end *time.Time, star int64, channel string) (int64, error)
		CountUserOrderCntByTypeRangeCreateTime(ctx context.Context, start, end *time.Time, orderType int64, channel string) (int64, error)
		SumUserPaidAmountRangeCreateTimeLtv(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string, uids []int64) (int64, error)
		SumUserPaidAmountRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error)
	}

	customChatOrderModel struct {
		*defaultChatOrderModel
	}
)

//
var listenerGoodCommentRows = strings.Join(builder.RawFieldNames(&ListenerGoodComment{}, true), ",")

type ListenerGoodComment struct {
	OrderId          string        `db:"order_id"`
	ListenerUid      int64         `db:"listener_uid"`
	Uid              int64         `db:"uid"`
	ListenerNickName string        `db:"listener_nick_name"` // XXX昵称
	NickName         string        `db:"nick_name"`          // 用户昵称
	CommentTime      sql.NullTime  `db:"comment_time"`       // 用户评价时间
	CommentTag       pq.Int64Array `db:"comment_tag"`        // 评价标签
}

//
var listenerFeedbackRows = strings.Join(builder.RawFieldNames(&ListenerFeedback{}, true), ",")

type ListenerFeedback struct {
	CreateTime       time.Time    `db:"create_time"`
	OrderId          string       `db:"order_id"`
	ListenerUid      int64        `db:"listener_uid"`
	ListenerNickName string       `db:"listener_nick_name"` // XXX昵称
	ListenerAvatar   string       `db:"listener_avatar"`    // XXX头像
	FeedbackTime     sql.NullTime `db:"feedback_time"`      // XXX反馈时间
	Feedback         string       `db:"feedback"`           // XXX反馈
}

// TODO 自定义select 字段 需要定义对应的结构体 用于接收结果 否则orm会报错 not matching destination to scan
var listenerOrderOpinionRows = strings.Join(builder.RawFieldNames(&ListenerOrderOpinion{}, true), ",")

type ListenerOrderOpinion struct {
	OrderId          string        `db:"order_id"`
	ListenerUid      int64         `db:"listener_uid"`
	ListenerNickName string        `db:"listener_nick_name"` // XXX昵称
	ListenerAvatar   string        `db:"listener_avatar"`    // XXX头像
	Uid              int64         `db:"uid"`
	OrderType        int64         `db:"order_type"`    // 聊天类型2文字4语音
	NickName         string        `db:"nick_name"`     // 用户昵称
	Avatar           string        `db:"avatar"`        // 用户头像
	CommentTime      sql.NullTime  `db:"comment_time"`  // 用户评价时间
	Comment          string        `db:"comment"`       // 用户评价内容
	CommentTag       pq.Int64Array `db:"comment_tag"`   // 评价标签
	Star             int64         `db:"star"`          // 用户评价1不满意3一般5满意
	ReplyTime        sql.NullTime  `db:"reply_time"`    // XXX回复时间
	Reply            string        `db:"reply"`         // XXX回复内容
	FeedbackTime     sql.NullTime  `db:"feedback_time"` // XXX反馈时间
	Feedback         string        `db:"feedback"`      // XXX反馈
}

// 自动进入下一阶段状态的订单
var autoProcessOrderRows = strings.Join(builder.RawFieldNames(&AutoProcessOrder{}, true), ",")

type AutoProcessOrder struct {
	CreateTime  time.Time `db:"create_time"`
	UpdateTime  time.Time `db:"update_time"` // 更新状态时间
	OrderId     string    `db:"order_id"`
	ListenerUid int64     `db:"listener_uid"`
	Uid         int64     `db:"uid"`
	OrderType   int64     `db:"order_type"`  // 聊天类型2文字4语音
	OrderState  int64     `db:"order_state"` // 订单状态
}

// 到期的订单
var expireVoiceChatOrderRows = strings.Join(builder.RawFieldNames(&ExpireVoiceChatOrder{}, true), ",")

type ExpireVoiceChatOrder struct {
	OrderId        string `db:"order_id"`
	Uid            int64  `db:"uid"`
	ListenerUid    int64  `db:"listener_uid"`
	BuyUnit        int64  `db:"buy_unit"`
	ChatUnitMinute int64  `db:"chat_unit_minute"`
	UsedChatMinute int64  `db:"used_chat_minute"`
}

// NewChatOrderModel returns a model for the database table.
func NewChatOrderModel(conn sqlx.SqlConn, c cache.CacheConf) ChatOrderModel {
	return &customChatOrderModel{
		defaultChatOrderModel: newChatOrderModel(conn, c),
	}
}

func (m *defaultChatOrderModel) FindCount(ctx context.Context, uid int64, listenerUid int64, orderType int64, state []int64, sec int64) (int64, error) {
	rb := squirrel.Select("COUNT(order_id)").From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf("uid > $%d", argNo), userkey.UidStart)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if orderType != 0 {
		rb = rb.Where(fmt.Sprintf("order_type = $%d", argNo), orderType)
		argNo++
	}
	if len(state) == 1 {
		rb = rb.Where(fmt.Sprintf("order_state = $%d", argNo), state[0])
		argNo++
	} else if len(state) > 1 {
		rb = rb.Where(fmt.Sprintf("order_state = ANY($%d)", argNo), pq.Int64Array(state))
		argNo++
	} else {
		rb = rb.Where("order_state != ALL('{1,2}')")
	}
	if sec != 0 { // 多少秒内
		t := time.Now().Add(-time.Duration(sec) * time.Second)
		rb = rb.Where(fmt.Sprintf("create_time > $%d", argNo), t)
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) FindCount2(ctx context.Context, uid int64, listenerUid int64, orderType int64, usedMinCond bool, notState []int64, sec int64) (int64, error) {
	rb := squirrel.Select("COUNT(order_id)").From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf("uid > $%d", argNo), userkey.UidStart)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if orderType != 0 {
		rb = rb.Where(fmt.Sprintf("order_type = $%d", argNo), orderType)
		argNo++
	}
	if len(notState) == 1 {
		rb = rb.Where(fmt.Sprintf("order_state != $%d", argNo), notState[0])
		argNo++
	} else if len(notState) > 1 {
		rb = rb.Where(fmt.Sprintf("order_state != ALL($%d)", argNo), pq.Int64Array(notState))
		argNo++
	}
	if usedMinCond {
		rb = rb.Where("used_chat_minute > 0")
	}
	if sec != 0 { // 多少秒内
		t := time.Now().Add(-time.Duration(sec) * time.Second)
		rb = rb.Where(fmt.Sprintf("create_time > $%d", argNo), t)
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) UpdateOrderState(ctx context.Context, orderId string, state int64, statesCond []int64) error {
	rb := squirrel.Update(m.table).Set("order_state", squirrel.Expr("$1", state)).Where("order_id = $2", orderId)

	if len(statesCond) > 0 {
		rb = rb.Where("order_state = ANY($3)", pq.Int64Array(statesCond))
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) ResetOrderAmount(ctx context.Context, orderId string, sumBuyMin int64, taxAmount int64) error {
	rb := squirrel.Update(m.table).Set("buy_minute_sum", squirrel.Expr("$1", sumBuyMin)).Set("tax_amount", squirrel.Expr("$2", taxAmount)).Where("order_id = $3", orderId).Where("order_state != $4", orderkey.ChatOrderStateSettle21)

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) ResetOrderAmount2(ctx context.Context, orderId string, sumBuyMin int64, taxAmount, platAmount, listenerAmount int64) error {
	rb := squirrel.Update(m.table).Set("buy_minute_sum", squirrel.Expr("$1", sumBuyMin)).Set("tax_amount", squirrel.Expr("$2", taxAmount)).Set("platform_share_amount", squirrel.Expr("$3", platAmount)).Set("listener_amount", squirrel.Expr("$4", listenerAmount)).Set("order_state", squirrel.Expr("$5", orderkey.ChatOrderStateSettle21)).Where("order_id = $6", orderId).Where("order_state = ANY($7)", pq.Int64Array{orderkey.ChatOrderStateSettle21, orderkey.ChatOrderStateAutoCommentFinish18, orderkey.ChatOrderStateRefundRefuseAutoFinish20,
		orderkey.ChatOrderStateAutoConfirmFinish19, orderkey.ChatOrderStateConfirmFinish16})

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) UpdateOrderPaySuccess(ctx context.Context, orderId string, state, buyMinSum int64, startTime, expiryTime *time.Time, statesCond []int64) error {
	rb := squirrel.Update(m.table).Set("order_state", squirrel.Expr("$1", state)).Set("expiry_time", squirrel.Expr("$2", expiryTime)).Set("buy_minute_sum", squirrel.Expr("$3", buyMinSum))
	argNo := 4
	if startTime != nil {
		rb = rb.Set("start_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), startTime))
		argNo++
	}
	rb = rb.Where(fmt.Sprintf("order_id = $%d", argNo), orderId)
	argNo++
	if len(statesCond) > 0 {
		rb = rb.Where(fmt.Sprintf("order_state = ANY($%d)", argNo), pq.Int64Array(statesCond))
		argNo++
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) Find(ctx context.Context, uid int64, listenerUid int64, orderType int64, state []int64, listType int64, sort string, pageNo, pageSize int64) ([]*ChatOrder, error) {
	rb := squirrel.Select(chatOrderRows).From(m.table)
	argNo := 1
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf("uid > $%d", argNo), userkey.UidStart)
		argNo++
	}
	if orderType != 0 {
		rb = rb.Where(fmt.Sprintf("order_type = $%d", argNo), orderType)
		argNo++
	}

	switch listType {
	case orderkey.OrderListTypeNeedFeedback:
		rb = rb.Where("order_state != ALL('{1,2,3,4,5,24}')").Where("feedback_time is null")

	default:
		if len(state) > 0 {
			rb = rb.Where(fmt.Sprintf("order_state = ANY($%d)", argNo), pq.Int64Array(state))
			argNo++
		} else if len(state) == 0 {
			rb = rb.Where("order_state != ALL('{1,2}')")
		}
	}

	if sort != "" {
		rb = rb.OrderBy(sort)
	} else {
		rb = rb.OrderBy("create_time desc")
	}

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query, values, err := rb.Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatOrder, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatOrderModel) FindAllPaidOrder(ctx context.Context, pageNo, pageSize int64) ([]*ChatOrder, error) {
	rb := squirrel.Select(chatOrderRows).From(m.table).Where("uid > $1", userkey.UidStart).Where("order_state > 2").Where("actual_amount > 0").OrderBy("create_time desc")

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query, values, err := rb.Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatOrder, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatOrderModel) UpdateOrderUse(ctx context.Context, orderId string, usedChatMinute int64, startTime, endTime *time.Time, state int64, statesCond []int64) error {
	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	rb := squirrel.Update(m.table)
	argNo := 1
	if usedChatMinute != 0 {
		rb = rb.Set("used_chat_minute", squirrel.Expr(fmt.Sprintf("$%d", argNo), usedChatMinute))
		argNo++
	}
	if startTime != nil {
		rb = rb.Set("start_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), startTime))
		argNo++
	}
	if endTime != nil {
		rb = rb.Set("end_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), endTime))
		argNo++
	}
	if state != 0 {
		rb = rb.Set("order_state", squirrel.Expr(fmt.Sprintf("$%d", argNo), state))
		argNo++
	}

	rb = rb.Where(fmt.Sprintf("order_id = $%d", argNo), orderId)
	argNo++

	if len(statesCond) > 0 {
		rb = rb.Where(fmt.Sprintf("order_state = ANY($%d)", argNo), pq.Int64Array(statesCond))
		argNo++
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) FindExpireOrder(ctx context.Context, expiryTimeStart, expiryTimeEnd *time.Time, orderType int64, state []int64, pageNo, pageSize int64) ([]*ExpireVoiceChatOrder, error) {
	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	query, values, err := squirrel.Select(expireVoiceChatOrderRows).From(m.table).Where("expiry_time between $1 and $2", expiryTimeStart, expiryTimeEnd).Where("order_type = $3", orderType).Where("order_state = ANY($4)", pq.Int64Array(state)).OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ExpireVoiceChatOrder, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
func (m *defaultChatOrderModel) SettleOrder(ctx context.Context, orderId string, shareAmount, amount, state int64, settleTime *time.Time) error {
	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set platform_share_amount = $1, listener_amount = $2, order_state = $3, settle_time = $4 where order_id = $5", m.table)
		return conn.ExecCtx(ctx, query, shareAmount, amount, state, settleTime, orderId)
	}, jakartaChatOrderOrderIdKey)
	return err
}

func (m *defaultChatOrderModel) UpdateOrderOpinionAndState(ctx context.Context, orderId string, newData *ChatOrder, statesCond []int64) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.Comment != "" || newData.Star != 0 {
		if newData.Comment != "" {
			rb = rb.Set("comment", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Comment))
			argNo++
		}
		rb = rb.Set("comment_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), time.Now())).Set("star", squirrel.Expr(fmt.Sprintf("$%d", argNo+1), newData.Star))
		argNo += 2
		if len(newData.CommentTag) > 0 {
			rb = rb.Set("comment_tag", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CommentTag))
			argNo++
		}
	}
	if newData.Reply != "" {
		rb = rb.Set("reply", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Reply)).Set("reply_time", squirrel.Expr(fmt.Sprintf("$%d", argNo+1), time.Now()))
		argNo += 2
	}
	if newData.Feedback != "" {
		rb = rb.Set("feedback", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Feedback)).Set("feedback_time", squirrel.Expr(fmt.Sprintf("$%d", argNo+1), time.Now()))
		argNo += 2
	}
	if newData.OrderState != 0 {
		rb = rb.Set("order_state", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.OrderState))
		argNo++
	}
	rb = rb.Where(fmt.Sprintf("order_id = $%d", argNo), orderId)
	argNo++

	if len(statesCond) > 0 {
		rb = rb.Where(fmt.Sprintf("order_state = ANY($%d)", argNo), pq.Int64Array(statesCond))
		argNo++
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}

	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) FindListenerOrderOpinionList(ctx context.Context, listenerUid int64, star int64, pageNo int64, pageSize int64) ([]*ListenerOrderOpinion, error) {
	resp := make([]*ListenerOrderOpinion, 0)
	var err error
	if pageSize == 1 {
		if star != 0 {
			query := fmt.Sprintf("select %s from %s where listener_uid = $1 and comment !='' and star = $2 order by comment_time desc limit 1", listenerOrderOpinionRows, m.table)
			err = m.QueryRowsNoCacheCtx(ctx, &resp, query, listenerUid, star)
		} else {
			query := fmt.Sprintf("select %s from %s where listener_uid = $1 and comment != '' order by comment_time desc limit 1", listenerOrderOpinionRows, m.table)
			err = m.QueryRowsNoCacheCtx(ctx, &resp, query, listenerUid)
		}
	} else {
		if star != 0 {
			query := fmt.Sprintf("select %s from %s where listener_uid = $1 and comment_time is not null and star = $2 order by comment_time desc limit $3 offset $4", listenerOrderOpinionRows, m.table)
			err = m.QueryRowsNoCacheCtx(ctx, &resp, query, listenerUid, star, pageSize, (pageNo-1)*pageSize)
		} else {
			query := fmt.Sprintf("select %s from %s where listener_uid = $1 and comment_time is not null order by comment_time desc limit $2 offset $3", listenerOrderOpinionRows, m.table)
			err = m.QueryRowsNoCacheCtx(ctx, &resp, query, listenerUid, pageSize, (pageNo-1)*pageSize)
		}
	}

	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatOrderModel) UpdateApplyRefund(ctx context.Context, orderId string, newData *ChatOrder, statesCond []int64) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.ApplyRefundTime.Valid {
		rb = rb.Set("apply_refund_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ApplyRefundTime))
		argNo++
	}
	if newData.RefundSuccessTime.Valid {
		rb = rb.Set("refund_success_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RefundSuccessTime))
		argNo++
	}
	if newData.RefundReason != "" {
		rb = rb.Set("refund_reason", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RefundReason))
		argNo++
	}
	if newData.RefundReasonTag != 0 {
		rb = rb.Set("refund_reason_tag", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RefundReasonTag))
		argNo++
	}
	if newData.OrderState != 0 {
		rb = rb.Set("order_state", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.OrderState))
		argNo++
	}
	if newData.Additional != "" {
		rb = rb.Set("additional", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Additional))
		argNo++
	}
	if newData.Attachment != "" {
		rb = rb.Set("attachment", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Attachment))
		argNo++
	}
	if newData.RefundCheckRemark != "" {
		rb = rb.Set("refund_check_remark", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RefundCheckRemark))
		argNo++
	}
	if argNo == 1 {
		return nil
	}

	rb = rb.Where(fmt.Sprintf("order_id = $%d", argNo), orderId)
	argNo++

	if len(statesCond) > 0 {
		rb = rb.Where(fmt.Sprintf("order_state = ANY($%d)", argNo), pq.Int64Array(statesCond))
		argNo++
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}

	jakartaChatOrderOrderIdKey := fmt.Sprintf("%s%v", cacheJakartaChatOrderOrderIdPrefix, orderId)
	var rs sql.Result
	rs, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatOrderOrderIdKey)
	if err != nil {
		return err
	}
	af, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if af <= 0 {
		return xerr.NewGrpcErrCodeMsg(xerr.DbError, "更新订单状态异常")
	}
	return nil
}

func (m *defaultChatOrderModel) FindNeedAutoProcessOrderList(ctx context.Context, state []int64, beforeTime *time.Time, pageNo, pageSize int64) ([]*AutoProcessOrder, error) {
	resp := make([]*AutoProcessOrder, 0)
	query := fmt.Sprintf("select %s from %s where update_time < $1 and order_state = ANY($2) order by create_time asc limit $3 offset $4", autoProcessOrderRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, beforeTime, pq.Int64Array(state), pageSize, (pageNo-1)*pageSize)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

const findPaidUserCountSQL2 = "select count(*) from (select uid from %s where listener_uid = $1 and order_state != ALL('{1, 2}') and create_time between $2 and $3 and uid > $4 group by uid) t"

func (m *defaultChatOrderModel) CountPaidUser(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error) {
	query := fmt.Sprintf(findPaidUserCountSQL2, m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid, start, end, userkey.UidStart)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const findPaidUserCountSQL = "select count(*) from (select uid from %s where uid > 0 and listener_uid = $1 and order_state != ALL('{1, 2}') group by uid) t"

func (m *defaultChatOrderModel) CountPaidUser2(ctx context.Context, listenerUid int64) (int64, error) {
	query := fmt.Sprintf(findPaidUserCountSQL, m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

//
const findRepeatPaidUserCountSQL = "select count(uid) from (select uid from %s where listener_uid = $1 and order_state != ALL('{1, 2}') and create_time between $2 and $3 and uid > $4 group by uid having count(uid) >= 2) t"

func (m *defaultChatOrderModel) CountRepeatPaidUser(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error) {
	query := fmt.Sprintf(findRepeatPaidUserCountSQL, m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid, start, end, userkey.UidStart)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const findRepeatPaidUserCountSQL2 = "select count(uid) from (select uid from %s where uid > 0 and listener_uid = $1 and order_state != ALL('{1, 2}') group by uid having count(uid) >= 2) t"

func (m *defaultChatOrderModel) CountRepeatPaidUser2(ctx context.Context, listenerUid int64) (int64, error) {
	query := fmt.Sprintf(findRepeatPaidUserCountSQL2, m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const findAveragePaidAmountPerDaySQL = "select coalesce(sum(actual_amount), 0) from %s where listener_uid = $1 and order_state != ALL('{1, 2}') and create_time between $2 and $3 and uid > $4"

func (m *defaultChatOrderModel) SumListenerPaidAmount(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error) {
	query := fmt.Sprintf(findAveragePaidAmountPerDaySQL, m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid, start, end, userkey.UidStart)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

const countPaidUserCntSQL1 = "select count(uid) from (select uid from %s where create_time between $1 and $2 and user_channel = $3 and order_state != ALL('{1, 2}') and uid > $4 group by uid having count(uid) >= 1) t"
const countPaidUserCntSQL2 = "select count(uid) from (select uid from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and uid > $3 group by uid having count(uid) >= 1) t"

func (m *defaultChatOrderModel) CountPaidUserCnt(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countPaidUserCntSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countPaidUserCntSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countPaidNewUserCntSQL1 = "select count(uid) from (select uid from %s where uid between $1 and $2 and create_time between $3 and $4 and user_channel = $5 and order_state != ALL('{1, 2}') group by uid having count(uid) >= 1) t"
const countPaidNewUserCntSQL2 = "select count(uid) from (select uid from %s where uid between $1 and $2 and create_time between $3 and $4 and order_state != ALL('{1, 2}') group by uid having count(uid) >= 1) t"

func (m *defaultChatOrderModel) CountPaidNewUserCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, uidStart, uidEnd int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countPaidNewUserCntSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, createTimeStart, createTimeEnd, channel)
	} else {
		query := fmt.Sprintf(countPaidNewUserCntSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, createTimeStart, createTimeEnd)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countPaidUserCntRangeCreateTimeSQL1 = "select count(uid) from (select uid from %s where create_time between $1 and $2 and user_channel = $3 and uid > $4 and order_state != ALL('{1, 2}') group by uid having count(uid) >= 1) t"
const countPaidUserCntRangeCreateTimeSQL2 = "select count(uid) from (select uid from %s where create_time between $1 and $2 and uid > $3 and order_state != ALL('{1, 2}') group by uid having count(uid) >= 1) t"

func (m *defaultChatOrderModel) CountPaidUserCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countPaidUserCntRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countPaidUserCntRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countPaidUserCntSQL3 = "select count(uid) from (select uid from %s where listener_uid = $1 and create_time between $2 and $3 and order_state != ALL('{1, 2}') group by uid having count(uid) >= 1) t"

func (m *defaultChatOrderModel) CountPaidUserCnt2(ctx context.Context, listenerUid int64, createTimeStart *time.Time, createTimeEnd *time.Time) (int64, error) {
	var resp int64
	var err error
	query := fmt.Sprintf(countPaidUserCntSQL3, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid, createTimeStart, createTimeEnd)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countPaidOrderCntSQL1 = "select count(order_id) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and user_channel = $3 and uid > $4"
const countPaidOrderCntSQL2 = "select count(order_id) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and uid > $3"

func (m *defaultChatOrderModel) CountPaidOrderCnt(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countPaidOrderCntSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countPaidOrderCntSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countNewUserPaidOrderCntSQL1 = "select count(order_id) from %s where uid between $1 and $2 and order_state != ALL('{1, 2}') and create_time between $3 and $4 and user_channel = $5"
const countNewUserPaidOrderCntSQL2 = "select count(order_id) from %s where uid between $1 and $2 and order_state != ALL('{1, 2}') and create_time between $3 and $4"

func (m *defaultChatOrderModel) CountNewUserPaidOrderCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, uidStart, uidEnd int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countNewUserPaidOrderCntSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, createTimeStart, createTimeEnd, channel)
	} else {
		query := fmt.Sprintf(countNewUserPaidOrderCntSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, createTimeStart, createTimeEnd)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countUserPaidOrderCntRangeCreateTimeSQL1 = "select count(order_id) from %s where order_state != ALL('{1, 2}') and create_time between $1 and $2 and user_channel = $3 and uid > $4"
const countUserPaidOrderCntRangeCreateTimeSQL2 = "select count(order_id) from %s where order_state != ALL('{1, 2}') and create_time between $1 and $2 and uid > $3"

func (m *defaultChatOrderModel) CountUserPaidOrderCntRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countUserPaidOrderCntRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countUserPaidOrderCntRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countPaidOrderCntSQL3 = "select count(order_id) from %s where listener_uid = $1 and order_state != ALL('{1, 2}') and uid > 0"

func (m *defaultChatOrderModel) CountPaidOrderCnt2(ctx context.Context, listenerUid int64) (int64, error) {
	var resp int64
	var err error
	query := fmt.Sprintf(countPaidOrderCntSQL3, m.table)
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumPaidAmountSQL1 = "select coalesce(sum(actual_amount), 0) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and user_channel = $3 and uid > $4"
const sumPaidAmountSQL2 = "select coalesce(sum(actual_amount), 0) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and uid > $3"

func (m *defaultChatOrderModel) SumPaidAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumPaidAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumPaidAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumApplyRefundAmountSQL1 = "select coalesce(sum(actual_amount), 0) from %s where order_state = ANY($1) and apply_refund_time between $2 and $3 and user_channel = $4 and uid > $5"
const sumApplyRefundAmountSQL2 = "select coalesce(sum(actual_amount), 0) from %s where order_state = ANY($1) and apply_refund_time between $2 and $3 and uid > $4"

func (m *defaultChatOrderModel) SumApplyRefundAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumApplyRefundAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, pq.Int64Array(orderkey.ChatOrderRefundState), createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumApplyRefundAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, pq.Int64Array(orderkey.ChatOrderRefundState), createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumRefundSuccessAmountSQL1 = "select coalesce(sum(actual_amount), 0) from %s where order_state = $1 and refund_success_time between $2 and $3 and user_channel = $4 and uid > $5"
const sumRefundSuccessAmountSQL2 = "select coalesce(sum(actual_amount), 0) from %s where order_state = $1 and refund_success_time between $2 and $3 and uid > $4"

func (m *defaultChatOrderModel) SumRefundSuccessAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumRefundSuccessAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateFinishRefund8, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumRefundSuccessAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateFinishRefund8, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumConfirmAmountSQL1 = "select coalesce(sum(actual_amount), 0) from %s where order_state = $1 and settle_time between $2 and $3 and user_channel = $4 and uid > $5"
const sumConfirmAmountSQL2 = "select coalesce(sum(actual_amount), 0) from %s where order_state = $1 and settle_time between $2 and $3 and uid > $4"

func (m *defaultChatOrderModel) SumConfirmAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumConfirmAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateSettle21, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumConfirmAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateSettle21, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumListenerAmountSQL1 = "select coalesce(sum(listener_amount), 0) from %s where order_state = $1 and settle_time between $2 and $3 and user_channel = $4 and uid > $5"
const sumListenerAmountSQL2 = "select coalesce(sum(listener_amount), 0) from %s where order_state = $1 and settle_time between $2 and $3 and uid > $4"

func (m *defaultChatOrderModel) SumListenerAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumListenerAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateSettle21, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumListenerAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateSettle21, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumPlatformAmountSQL1 = "select coalesce(sum(platform_share_amount), 0) from %s where order_state = $1 and settle_time between $2 and $3 and user_channel = $4 and uid > $5"
const sumPlatformAmountSQL2 = "select coalesce(sum(platform_share_amount), 0) from %s where order_state = $1 and settle_time between $2 and $3 and uid > $4"

func (m *defaultChatOrderModel) SumPlatformAmount(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumPlatformAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateSettle21, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumPlatformAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, orderkey.ChatOrderStateSettle21, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const queryGoodCommentSQL = `select DISTINCT(uid),order_id, listener_uid,uid, listener_nick_name, nick_name,comment_time,comment_tag from %s where star = 5 and comment != '' order by comment_time desc limit $1 offset $2`

func (m *defaultChatOrderModel) FindGoodComment(ctx context.Context, pageNo, pageSize int64) ([]*ListenerGoodComment, error) {
	resp := make([]*ListenerGoodComment, 0)
	query := fmt.Sprintf(queryGoodCommentSQL, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, pageSize, (pageNo-1)*pageSize)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatOrderModel) CountOrderRatingCnt(ctx context.Context, listenerUid int64, star int64) (int64, error) {
	rb := squirrel.Select("count(order_id)").From(m.table).Where("listener_uid = $1", listenerUid).Where("uid > 0")
	if star != 0 {
		rb = rb.Where("star = $2", star)
	} else {
		rb = rb.Where("star != 0")
	}
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) SumOrderBuyUnitCnt(ctx context.Context, listenerUid int64) (int64, error) {
	rb := squirrel.Select("coalesce(sum(buy_unit), 0)").From(m.table).Where("uid > 0").Where("listener_uid = $1", listenerUid).Where("order_state > 2").Where("actual_amount > 0")
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) SumUserOrderBuyUnitCnt(ctx context.Context, uid, listenerUid int64) (int64, error) {
	rb := squirrel.Select("coalesce(sum(buy_unit), 0)").From(m.table).Where("uid = $1", uid).Where("listener_uid = $2", listenerUid).Where("order_state > 2").Where("actual_amount > 0")
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) SumUserOrderBuyUnitCntBefore(ctx context.Context, uid, listenerUid int64, endCreateTime *time.Time, notOrderId string) (int64, error) {
	rb := squirrel.Select("coalesce(sum(buy_unit), 0)").From(m.table).Where("uid = $1", uid).Where("listener_uid = $2", listenerUid).Where("order_state > 2").Where("actual_amount > 0").Where("create_time < $3", endCreateTime).Where("order_id != $4", notOrderId)
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) SumOrderBuyUnitCntByState(ctx context.Context, uid, listenerUid, state int64) (int64, error) {
	rb := squirrel.Select("coalesce(sum(buy_unit), 0)").From(m.table).Where("uid = $1", uid).Where("listener_uid = $2", listenerUid).Where("order_state = $3", state)
	query, args, err := rb.ToSql()
	if err != nil {
		return 0, err
	}
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const findOrderCountSQL = `select count(order_id) from %s where listener_uid = $1 and order_state = ANY($2) and uid > 0`

func (m *defaultChatOrderModel) CountOrder(ctx context.Context, listenerUid int64, state []int64) (int64, error) {
	query := fmt.Sprintf(findOrderCountSQL, m.table)
	var resp int64
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, listenerUid, pq.Int64Array(state))
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) FindAllComment(ctx context.Context, listenerUid int64, pageNo, pageSize int64) ([]*ListenerGoodComment, error) {
	query, args, err := squirrel.Select(listenerGoodCommentRows).From(m.table).Where("listener_uid = $1", listenerUid).Where("comment_time is not null").OrderBy("comment_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	resp := make([]*ListenerGoodComment, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatOrderModel) FindLastCommentOrder(ctx context.Context, listenerUid int64, star int64) (*ListenerGoodComment, error) {
	query, args, err := squirrel.Select(listenerGoodCommentRows).From(m.table).Where("listener_uid = $1", listenerUid).Where("star = $2", star).Where("comment_time is not null").OrderBy("comment_time DESC").Limit(1).ToSql()
	if err != nil {
		return nil, err
	}
	var resp ListenerGoodComment
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

const countNewUserOrderSQL = `select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and order_state = ANY($5)`
const countNewUserOrderSQL2 = `select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and order_state = ANY($5) and user_channel = $6`

func (m *defaultChatOrderModel) CountNewUserOrderRangeCreateTime(ctx context.Context, startTime, endTime *time.Time, uidStart, uidEnd int64, channel string, state []int64) (int64, error) {

	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countNewUserOrderSQL, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, startTime, endTime, pq.Int64Array(state), channel)
	} else {
		query := fmt.Sprintf(countNewUserOrderSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, startTime, endTime, pq.Int64Array(state))
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countUserOrderRangeCreateTimeSQL1 = `select count(order_id) from %s where create_time between $1 and $2 and order_state = ANY($3) and uid > $4`
const countUserOrderRangeCreateTimeSQL2 = `select count(order_id) from %s where create_time between $1 and $2 and order_state = ANY($3) and user_channel = $4 and uid > $5`

func (m *defaultChatOrderModel) CountUserOrderRangeCreateTime(ctx context.Context, startTime, endTime *time.Time, channel string, state []int64) (int64, error) {

	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countUserOrderRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, startTime, endTime, pq.Int64Array(state), channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countUserOrderRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, startTime, endTime, pq.Int64Array(state), userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countRepeatPaidNewUserSQL = "select count(uid) from (select uid from %s where uid between $1 and $2 and order_state != ALL('{1, 2}') and create_time  between $3 and $4 and user_channel = $5 group by uid having count(uid) >= 2) t"
const countRepeatPaidNewUserSQL2 = "select count(uid) from (select uid from %s where uid between $1 and $2 and order_state != ALL('{1, 2}') and create_time between $3 and $4 group by uid having count(uid) >= 2) t"

func (m *defaultChatOrderModel) CountRepeatPaidNewUserRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd int64, channel string) (int64, error) {

	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countRepeatPaidNewUserSQL, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end, channel)
	} else {
		query := fmt.Sprintf(countRepeatPaidNewUserSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countRepeatPaidUserRangeCreateTimeSQL1 = "select count(uid) from (select uid from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and user_channel = $3 and uid > $4 group by uid having count(uid) >= 2) t"
const countRepeatPaidUserRangeCreateTimeSQL2 = "select count(uid) from (select uid from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and uid > $3 group by uid having count(uid) >= 2) t"

func (m *defaultChatOrderModel) CountRepeatPaidUserRangeCreateTime(ctx context.Context, start, end *time.Time, channel string) (int64, error) {

	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countRepeatPaidUserRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countRepeatPaidUserRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countNewUserCommentOrderSQL = "select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and comment_time is not null and user_channel = $5"
const countNewUserCommentOrderSQL2 = "select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and comment_time is not null"

func (m *defaultChatOrderModel) CountNewUserCommentOrderRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countNewUserCommentOrderSQL, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end, channel)
	} else {
		query := fmt.Sprintf(countNewUserCommentOrderSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countUserCommentOrderRangeCreateTimeSQL1 = "select count(order_id) from %s where create_time between $1 and $2 and comment_time is not null and user_channel = $3 and uid > $4"
const countUserCommentOrderRangeCreateTimeSQL2 = "select count(order_id) from %s where create_time between $1 and $2 and comment_time is not null and uid > $3"

func (m *defaultChatOrderModel) CountUserCommentOrderRangeCreateTime(ctx context.Context, start, end *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countUserCommentOrderRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countUserCommentOrderRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countNewUserCommentOrderByStarSQL = "select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and star = $5 and user_channel = $6"
const countNewUserCommentOrderByStarSQL2 = "select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and star = $5"

func (m *defaultChatOrderModel) CountNewUserCommentOrderByStarRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd, star int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countNewUserCommentOrderByStarSQL, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end, star, channel)
	} else {
		query := fmt.Sprintf(countNewUserCommentOrderByStarSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end, star)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countUserCommentOrderByStarRangeCreateTimeSQL1 = "select count(order_id) from %s where create_time between $1 and $2 and star = $3 and user_channel = $4 and uid > $5"
const countUserCommentOrderByStarRangeCreateTimeSQL2 = "select count(order_id) from %s where create_time between $1 and $2 and star = $3 and uid > $4"

func (m *defaultChatOrderModel) CountUserCommentOrderByStarRangeCreateTime(ctx context.Context, start, end *time.Time, star int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countUserCommentOrderByStarRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, star, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countUserCommentOrderByStarRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, star, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countNewUserOrderCntByTypeSQL = "select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and order_type = $5 and user_channel = $6 and order_state != ALL('{1, 2}')"
const countNewUserOrderCntByTypeSQL2 = "select count(order_id) from %s where uid between $1 and $2 and create_time between $3 and $4 and order_type = $5 and order_state != ALL('{1, 2}')"

func (m *defaultChatOrderModel) CountNewUserOrderCntByTypeRangeCreateTime(ctx context.Context, start, end *time.Time, uidStart, uidEnd, orderType int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countNewUserOrderCntByTypeSQL, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end, orderType, channel)
	} else {
		query := fmt.Sprintf(countNewUserOrderCntByTypeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, start, end, orderType)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const countUserOrderCntByTypeRangeCreateTimeSQL1 = "select count(order_id) from %s where create_time between $1 and $2 and order_type = $3 and user_channel = $4 and order_state != ALL('{1, 2}') and uid > $5"
const countUserOrderCntByTypeRangeCreateTimeSQL2 = "select count(order_id) from %s where create_time between $1 and $2 and order_type = $3 and order_state != ALL('{1, 2}') and uid > $4"

func (m *defaultChatOrderModel) CountUserOrderCntByTypeRangeCreateTime(ctx context.Context, start, end *time.Time, orderType int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(countUserOrderCntByTypeRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, orderType, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(countUserOrderCntByTypeRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, start, end, orderType, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumNewUserPaidAmountSQL1 = "select coalesce(sum(actual_amount), 0) from %s where uid between $1 and $2 and create_time between $3 and $4 and order_state != ALL('{1, 2}') and user_channel = $5"
const sumNewUserPaidAmountSQL2 = "select coalesce(sum(actual_amount), 0) from %s where uid between $1 and $2 and create_time between $3 and $4 and order_state != ALL('{1, 2}')"

func (m *defaultChatOrderModel) SumNewUserPaidAmountRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, uidStart, uidEnd int64, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumNewUserPaidAmountSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, createTimeStart, createTimeEnd, channel)
	} else {
		query := fmt.Sprintf(sumNewUserPaidAmountSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, uidStart, uidEnd, createTimeStart, createTimeEnd)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumUserPaidAmountRangeCreateTimeLtvSQL1 = "select coalesce(sum(actual_amount), 0) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and user_channel = $3 and uid = ANY($4)"
const sumUserPaidAmountRangeCreateTimeLtvSQL2 = "select coalesce(sum(actual_amount), 0) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and uid = ANY($3)"

func (m *defaultChatOrderModel) SumUserPaidAmountRangeCreateTimeLtv(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string, uids []int64) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumUserPaidAmountRangeCreateTimeLtvSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, pq.Int64Array(uids))
	} else {
		query := fmt.Sprintf(sumUserPaidAmountRangeCreateTimeLtvSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, pq.Int64Array(uids))
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

const sumUserPaidAmountRangeCreateTimeSQL1 = "select coalesce(sum(actual_amount), 0) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and user_channel = $3 and uid > $4"
const sumUserPaidAmountRangeCreateTimeSQL2 = "select coalesce(sum(actual_amount), 0) from %s where create_time between $1 and $2 and order_state != ALL('{1, 2}') and uid > $3"

func (m *defaultChatOrderModel) SumUserPaidAmountRangeCreateTime(ctx context.Context, createTimeStart *time.Time, createTimeEnd *time.Time, channel string) (int64, error) {
	var resp int64
	var err error
	if channel != "" {
		query := fmt.Sprintf(sumUserPaidAmountRangeCreateTimeSQL1, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, channel, userkey.UidStart)
	} else {
		query := fmt.Sprintf(sumUserPaidAmountRangeCreateTimeSQL2, m.table)
		err = m.QueryRowNoCacheCtx(ctx, &resp, query, createTimeStart, createTimeEnd, userkey.UidStart)
	}
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultChatOrderModel) FindChatOrderFeedback(ctx context.Context, uid, pageNo, pageSize int64) ([]*ListenerFeedback, error) {
	q, a, err := squirrel.Select(listenerFeedbackRows).From(m.table).Where("uid = $1", uid).Where("feedback != ''").OrderBy("create_time desc").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	resp := make([]*ListenerFeedback, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, q, a...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
