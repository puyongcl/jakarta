package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"jakarta/common/tool"
)

var _ ListenerProfileDraftModel = (*customListenerProfileDraftModel)(nil)

type (
	// ListenerProfileDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customListenerProfileDraftModel.
	ListenerProfileDraftModel interface {
		listenerProfileDraftModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		FindByFilter(ctx context.Context, checkStatus []int64, listenerName string, listenerUid, certType, pageNo, pageSize int64) ([]*ListenerProfileDraft, error)
		FindByFilterCount(ctx context.Context, checkStatus []int64, listenerName string, listenerUid, certType int64) (int64, error)
		InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerProfileDraft) (sql.Result, error)
		UpdateTrans(ctx context.Context, session sqlx.Session, newData *ListenerProfileDraft) error
		UpdateCheckStatus(ctx context.Context, listenerUid int64, checkStatus int64) error
	}

	customListenerProfileDraftModel struct {
		*defaultListenerProfileDraftModel
	}
)

// NewListenerProfileDraftModel returns a model for the database table.
func NewListenerProfileDraftModel(conn sqlx.SqlConn, c cache.CacheConf) ListenerProfileDraftModel {
	return &customListenerProfileDraftModel{
		defaultListenerProfileDraftModel: newListenerProfileDraftModel(conn, c),
	}
}

// export logic
func (m *defaultListenerProfileDraftModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultListenerProfileDraftModel) FindByFilter(ctx context.Context, checkStatus []int64, listenerName string, listenerUid, certType, pageNo, pageSize int64) ([]*ListenerProfileDraft, error) {
	argNo := 1
	rb := squirrel.Select(listenerProfileDraftRows).From(m.table)
	if len(checkStatus) > 0 {
		rb = rb.Where(fmt.Sprintf("check_status = ANY($%d)", argNo), pq.Int64Array(checkStatus))
		argNo++
	}
	if listenerName != "" {
		rb = rb.Where(fmt.Sprintf("nick_name = $%d", argNo), listenerName)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if certType != 0 {
		rb = rb.Where(fmt.Sprintf("cert_type = $%d", argNo), certType)
		argNo++
	}
	// TODO
	//rb = rb.Where(fmt.Sprintf("update_time < $%d", argNo), time.Now().Add(time.Minute*(-5)))

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query, values, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ListenerProfileDraft, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultListenerProfileDraftModel) FindByFilterCount(ctx context.Context, checkStatus []int64, listenerName string, listenerUid int64, certType int64) (int64, error) {
	argNo := 1
	rb := squirrel.Select("COUNT(listener_uid)").From(m.table)
	if len(checkStatus) > 0 {
		rb = rb.Where(fmt.Sprintf("check_status = ANY($%d)", argNo), pq.Int64Array(checkStatus))
		argNo++
	}
	if listenerName != "" {
		rb = rb.Where(fmt.Sprintf("nick_name = $%d", argNo), listenerName)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if certType != 0 {
		rb = rb.Where(fmt.Sprintf("cert_type = $%d", argNo), certType)
		argNo++
	}
	// TODO
	// rb = rb.Where(fmt.Sprintf("update_time < $%d", argNo), time.Now().Add(time.Minute*(-1)))

	query, values, err := rb.ToSql()
	if err != nil {
		return 0, err
	}

	var cnt int64
	err = m.QueryRowNoCacheCtx(ctx, &cnt, query, values...)
	switch err {
	case nil:
		return cnt, nil
	default:
		return 0, err
	}
}

func (m *defaultListenerProfileDraftModel) InsertTrans(ctx context.Context, session sqlx.Session, data *ListenerProfileDraft) (sql.Result, error) {
	jakartaListenerProfileDraftListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileDraftListenerUidPrefix, data.ListenerUid)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42)", m.table, listenerProfileDraftRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.ListenerUid, data.NickName, data.ListenerName, data.Avatar, data.SmallAvatar, data.MaritalStatus, data.PhoneNumber, data.Constellation, data.Province, data.City, data.Job, data.Education, data.Gender, data.Birthday, data.IdNo, data.IdPhoto1, data.IdPhoto2, data.IdPhoto3, data.Specialties, data.Introduction, data.VoiceFile, data.Experience1, data.Experience2, data.CertType, data.OtherPlatformAccount, data.CertFiles1, data.CertFiles2, data.CertFiles3, data.CertFiles4, data.CertFiles5, data.AutoReplyNew, data.AutoReplyProcessing, data.AutoReplyFinish, data.TextChatPrice, data.VoiceChatPrice, data.TextChatSwitch, data.VoiceChatSwitch, data.CheckFailField, data.CheckingField, data.CheckStatus, data.DraftVersion, data.CheckVersion)
		}
		return conn.ExecCtx(ctx, query, data.ListenerUid, data.NickName, data.ListenerName, data.Avatar, data.SmallAvatar, data.MaritalStatus, data.PhoneNumber, data.Constellation, data.Province, data.City, data.Job, data.Education, data.Gender, data.Birthday, data.IdNo, data.IdPhoto1, data.IdPhoto2, data.IdPhoto3, data.Specialties, data.Introduction, data.VoiceFile, data.Experience1, data.Experience2, data.CertType, data.OtherPlatformAccount, data.CertFiles1, data.CertFiles2, data.CertFiles3, data.CertFiles4, data.CertFiles5, data.AutoReplyNew, data.AutoReplyProcessing, data.AutoReplyFinish, data.TextChatPrice, data.VoiceChatPrice, data.TextChatSwitch, data.VoiceChatSwitch, data.CheckFailField, data.CheckingField, data.CheckStatus, data.DraftVersion, data.CheckVersion)
	}, jakartaListenerProfileDraftListenerUidKey)
	return ret, err
}

