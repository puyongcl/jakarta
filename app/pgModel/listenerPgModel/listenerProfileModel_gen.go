// Code generated by goctl. DO NOT EDIT!

package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	listenerProfileFieldNames          = builder.RawFieldNames(&ListenerProfile{}, true)
	listenerProfileRows                = strings.Join(listenerProfileFieldNames, ",")
	listenerProfileRowsExpectAutoSet   = strings.Join(stringx.Remove(listenerProfileFieldNames, "create_time", "update_time", "create_at", "update_at"), ",")
	listenerProfileRowsWithPlaceHolder = builder.PostgreSqlJoin(stringx.Remove(listenerProfileFieldNames, "listener_uid", "create_time", "update_time", "create_at", "update_at"))

	cacheJakartaListenerProfileListenerUidPrefix = "cache:jakarta:listenerProfile:listenerUid:"
)

type (
	listenerProfileModel interface {
		Insert(ctx context.Context, data *ListenerProfile) (sql.Result, error)
		FindOne(ctx context.Context, listenerUid int64) (*ListenerProfile, error)
		Update(ctx context.Context, newData *ListenerProfile) error
		Delete(ctx context.Context, listenerUid int64) error
	}

	defaultListenerProfileModel struct {
		sqlc.CachedConn
		table string
	}

	ListenerProfile struct {
		CreateTime           time.Time     `db:"create_time"`
		UpdateTime           time.Time     `db:"update_time"`
		ListenerUid          int64         `db:"listener_uid"`
		NickName             string        `db:"nick_name"`
		ListenerName         string        `db:"listener_name"`
		OpenId               string        `db:"open_id"`
		Avatar               string        `db:"avatar"`
		MaritalStatus        int64         `db:"marital_status"`
		PhoneNumber          string        `db:"phone_number"`
		Constellation        int64         `db:"constellation"`
		Province             string        `db:"province"`
		City                 string        `db:"city"`
		Job                  string        `db:"job"`
		Education            int64         `db:"education"`
		Gender               int64         `db:"gender"`
		Birthday             time.Time     `db:"birthday"`
		IdNo                 string        `db:"id_no"`
		IdPhoto1             string        `db:"id_photo1"`
		IdPhoto2             string        `db:"id_photo2"`
		IdPhoto3             string        `db:"id_photo3"`
		Specialties          pq.Int64Array `db:"specialties"`
		Introduction         string        `db:"introduction"`
		VoiceFile            string        `db:"voice_file"`
		Experience1          string        `db:"experience1"`
		Experience2          string        `db:"experience2"`
		CertType             int64         `db:"cert_type"`
		OtherPlatformAccount string        `db:"other_platform_account"`
		CertFiles1           string        `db:"cert_files1"`
		CertFiles2           string        `db:"cert_files2"`
		CertFiles3           string        `db:"cert_files3"`
		CertFiles4           string        `db:"cert_files4"`
		CertFiles5           string        `db:"cert_files5"`
		AutoReplyNew         string        `db:"auto_reply_new"`
		AutoReplyProcessing  string        `db:"auto_reply_processing"`
		AutoReplyFinish      string        `db:"auto_reply_finish"`
		TextChatPrice        int64         `db:"text_chat_price"`
		VoiceChatPrice       int64         `db:"voice_chat_price"`
		TextChatSwitch       int64         `db:"text_chat_switch"`
		VoiceChatSwitch      int64         `db:"voice_chat_switch"`
		UserCount            int64         `db:"user_count"`
		ChatDuration         int64         `db:"chat_duration"`
		RatingSum            int64         `db:"rating_sum"`
		FiveStar             int64         `db:"five_star"`
		ThreeStar            int64         `db:"three_star"`
		OneStar              int64         `db:"one_star"`
		RestingTimeEnable    int64         `db:"resting_time_enable"`
		StopWorkTime         string        `db:"stop_work_time"`
		StartWorkTime        string        `db:"start_work_time"`
		WorkDays             pq.Int64Array `db:"work_days"`
		WorkState            int64         `db:"work_state"`
		RefundOrderCnt       int64         `db:"refund_order_cnt"`     // 退款的订单数
		FinishOrderCnt       int64         `db:"finish_order_cnt"`     // 确认完成的订单数
		PaidOrderCnt         int64         `db:"paid_order_cnt"`       // 支付成功的订单数
		RepeatPaidUserCnt    int64         `db:"repeat_paid_user_cnt"` // 复购用户数
		OnlineState          int64         `db:"online_state"`         // 在线状态
		Channel              string        `db:"channel"`              // 获客渠道
	}
)

func newListenerProfileModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultListenerProfileModel {
	return &defaultListenerProfileModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      `"jakarta"."listener_profile"`,
	}
}

func (m *defaultListenerProfileModel) Delete(ctx context.Context, listenerUid int64) error {
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where listener_uid = $1", m.table)
		return conn.ExecCtx(ctx, query, listenerUid)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) FindOne(ctx context.Context, listenerUid int64) (*ListenerProfile, error) {
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, listenerUid)
	var resp ListenerProfile
	err := m.QueryRowCtx(ctx, &resp, jakartaListenerProfileListenerUidKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where listener_uid = $1 limit 1", listenerProfileRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, listenerUid)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultListenerProfileModel) Insert(ctx context.Context, data *ListenerProfile) (sql.Result, error) {
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54)", m.table, listenerProfileRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.ListenerUid, data.NickName, data.ListenerName, data.OpenId, data.Avatar, data.MaritalStatus, data.PhoneNumber, data.Constellation, data.Province, data.City, data.Job, data.Education, data.Gender, data.Birthday, data.IdNo, data.IdPhoto1, data.IdPhoto2, data.IdPhoto3, data.Specialties, data.Introduction, data.VoiceFile, data.Experience1, data.Experience2, data.CertType, data.OtherPlatformAccount, data.CertFiles1, data.CertFiles2, data.CertFiles3, data.CertFiles4, data.CertFiles5, data.AutoReplyNew, data.AutoReplyProcessing, data.AutoReplyFinish, data.TextChatPrice, data.VoiceChatPrice, data.TextChatSwitch, data.VoiceChatSwitch, data.UserCount, data.ChatDuration, data.RatingSum, data.FiveStar, data.ThreeStar, data.OneStar, data.RestingTimeEnable, data.StopWorkTime, data.StartWorkTime, data.WorkDays, data.WorkState, data.RefundOrderCnt, data.FinishOrderCnt, data.PaidOrderCnt, data.RepeatPaidUserCnt, data.OnlineState, data.Channel)
	}, jakartaListenerProfileListenerUidKey)
	return ret, err
}

func (m *defaultListenerProfileModel) Update(ctx context.Context, data *ListenerProfile) error {
	jakartaListenerProfileListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, data.ListenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where listener_uid = $1", m.table, listenerProfileRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, data.ListenerUid, data.NickName, data.ListenerName, data.OpenId, data.Avatar, data.MaritalStatus, data.PhoneNumber, data.Constellation, data.Province, data.City, data.Job, data.Education, data.Gender, data.Birthday, data.IdNo, data.IdPhoto1, data.IdPhoto2, data.IdPhoto3, data.Specialties, data.Introduction, data.VoiceFile, data.Experience1, data.Experience2, data.CertType, data.OtherPlatformAccount, data.CertFiles1, data.CertFiles2, data.CertFiles3, data.CertFiles4, data.CertFiles5, data.AutoReplyNew, data.AutoReplyProcessing, data.AutoReplyFinish, data.TextChatPrice, data.VoiceChatPrice, data.TextChatSwitch, data.VoiceChatSwitch, data.UserCount, data.ChatDuration, data.RatingSum, data.FiveStar, data.ThreeStar, data.OneStar, data.Specialties, data.StopWorkTime, data.StartWorkTime, data.WorkDays, data.WorkState, data.RefundOrderCnt, data.FinishOrderCnt, data.PaidOrderCnt, data.RepeatPaidUserCnt, data.OnlineState, data.Channel)
	}, jakartaListenerProfileListenerUidKey)
	return err
}

func (m *defaultListenerProfileModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheJakartaListenerProfileListenerUidPrefix, primary)
}

func (m *defaultListenerProfileModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where listener_uid = $1 limit 1", listenerProfileRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultListenerProfileModel) tableName() string {
	return m.table
}
