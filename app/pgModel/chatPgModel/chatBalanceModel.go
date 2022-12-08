package chatPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/chatkey"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"
)

var _ ChatBalanceModel = (*customChatBalanceModel)(nil)

type (
	// ChatBalanceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatBalanceModel.
	ChatBalanceModel interface {
		chatBalanceModel
		Trans(ctx context.Context, fn func(context context.Context, session sqlx.Session) error) error
		UpdateCurrentChatStop(ctx context.Context, session sqlx.Session, uid, listenerUid int64, currentChatLogId string, stopTime *time.Time) (error, string, int64, int64)
		UpdateCurrentChatStart(ctx context.Context, session sqlx.Session, uid, listenerUid int64, currentChatLogId string, startTime *time.Time) (error, string, int64, int64)
		UpdateEnterChatTime(ctx context.Context, id string) error
		AddVoiceChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, addMin int64) (error, string, int64, string)
		ReduceVoiceChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, expiryMin int64) (error, string, int64, string)
		StopTextChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, expiryMin int64) (error, string, int64, string)
		UpdateTextChatTimeOver(ctx context.Context, uid, listenerUid int64) (error, string, int64)
		UpdateUsedVoiceChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, addMin int64) (error, string, int64, string)
		AddTextChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, addMin int64, expiryTime *time.Time) (error, string, int64, string)
		FindVoiceChat(ctx context.Context, uid int64, listenerUid int64, state int64, expiryTimeStart *time.Time, expiryTimeEnd *time.Time, pageNo, pageSize int64) ([]*ChatBalance, error)
		FindTextChat(ctx context.Context, uid int64, listenerUid int64, expiryTimeStart *time.Time, expiryTimeEnd *time.Time, pageNo, pageSize int64) ([]*ChatBalance, error)
		FindByUidAndState(ctx context.Context, uid int64, listenerUid int64, state int64, pageNo, pageSize int64) ([]*ChatBalance, error)
		FindCountEnterChatUserCntRangeUpdateTime(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error)
	}

	customChatBalanceModel struct {
		*defaultChatBalanceModel
	}
)

// NewChatBalanceModel returns a model for the database table.
func NewChatBalanceModel(conn sqlx.SqlConn, c cache.CacheConf) ChatBalanceModel {
	return &customChatBalanceModel{
		defaultChatBalanceModel: newChatBalanceModel(conn, c),
	}
}

func (m *defaultChatBalanceModel) UpdateCurrentChatStart(ctx context.Context, session sqlx.Session, uid, listenerUid int64, currentChatLogId string, startTime *time.Time) (error, string, int64, int64) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, 0
	}

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set current_chat_log_id = $1, current_start_time = $2, current_voice_chat_expiry_time = $3, current_voice_chat_state = $4 where id = $5", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, currentChatLogId, startTime, startTime.Add(time.Duration(data.AvailableVoiceChatMinute)*time.Minute), chatkey.VoiceChatStateStart, data.Id)
		}
		return conn.ExecCtx(ctx, query, currentChatLogId, startTime, time.Now().Add(time.Duration(data.AvailableVoiceChatMinute)*time.Minute), chatkey.VoiceChatStateStart, data.Id)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), data.AvailableVoiceChatMinute, data.UsedVoiceChatMinute
}

func (m *defaultChatBalanceModel) UpdateCurrentChatStop(ctx context.Context, session sqlx.Session, uid, listenerUid int64, currentChatLogId string, stopTime *time.Time) (error, string, int64, int64) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, 0
	}

	if currentChatLogId != data.CurrentChatLogId {
		return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("current_chat_log_id not match %d:%d", data.CurrentChatLogId, currentChatLogId)), "", 0, 0
	}
	var availableMin int64
	var usedMin int64
	if data.CurrentStartTime.Valid {
		usedMin = int64(time.Now().Sub(data.CurrentStartTime.Time).Minutes())
		availableMin = data.AvailableVoiceChatMinute - usedMin
		if availableMin < 0 {
			availableMin = 0
		}
	}
	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set current_stop_time = $1, current_voice_chat_state = $2 where id = $3", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, stopTime, chatkey.VoiceChatStateStop, data.Id)
		}
		return conn.ExecCtx(ctx, query, stopTime, chatkey.VoiceChatStateStop, data.Id)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), availableMin, data.UsedVoiceChatMinute + usedMin
}

