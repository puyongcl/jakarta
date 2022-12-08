package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/userkey"
	"jakarta/common/tool"
	"strings"
	"time"
)

var _ ListenerProfileModel = (*customListenerProfileModel)(nil)

// 增加统计值
type AddListenerStat struct {
	ListenerUid          int64 `json:"listenerUid"`
	AddUserCount         int64 `json:"addUserCount"`
	AddChatDuration      int64 `json:"addChatDuration"`
	AddRatingSum         int64 `json:"addRatingSum"`
	AddFiveStar          int64 `json:"addFiveStar"`
	AddThreeStar         int64 `json:"addThreeStar"`
	AddOneStar           int64 `json:"addOneStar"`
	AddRefundOrderCnt    int64 `json:"addRefundOrderCnt"`
	AddFinishOrderCnt    int64 `json:"addFinishOrderCnt"`
	AddPaidOrderCnt      int64 `json:"addPaidOrderCnt"`
	AddRepeatPaidUserCnt int64 `json:"addRepeatPaidUserCnt"`
}

type (
	// ListenerProfileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerProfileModel.
	ListenerProfileModel interface {
		listenerProfileModel
		FindRecommendListenerList(ctx context.Context, pageNo int64, pageSize int64, specialties int64, chatType int64, gender int64, age int64, workstate []int64, onlineState []int64, sortorder int64, blkUid []int64) ([]*ListenerProfile, error)
		UpdateWorkState(ctx context.Context, session sqlx.Session, data *ListenerProfile) error
		UpdateListenerStat(ctx context.Context, addData *AddListenerStat) error
		ResetListenerStat(ctx context.Context, addData *AddListenerStat) error
		UpdateNoNeedCheckProfile(ctx context.Context, session sqlx.Session, data *ListenerProfile) error
		InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerProfile) (sql.Result, error)
		UpdateTrans(ctx context.Context, session sqlx.Session, newData *ListenerProfile) error
		FindListenerUidRangeUpdateTime(ctx context.Context, pageNo, pageSize, worksState int64, start, end *time.Time) ([]*ListenerShortProfile, error)
		UpdateOnlineState(ctx context.Context, listenerUid, state int64) error
		FindListenerRangeCreateTime(ctx context.Context, pageNo, pageSize int64, start, end *time.Time) ([]*ListenerProfile, error)
		FindRecentActive(ctx context.Context, interval int, limit uint64, state []int64) ([]int64, error)
		CountListenerByWorkState(ctx context.Context, workState int64) (int64, error)
		CountListenerByOnlineState(ctx context.Context, onlineState int64) (int64, error)
		FindAllListener(ctx context.Context, pageNo, pageSize int64) ([]int64, error)
	}

	customListenerProfileModel struct {
		*defaultListenerProfileModel
	}
)

// NewListenerProfileModel returns a model for the database table.
func NewListenerProfileModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerProfileModel {
	return &customListenerProfileModel{
		defaultListenerProfileModel: newListenerProfileModel(conn, c),
	}
}

type ListenerShortProfile struct {
	ListenerUid int64  `db:"listener_uid"`
	NickName    string `db:"nick_name"`
	Avatar      string `db:"avatar"`
}