func (m *defaultListenerProfileDraftModel) UpdateTrans(ctx context.Context, session sqlx.Session, newData *ListenerProfileDraft) error {
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
	if newData.SmallAvatar != data.SmallAvatar {
		rb = rb.Set("small_avatar", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.SmallAvatar))
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
	if !tool.IsEqualArrayString(newData.CheckFailField, data.CheckFailField) {
		rb = rb.Set("check_fail_field", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CheckFailField))
		argNo++
	}
	if !tool.IsEqualArrayString(newData.CheckingField, data.CheckingField) {
		rb = rb.Set("checking_field", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CheckingField))
		argNo++
	}
	if newData.CheckStatus != data.CheckStatus {
		rb = rb.Set("check_status", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CheckStatus))
		argNo++
	}
	if newData.DraftVersion != data.DraftVersion {
		rb = rb.Set("draft_version", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.DraftVersion))
		argNo++
	}
	if newData.CheckVersion != data.CheckVersion {
		rb = rb.Set("check_version", squirrel.Expr(fmt.Sprintf("$%d", argNo), newData.CheckVersion))
		argNo++
	}
	query, args, err := rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), newData.ListenerUid).ToSql()
	if err != nil {
		return err
	}
	jakartaListenerProfileDraftListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileDraftListenerUidPrefix, newData.ListenerUid)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaListenerProfileDraftListenerUidKey)
	return err
}

const updateCheckStatusSql = "update %s set check_status = $1 where listener_uid = $2"

func (m *defaultListenerProfileDraftModel) UpdateCheckStatus(ctx context.Context, listenerUid int64, checkStatus int64) error {
	query := fmt.Sprintf(updateCheckStatusSql, m.table)
	jakartaListenerProfileDraftListenerUidKey := fmt.Sprintf("%s%v", cacheJakartaListenerProfileDraftListenerUidPrefix, listenerUid)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, checkStatus, listenerUid)
	}, jakartaListenerProfileDraftListenerUidKey)
	return err
}