func (m *defaultChatBalanceModel) AddVoiceChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, addMin int64) (error, string, int64, string) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, ""
	}
	var query string
	args := make([]interface{}, 0)
	if data.CurrentVoiceChatState == chatkey.VoiceChatStateStart {
		query = fmt.Sprintf("update %s set available_voice_chat_minute = available_voice_chat_minute + $1, current_voice_chat_expiry_time = $2 where id = $3", m.table)
		var currExpiryTime time.Time
		if data.CurrentVoiceChatExpiryTime.Valid {
			currExpiryTime = data.CurrentVoiceChatExpiryTime.Time.Add(time.Duration(addMin) * time.Minute)
		}
		args = append(args, addMin, currExpiryTime, data.Id)
	} else {
		query = fmt.Sprintf("update %s set available_voice_chat_minute = available_voice_chat_minute + $1 where id = $2", m.table)
		args = append(args, addMin, data.Id)
	}

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), data.AvailableVoiceChatMinute + addMin, data.Id
}

func (m *defaultChatBalanceModel) AddTextChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, addMin int64, expiryTime *time.Time) (error, string, int64, string) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, ""
	}

	var query string
	now := time.Now()
	var args []interface{}
	if data.TextChatExpiryTime.Valid && data.TextChatExpiryTime.Time.After(now) {
		data.TextChatExpiryTime.Time = data.TextChatExpiryTime.Time.Add(time.Duration(addMin) * time.Minute)
		query = fmt.Sprintf("update %s set current_text_chat_state = 1, text_chat_expiry_time = text_chat_expiry_time + interval '%d minute' where id = $1", m.table, addMin)
		args = append(args, data.Id)
	} else {
		data.TextChatExpiryTime.Time = time.Date(expiryTime.Year(), expiryTime.Month(), expiryTime.Day(), expiryTime.Hour(), expiryTime.Minute()+1, 0, 0, time.Local)
		query = fmt.Sprintf("update %s set current_text_chat_state = 1, text_chat_expiry_time = $1 where id = $2", m.table)
		args = append(args, data.TextChatExpiryTime.Time, data.Id)
	}

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), data.AvailableVoiceChatMinute, data.Id
}

func (m *defaultChatBalanceModel) UpdateUsedVoiceChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, usedMin int64) (error, string, int64, string) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, ""
	}
	decr := usedMin
	if data.AvailableVoiceChatMinute <= usedMin {
		decr = data.AvailableVoiceChatMinute
		data.AvailableVoiceChatMinute = 0
	} else {
		data.AvailableVoiceChatMinute = data.AvailableVoiceChatMinute - decr
	}

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set used_voice_chat_minute = used_voice_chat_minute + $1, available_voice_chat_minute = available_voice_chat_minute - $2, current_voice_chat_state = $3 where id = $4", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, usedMin, decr, chatkey.VoiceChatStateSettle, data.Id)
		}
		return conn.ExecCtx(ctx, query, usedMin, decr, chatkey.VoiceChatStateSettle, data.Id)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), data.AvailableVoiceChatMinute, data.Id
}

func (m *defaultChatBalanceModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *defaultChatBalanceModel) FindVoiceChat(ctx context.Context, uid int64, listenerUid int64, state int64, expiryTimeStart *time.Time, expiryTimeEnd *time.Time, pageNo, pageSize int64) ([]*ChatBalance, error) {
	rb := squirrel.Select(chatBalanceRows).From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if state != 0 {
		rb = rb.Where(fmt.Sprintf("current_voice_chat_state = $%d", argNo), state)
		argNo++
	}
	if expiryTimeStart != nil && expiryTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("current_voice_chat_expiry_time between $%d and $%d", argNo, argNo+1), expiryTimeStart, expiryTimeEnd)
		argNo += 2
	}

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	query, values, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatBalance, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatBalanceModel) FindTextChat(ctx context.Context, uid int64, listenerUid int64, expiryTimeStart *time.Time, expiryTimeEnd *time.Time, pageNo, pageSize int64) ([]*ChatBalance, error) {
	rb := squirrel.Select(chatBalanceRows).From(m.table).Where("current_text_chat_state = 1")
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if expiryTimeStart != nil && expiryTimeEnd != nil {
		rb = rb.Where(fmt.Sprintf("text_chat_expiry_time between $%d and $%d", argNo, argNo+1), expiryTimeStart, expiryTimeEnd)
		argNo += 2
	}

	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	query, values, err := rb.OrderBy("create_time ASC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatBalance, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *defaultChatBalanceModel) ReduceVoiceChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, expiryMin int64) (error, string, int64, string) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, ""
	}
	availableMin := data.AvailableVoiceChatMinute - expiryMin
	if availableMin < 0 {
		availableMin = 0
		expiryMin = data.AvailableVoiceChatMinute
	}

	var query string
	args := make([]interface{}, 0)
	if data.CurrentVoiceChatState == chatkey.VoiceChatStateStart {
		return xerr.NewGrpcErrCodeMsg(xerr.OrderErrorNotAllowStopOrder, "当前在通话中"), data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), data.AvailableVoiceChatMinute, data.Id
	} else {
		query = fmt.Sprintf("update %s set available_voice_chat_minute = available_voice_chat_minute - $1 where id = $2", m.table)
		args = append(args, expiryMin, data.Id)
	}

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), availableMin, data.Id
}