func (m *defaultListenerProfileModel) FindRecommendListenerList(ctx context.Context, pageNo int64, pageSize int64,
	specialties int64, chatType int64, gender int64, age int64, workstate []int64, onlineState []int64, sortorder int64, blkUid []int64) ([]*ListenerProfile, error) {
	rb := squirrel.Select(listenerProfileRows).From(m.table)
	argNo := 1
	if len(blkUid) > 0 {
		cond := fmt.Sprintf("listener_uid != ALL($%d)", argNo)
		rb = rb.Where(cond, pq.Int64Array(blkUid))
		argNo++
	}
	if chatType != 0 {
		if chatType == orderkey.ListenerOrderTypeTextChat {
			rb = rb.Where(fmt.Sprintf("text_chat_switch = %d", db.Enable))
		} else if chatType == orderkey.ListenerOrderTypeVoiceChat {
			rb = rb.Where(fmt.Sprintf("voice_chat_switch = %d", db.Enable))
		}
	}
	if gender != 0 {
		cond := fmt.Sprintf("gender = $%d", argNo)
		rb = rb.Where(cond, gender)
		argNo++
	}
	if age != 0 { //
		switch age {
		case listenerkey.AgeRange1:
			rb = rb.Where("birthday >= '1960-01-01'::date AND birthday <= '1969-12-31'::date")
		case listenerkey.AgeRange2:
			rb = rb.Where("birthday >= '1970-01-01'::date AND birthday <= '1979-12-31'::date")
		case listenerkey.AgeRange3:
			rb = rb.Where("birthday >= '1980-01-01'::date AND birthday <= '1989-12-31'::date")
		case listenerkey.AgeRange4:
			rb = rb.Where("birthday >= '1990-01-01'::date AND birthday <= '1999-12-31'::date")
		default:
		}
	}
	if len(workstate) != 0 {
		cond := fmt.Sprintf("work_state = ANY($%d)", argNo)
		rb = rb.Where(cond, pq.Int64Array(workstate))
		argNo++
	} else {
		rb = rb.Where(fmt.Sprintf("work_state != %d", listenerkey.ListenerWorkStateAccountDeleted))
	}
	if specialties != 0 {
		cond := fmt.Sprintf("$%d = ANY(specialties)", argNo)
		rb = rb.Where(cond, specialties)
		argNo++
	}
	if len(onlineState) != 0 {
		cond := fmt.Sprintf("online_state = ANY($%d)", argNo)
		rb = rb.Where(cond, pq.Int64Array(onlineState))
		argNo++
	}
	if sortorder != 0 { // TODO 综合排序 回头客 需要再规划
		switch sortorder {
		case listenerkey.ListenerSortOrderDefault:
			rb = rb.OrderByClause("listener_uid ASC")
		case listenerkey.ListenerSortOrderRatingStar:
			rb = rb.OrderByClause("five_star DESC")
		case listenerkey.ListenerSortOrderRepeatCustomer:
			rb = rb.OrderBy("repeat_paid_user_cnt DESC")
		case listenerkey.ListenerSortOrderChatMinute:
			rb = rb.OrderBy("chat_duration DESC")
		}
	} else {
		rb = rb.OrderBy("repeat_paid_user_cnt DESC")
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

	resp := make([]*ListenerProfile, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerProfileModel) UpdateWorkState(ctx context.Context, session sqlx.Session, newData *ListenerProfile) error {
	data, err := m.FindOne(ctx, newData.ListenerUid)
	if err != nil {
		return err
	}
	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.WorkState != 0 && data.WorkState != newData.WorkState {
		rb = rb.Set("work_state", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.WorkState))
		argNo++
	}
	if newData.RestingTimeEnable != 0 && data.RestingTimeEnable != newData.RestingTimeEnable {
		rb = rb.Set("resting_time_enable", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.RestingTimeEnable))
		argNo++
	}
	if newData.StartWorkTime != "" && data.StartWorkTime != newData.StartWorkTime {
		rb = rb.Set("start_work_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.StartWorkTime))
		argNo++
	}
	if newData.StopWorkTime != "" && data.StopWorkTime != newData.StopWorkTime {
		rb = rb.Set("stop_work_time", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.StopWorkTime))
		argNo++
	}
	if len(newData.WorkDays) > 0 && !tool.IsEqualArrayInt64(data.WorkDays, newData.WorkDays) {
		rb = rb.Set("work_days", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.WorkDays))
		argNo++
	}

	if argNo == 1 {
		return nil
	}

	query, args, err := rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), newData.ListenerUid).ToSql()
	if err != nil {
		return err
	}
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, newData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) UpdateListenerStat(ctx context.Context, addData *AddListenerStat) error {
	rb := squirrel.Update(m.table)
	argNo := 1
	if addData.AddUserCount != 0 {
		rb = rb.Set("user_count", squirrel.Expr(fmt.Sprintf("user_count + $%d", argNo), addData.AddUserCount))
		argNo++
	}
	if addData.AddChatDuration != 0 {
		rb = rb.Set("chat_duration", squirrel.Expr(fmt.Sprintf("chat_duration + $%d", argNo), addData.AddChatDuration))
		argNo++
	}
	if addData.AddRatingSum != 0 {
		rb = rb.Set("rating_sum", squirrel.Expr(fmt.Sprintf("rating_sum + $%d", argNo), addData.AddRatingSum))
		argNo++
	}
	if addData.AddFiveStar != 0 {
		rb = rb.Set("five_star", squirrel.Expr(fmt.Sprintf("five_star + $%d", argNo), addData.AddFiveStar))
		argNo++
	}
	if addData.AddThreeStar != 0 {
		rb = rb.Set("three_star", squirrel.Expr(fmt.Sprintf("three_star + $%d", argNo), addData.AddThreeStar))
		argNo++
	}
	if addData.AddOneStar != 0 {
		rb = rb.Set("one_star", squirrel.Expr(fmt.Sprintf("one_star + $%d", argNo), addData.AddOneStar))
		argNo++
	}
	if addData.AddRefundOrderCnt != 0 {
		rb = rb.Set("refund_order_cnt", squirrel.Expr(fmt.Sprintf("refund_order_cnt + $%d", argNo), addData.AddRefundOrderCnt))
		argNo++
	}
	if addData.AddFinishOrderCnt != 0 {
		rb = rb.Set("finish_order_cnt", squirrel.Expr(fmt.Sprintf("finish_order_cnt + $%d", argNo), addData.AddFinishOrderCnt))
		argNo++
	}
	if addData.AddPaidOrderCnt != 0 {
		rb = rb.Set("paid_order_cnt", squirrel.Expr(fmt.Sprintf("paid_order_cnt + $%d", argNo), addData.AddPaidOrderCnt))
		argNo++
	}
	if addData.AddRepeatPaidUserCnt != 0 {
		rb = rb.Set("repeat_paid_user_cnt", squirrel.Expr(fmt.Sprintf("repeat_paid_user_cnt + $%d", argNo), addData.AddRepeatPaidUserCnt))
		argNo++
	}

	if argNo == 1 {
		return nil
	}

	rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), addData.ListenerUid)

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, addData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) ResetListenerStat(ctx context.Context, addData *AddListenerStat) error {
	rb := squirrel.Update(m.table)
	argNo := 1

	rb = rb.Set("user_count", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddUserCount))
	argNo++

	rb = rb.Set("chat_duration", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddChatDuration))
	argNo++

	rb = rb.Set("rating_sum", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddRatingSum))
	argNo++

	rb = rb.Set("five_star", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddFiveStar))
	argNo++

	rb = rb.Set("three_star", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddThreeStar))
	argNo++

	rb = rb.Set("one_star", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddOneStar))
	argNo++

	rb = rb.Set("refund_order_cnt", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddRefundOrderCnt))
	argNo++

	rb = rb.Set("finish_order_cnt", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddFinishOrderCnt))
	argNo++

	rb = rb.Set("paid_order_cnt", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddPaidOrderCnt))
	argNo++

	rb = rb.Set("repeat_paid_user_cnt", squirrel.Expr(fmt.Sprintf("$%d", argNo), addData.AddRepeatPaidUserCnt))
	argNo++

	rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), addData.ListenerUid)

	query, args, err := rb.ToSql()
	if err != nil {
		return err
	}
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, addData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) UpdateNoNeedCheckProfile(ctx context.Context, session sqlx.Session, newData *ListenerProfile) error {
	data, err := m.FindOne(ctx, newData.ListenerUid)
	if err != nil {
		return err
	}

	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.Job != data.Job {
		rb = rb.Set("job", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Job))
		argNo++
	}
	if newData.Education != data.Education {
		rb = rb.Set("education", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Education))
		argNo++
	}
	if newData.MaritalStatus != data.MaritalStatus {
		rb = rb.Set("marital_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.MaritalStatus))
		argNo++
	}
	if newData.TextChatSwitch != data.TextChatSwitch {
		rb = rb.Set("text_chat_switch", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TextChatSwitch))
		argNo++
	}
	if newData.VoiceChatSwitch != data.VoiceChatSwitch {
		rb = rb.Set("voice_chat_switch", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.VoiceChatSwitch))
		argNo++
	}

	query, args, err := rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), newData.ListenerUid).ToSql()
	if err != nil {
		return err
	}

	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, newData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerProfile) (sql.Result, error) {
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54)", m.table, listenerProfileRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.ListenerUid, data.NickName, data.ListenerName, data.OpenId, data.Avatar, data.MaritalStatus, data.PhoneNumber, data.Constellation, data.Province, data.City, data.Job, data.Education, data.Gender, data.Birthday, data.IdNo, data.IdPhoto1, data.IdPhoto2, data.IdPhoto3, data.Specialties, data.Introduction, data.VoiceFile, data.Experience1, data.Experience2, data.CertType, data.OtherPlatformAccount, data.CertFiles1, data.CertFiles2, data.CertFiles3, data.CertFiles4, data.CertFiles5, data.AutoReplyNew, data.AutoReplyProcessing, data.AutoReplyFinish, data.TextChatPrice, data.VoiceChatPrice, data.TextChatSwitch, data.VoiceChatSwitch, data.UserCount, data.ChatDuration, data.RatingSum, data.FiveStar, data.ThreeStar, data.OneStar, data.RestingTimeEnable, data.StopWorkTime, data.StartWorkTime, data.WorkDays, data.WorkState, data.RefundOrderCnt, data.FinishOrderCnt, data.PaidOrderCnt, data.RepeatPaidUserCnt, data.OnlineState, data.Channel)
	}, jakartaListenerProfileListenerUidKey)
	return ret, err
}