func (m *defaultChatBalanceModel) StopTextChatTime(ctx context.Context, session sqlx.Session, uid, listenerUid, expiryMin int64) (error, string, int64, string) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0, ""
	}
	var exTime time.Time
	if data.TextChatExpiryTime.Valid {
		exTime = data.TextChatExpiryTime.Time.Add(-time.Duration(expiryMin) * time.Minute)
	}
	var query string
	args := make([]interface{}, 0)
	interval := int64(exTime.Sub(time.Now()).Minutes())
	if tool.Abs(interval) < 1 { // 已经用完
		query = fmt.Sprintf("update %s set current_text_chat_state = 2, text_chat_expiry_time = text_chat_expiry_time - interval '%d minute' where id = $1", m.table, expiryMin)
		args = append(args, data.Id)
	} else {
		query = fmt.Sprintf("update %s set text_chat_expiry_time = text_chat_expiry_time - interval '%d minute' where id = $1", m.table, expiryMin)
		args = append(args, data.Id)
	}

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatBalanceIdKey)
	return err, exTime.Format(db.DateTimeFormat), data.AvailableVoiceChatMinute, data.Id
}

func (m *defaultChatBalanceModel) FindByUidAndState(ctx context.Context, uid int64, listenerUid int64, state int64, pageNo, pageSize int64) ([]*ChatBalance, error) {
	rb := squirrel.Select(chatBalanceRows).From(m.table)
	argNo := 1
	if uid != 0 {
		rb = rb.Where(fmt.Sprintf("uid = $%d", argNo), uid)
		argNo++
	}
	if listenerUid != 0 {
		rb = rb.Where(fmt.Sprintf("listener_uid = $%d", argNo), listenerUid)
		argNo++
	}
	if state != 0 {
		rb = rb.Where(fmt.Sprintf("current_voice_chat_state = $%d", argNo), state)
		argNo++
	}
	// 分页
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	query, values, err := rb.OrderBy("create_time DESC").Limit(uint64(pageSize)).Offset(uint64((pageNo - 1) * pageSize)).ToSql()
	if err != nil {
		return nil, err
	}

	resp := make([]*ChatBalance, 0)
	err = m.QueryRowsNoCacheCtx(ctx, &resp, query, values...)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

const queryFindEnterChatUserCntByCreateTimeSQL = "select count(id) from %s where listener_uid = $1 and update_time between $2 and $3"

func (m *defaultChatBalanceModel) FindCountEnterChatUserCntRangeUpdateTime(ctx context.Context, listenerUid int64, start, end *time.Time) (int64, error) {
	query := fmt.Sprintf(queryFindEnterChatUserCntByCreateTimeSQL, m.table)
	var cnt int64
	err := m.QueryRowNoCacheCtx(ctx, &cnt, query, listenerUid, start, end)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (m *defaultChatBalanceModel) UpdateTextChatTimeOver(ctx context.Context, uid, listenerUid int64) (error, string, int64) {
	data, err := m.FindOne(ctx, fmt.Sprintf(db.DBUidId, uid, listenerUid))
	if err != nil {
		return err, "", 0
	}
	if data.CurrentTextChatState != chatkey.TextChatStateStart {
		return nil, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), chatkey.TextChatStateAlreadyStop
	}
	if data.TextChatExpiryTime.Valid && data.TextChatExpiryTime.Time.After(time.Now()) { // 还未结束
		return nil, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), data.CurrentTextChatState
	}
	var query string
	args := make([]interface{}, 0)

	query = fmt.Sprintf("update %s set current_text_chat_state = 2 where id = $1", m.table)
	args = append(args, data.Id)

	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, jakartaChatBalanceIdKey)
	return err, data.TextChatExpiryTime.Time.Format(db.DateTimeFormat), chatkey.TextChatStateStop
}

func (m *defaultChatBalanceModel) UpdateEnterChatTime(ctx context.Context, id string) error {
	jakartaChatBalanceIdKey := fmt.Sprintf("%s%v", cacheJakartaChatBalanceIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set today_first_enter_time = $1 where id = $2", m.table)
		return conn.ExecCtx(ctx, query, time.Now(), id)
	}, jakartaChatBalanceIdKey)
	return err
}