func (m *defaultListenerProfileModel) UpdateTrans(ctx context.Context, session sqlx.Session, newData *ListenerProfile) error {
	data, err := m.FindOne(ctx, newData.ListenerUid)
	if err != nil {
		return err
	}

	rb := squirrel.Update(m.table)
	argNo := 1
	if newData.NickName != data.NickName {
		rb = rb.Set("nick_name", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.NickName))
		argNo++
	}
	if newData.ListenerName != data.ListenerName {
		rb = rb.Set("listener_name", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.ListenerName))
		argNo++
	}
	if newData.Avatar != data.Avatar {
		rb = rb.Set("avatar", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Avatar))
		argNo++
	}
	if newData.MaritalStatus != data.MaritalStatus {
		rb = rb.Set("marital_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.MaritalStatus))
		argNo++
	}
	if newData.PhoneNumber != data.PhoneNumber {
		rb = rb.Set("phone_number", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.PhoneNumber))
		argNo++
	}
	if newData.Constellation != data.Constellation {
		rb = rb.Set("constellation", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Constellation))
		argNo++
	}
	if newData.Province != data.Province {
		rb = rb.Set("province", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Province))
		argNo++
	}
	if newData.City != data.City {
		rb = rb.Set("city", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.City))
		argNo++
	}
	if newData.Job != data.Job {
		rb = rb.Set("job", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Job))
		argNo++
	}
	if newData.Education != data.Education {
		rb = rb.Set("education", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Education))
		argNo++
	}
	if newData.Gender != data.Gender {
		rb = rb.Set("gender", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Gender))
		argNo++
	}
	if newData.Birthday != data.Birthday {
		rb = rb.Set("birthday", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Birthday))
		argNo++
	}
	if newData.IdNo != data.IdNo {
		rb = rb.Set("id_no", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.IdNo))
		argNo++
	}
	if newData.IdPhoto1 != data.IdPhoto1 {
		rb = rb.Set("id_photo1", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.IdPhoto1))
		argNo++
	}
	if newData.IdPhoto2 != data.IdPhoto2 {
		rb = rb.Set("id_photo2", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.IdPhoto2))
		argNo++
	}
	if newData.IdPhoto3 != data.IdPhoto3 {
		rb = rb.Set("id_photo3", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.IdPhoto3))
		argNo++
	}
	if !tool.IsEqualArrayInt64(newData.Specialties, data.Specialties) {
		rb = rb.Set("specialties", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Specialties))
		argNo++
	}
	if newData.Introduction != data.Introduction {
		rb = rb.Set("introduction", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Introduction))
		argNo++
	}
	if newData.VoiceFile != data.VoiceFile {
		rb = rb.Set("voice_file", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.VoiceFile))
		argNo++
	}
	if newData.Experience1 != data.Experience1 {
		rb = rb.Set("experience1", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Experience1))
		argNo++
	}
	if newData.Experience2 != data.Experience2 {
		rb = rb.Set("experience2", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.Experience2))
		argNo++
	}
	if newData.CertType != data.CertType {
		rb = rb.Set("cert_type", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CertType))
		argNo++
	}
	if newData.OtherPlatformAccount != data.OtherPlatformAccount {
		rb = rb.Set("other_platform_account", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.OtherPlatformAccount))
		argNo++
	}
	if newData.CertFiles1 != data.CertFiles1 {
		rb = rb.Set("cert_files1", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CertFiles1))
		argNo++
	}
	if newData.CertFiles2 != data.CertFiles2 {
		rb = rb.Set("cert_files2", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CertFiles2))
		argNo++
	}
	if newData.CertFiles3 != data.CertFiles3 {
		rb = rb.Set("cert_files3", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CertFiles3))
		argNo++
	}
	if newData.CertFiles4 != data.CertFiles4 {
		rb = rb.Set("cert_files4", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CertFiles4))
		argNo++
	}
	if newData.CertFiles5 != data.CertFiles5 {
		rb = rb.Set("cert_files5", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CertFiles5))
		argNo++
	}
	if newData.AutoReplyNew != data.AutoReplyNew {
		rb = rb.Set("auto_reply_new", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.AutoReplyNew))
		argNo++
	}
	if newData.AutoReplyProcessing != data.AutoReplyProcessing {
		rb = rb.Set("auto_reply_processing", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.AutoReplyProcessing))
		argNo++
	}
	if newData.AutoReplyFinish != data.AutoReplyFinish {
		rb = rb.Set("auto_reply_finish", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.AutoReplyFinish))
		argNo++
	}
	if newData.TextChatPrice != data.TextChatPrice {
		rb = rb.Set("text_chat_price", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TextChatPrice))
		argNo++
	}
	if newData.VoiceChatPrice != data.VoiceChatPrice {
		rb = rb.Set("voice_chat_price", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.VoiceChatPrice))
		argNo++
	}
	if newData.TextChatSwitch != data.TextChatSwitch {
		rb = rb.Set("text_chat_switch", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.TextChatSwitch))
		argNo++
	}
	if newData.VoiceChatSwitch != data.VoiceChatSwitch {
		rb = rb.Set("voice_chat_switch", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.VoiceChatSwitch))
		argNo++
	}

	query, args, err := rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), newData.ListenerUid).ToSql()
	if err != nil {
		return err
	}

	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, newData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) FindListenerUidRangeUpdateTime(ctx context.Context, pageNo, pageSize, workState int64, start, end *time.Time) ([]*ListenerShortProfile, error) {
	listenerShortProfileFieldNames := builder.RawFieldNames(&ListenerShortProfile{}, true)
	listenerShortProfileRows := strings.Join(listenerShortProfileFieldNames, ",")
	rb := squirrel.Select(listenerShortProfileRows).From(m.table)
	argNo := 1
	if workState != 0 {
		rb = rb.Where(fmt.Sprintf("work_state = $%d", argNo), workState)
		argNo++
	}
	if start != nil && end != nil {
		rb = rb.Where(fmt.Sprintf("update_time between $%d and $%d", argNo, argNo+1), start, end)
		argNo += 2
	} else if start != nil {
		rb = rb.Where(fmt.Sprintf("update_time > $%d", argNo), start)
		argNo++
	} else if end != nil {
		rb = rb.Where(fmt.Sprintf("update_time < $%d", argNo), end)
		argNo++
	}

	query, args, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ListenerShortProfile, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

const updateListenerOnlineStateSQL = "update %s set online_state = $1 where listener_uid = $2"

func (m *defaultListenerProfileModel) UpdateOnlineState(ctx context.Context, listenerUid, state int64) error {
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf(updateListenerOnlineStateSQL, m.table)
		return conn.ExecCtx(ctx, query, state, listenerUid)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) FindListenerRangeCreateTime(ctx context.Context, pageNo, pageSize int64, start, end *time.Time) ([]*ListenerProfile, error) {
	rb := squirrel.Select(listenerProfileRows).From(m.table).Where("create_time between $1 and $2", start, end).Where(fmt.Sprintf("work_state != %d", listenerkey.ListenerWorkStateAccountDeleted))
	query, args, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	resp := make([]*ListenerProfile, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerProfileModel) FindRecentActive(ctx context.Context, interval int, limit uint64, state []int64) ([]int64, error) {
	rb := squirrel.Select("listener_uid").From(m.table).Where("update_time > $1", time.Now().AddDate(0, 0, -interval))
	if len(state) == 1 {
		rb = rb.Where(fmt.Sprintf("work_state = %d", state[0]))
	} else if len(state) > 1 {
		rb = rb.Where("work_state = ANY($2)", pq.Int64Array(state))
	} else {
		rb = rb.Where(fmt.Sprintf("work_state != %d", listenerkey.ListenerWorkStateAccountDeleted))
	}
	query, args, err := rb.OrderBy("create_time DESC").Limit(limit).ToSql()
	if err != nil {
		return nil, err
	}
	resp := make([]int64, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerProfileModel) FindAllListener(ctx context.Context, pageNo, pageSize int64) ([]int64, error) {
	rb := squirrel.Select("listener_uid").From(m.table).Where(fmt.Sprintf("work_state != %d", listenerkey.ListenerWorkStateAccountDeleted))

	query, args, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}
	resp := make([]int64, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerProfileModel) CountListenerByWorkState(ctx context.Context, workState int64) (int64, error) {
	rb := squirrel.Select("count(listener_uid)").From(m.table).Where("work_state=$1", workState)
	query, args, err := rb.ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}

func (m *defaultListenerProfileModel) CountListenerByOnlineState(ctx context.Context, onlineState int64) (int64, error) {
	rb := squirrel.Select("count(listener_uid)").From(m.table)
	if onlineState == userkey.Login {
		rb = rb.Where("online_state=$1", onlineState)
	} else {
		rb = rb.Where("online_state!=$1", userkey.Login)
	}
	query, args, err := rb.ToSql()
	var resp int64
	err = m.QueryRowNoCacheCtx(ctx, &resp, query, args...)
	switch err {
	case nil:
		return resp, nil
	default:
		return 0, err
	}
}
