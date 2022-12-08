package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/kqueue"
	"jakarta/common/notify"
	"jakarta/common/third_party/tim"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"strings"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCheckListenerProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminCheckListenerProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCheckListenerProfileLogic {
	return &AdminCheckListenerProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  管理员接口 审核XXX
func (l *AdminCheckListenerProfileLogic) AdminCheckListenerProfile(in *pb.CheckListenerProfileReq) (*pb.CheckListenerProfileResp, error) {
	// 加分布式锁
	rkey := fmt.Sprintf(rediskey.RedisLockEditListenerProfile, in.ListenerUid)
	rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
	rl.SetExpire(2)
	b, err := rl.AcquireCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
	}
	defer func() {
		_, err2 := rl.ReleaseCtx(l.ctx)
		if err2 != nil {
			logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
			return
		}
	}()

	var rsp *pb.CheckListenerProfileResp
	rsp, err = l.checkListenerProfile(in)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func (l *AdminCheckListenerProfileLogic) checkListenerProfile(in *pb.CheckListenerProfileReq) (rsp *pb.CheckListenerProfileResp, err error) {
	var data *listenerPgModel2.ListenerProfile
	data, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel2.ErrNotFound {
		return
	}

	var draft *listenerPgModel2.ListenerProfileDraft
	draft, err = l.svcCtx.ListenerProfileDraftModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil {
		return
	}

	// 检查状态
	if in.CheckStatus != draft.CheckStatus {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "审核状态不匹配")
	}
	if in.DraftVersion <= draft.CheckVersion {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "审核版本已经改变")
	}

	// 检查出审核的字段是否与当前数据库中的值相同
	var changeField []string
	changeField = l.veryDraft(in, draft)
	if len(changeField) > 0 {
		in.CheckFailField = stringx.Remove(in.CheckFailField, changeField...)
		in.CheckPassField = stringx.Remove(in.CheckPassField, changeField...)
	}

	// 审核失败的字段
	draft.CheckFailField = stringx.Remove(draft.CheckFailField, in.CheckPassField...)
	draft.CheckFailField = tool.CombineStringArray(draft.CheckFailField, in.CheckFailField)
	// 审核中的字段
	draft.CheckingField = stringx.Remove(draft.CheckingField, in.CheckPassField...)
	draft.CheckingField = stringx.Remove(draft.CheckingField, in.CheckFailField...)

	// 所有字段是否全部审核
	if len(draft.CheckingField) > 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("请审核完全部字段 %s", strings.Join(draft.CheckingField, ",")))
	}

	// 判断当前审核状态
	switch draft.CheckStatus {
	case listenerkey.CheckStatusFirstApplyEdit: // 首次申请 未提交
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "未提交审核")
	case listenerkey.CheckStatusFirstApplyChecking: // 首次申请 已经提交 未审核
		if len(in.CheckFailField) == 0 && len(changeField) == 0 && len(draft.CheckingField) == 0 { // 申请全部通过
			draft.CheckStatus = listenerkey.CheckStatusFirstApplyPass
		} else if len(in.CheckFailField) > 0 && len(changeField) == 0 { // 审核失败 并且当前未修改
			draft.CheckStatus = listenerkey.CheckStatusFirstApplyRefuse
		}

	case listenerkey.CheckStatusFirstApplyRefuse: // 未修改 已经审核
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "当前审核状态不通过，用户尚未作出修改")

	case listenerkey.CheckStatusFirstApplyPass: // 已经通过
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "当前审核状态未审核通过，不能审核")

	case listenerkey.CheckStatusEditWaitChecking: // XXX修改待审核
		if len(in.CheckFailField) == 0 && len(changeField) == 0 && len(draft.CheckingField) == 0 { // 申请全部通过
			draft.CheckStatus = listenerkey.CheckStatusEditPass
		} else if len(in.CheckFailField) > 0 && len(changeField) == 0 { // 审核失败 并且当前未修改
			draft.CheckStatus = listenerkey.CheckStatusEditRefuse
		}

	case listenerkey.CheckStatusEditRefuse: // 已经拒绝 未修改
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "当前审核状态不通过，用户尚未作出修改")

	case listenerkey.CheckStatusEditPass: // 已经审核通过
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "当前审核状态未审核通过，不能审核")
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("当前审核状态错误 %d，不能审核", draft.CheckStatus))
	}

	// 更新
	err = l.svcCtx.ListenerProfileDraftModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var err2 error
		// 更新profile
		if data == nil && len(in.CheckFailField) == 0 { // 首次申请 全部审核通过 才会更新到资料
			err2 = l.firstApply(ctx, session, in)
			if err2 != nil {
				return err2
			}
		} else if data != nil && len(in.CheckPassField) > 0 && len(in.CheckFailField) == 0 { // XXX修改 全部审核通过才可以更新
			imnick, imavatar := l.getUpdateProfile(in, data)
			err2 = l.svcCtx.ListenerProfileModel.UpdateTrans(ctx, session, data)
			if err2 != nil {
				return err2
			}
			// 更新IM
			err2 = l.svcCtx.TimClient.UpdateProfile(in.ListenerUid, strings.Replace(imnick, "#", "", -1), imavatar)
			if err2 != nil {
				return err2
			}
		}

		// 更新draft
		draft.CheckVersion = in.DraftVersion

		err2 = l.svcCtx.ListenerProfileDraftModel.UpdateTrans(ctx, session, draft)
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	rsp = &pb.CheckListenerProfileResp{CheckResult: draft.CheckStatus}
	// 消息通知
	l.notify(draft, in.Remark, in.CheckPassField)
	return
}

func (l *AdminCheckListenerProfileLogic) notify(draft *listenerPgModel2.ListenerProfileDraft, remark string, checkPassField []string) {
	msgs := make([]*kqueue.SendImDefineMessage, 0)
	switch draft.CheckStatus {
	case listenerkey.CheckStatusFirstApplyRefuse:
		msg := &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   draft.ListenerUid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg19,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle19,
			Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateSystemMsg19, remark, listenerkey.GetCheckFieldText(draft.CheckFailField)),
			Val1:    "",
			Val2:    "",
			Val3:    draft.ListenerName,
			Val4:    draft.PhoneNumber,
			Val5:    time.Now().Format(db.DateTimeFormat),
			Sync:    tim.TimMsgSyncFromNo,
		}
		msgs = append(msgs, msg)
	case listenerkey.CheckStatusFirstApplyPass:
		msg := &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   draft.ListenerUid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg20,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle20,
			Text:    notify.DefineNotifyMsgTemplateSystemMsg20,
			Val1:    "",
			Val2:    "",
			Val3:    draft.ListenerName,
			Val4:    draft.PhoneNumber,
			Val5:    time.Now().Format(db.DateTimeFormat),
			Sync:    tim.TimMsgSyncFromNo,
		}
		msg2 := &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   draft.ListenerUid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg26,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle26,
			Text:    notify.DefineNotifyMsgTemplateSystemMsg26,
			Val1:    "",
			Val2:    "",
			Sync:    tim.TimMsgSyncFromNo,
		}
		msgs = append(msgs, msg, msg2)
	case listenerkey.CheckStatusEditPass:
		msg := &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   draft.ListenerUid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg21,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle21,
			Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateSystemMsg21, listenerkey.GetCheckFieldText(checkPassField)),
			Val1:    "",
			Val2:    "",
			Val3:    draft.ListenerName,
			Val4:    draft.PhoneNumber,
			Val5:    time.Now().Format(db.DateTimeFormat),
			Sync:    tim.TimMsgSyncFromNo,
		}
		msgs = append(msgs, msg)
	case listenerkey.CheckStatusEditRefuse:
		msg := &kqueue.SendImDefineMessage{
			FromUid: notify.TimSystemNotifyUid,
			ToUid:   draft.ListenerUid,
			MsgType: notify.DefineNotifyMsgTypeSystemMsg22,
			Title:   notify.DefineNotifyMsgTemplateSystemMsgTitle22,
			Text:    fmt.Sprintf(notify.DefineNotifyMsgTemplateSystemMsg22, fmt.Sprintf(",%s", remark), listenerkey.GetCheckFieldText(draft.CheckFailField)),
			Val1:    "",
			Val2:    "",
			Val3:    draft.ListenerName,
			Val4:    draft.PhoneNumber,
			Val5:    time.Now().Format(db.DateTimeFormat),
			Sync:    tim.TimMsgSyncFromNo,
		}
		msgs = append(msgs, msg)
	default:
		return
	}
	for idx := 0; idx < len(msgs); idx++ {
		l.pushNotify(msgs[idx])
	}
}
func (l *AdminCheckListenerProfileLogic) pushNotify(msg *kqueue.SendImDefineMessage) {
	buf, err := json.Marshal(msg)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AdminCheckListenerProfileLogic notify json marshal err:%+v", err)
		return
	}
	err = l.svcCtx.KqueueSendDefineMsgClient.Push(string(buf))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("AdminCheckListenerProfileLogic notify kafka Push err:%+v", err)
		return
	}
	return
}

func (l *AdminCheckListenerProfileLogic) firstApply(ctx context.Context, session sqlx.Session, in *pb.CheckListenerProfileReq) error {
	var err error
	data := &listenerPgModel2.ListenerProfile{}
	_ = copier.Copy(data, in)
	data.WorkState = listenerkey.ListenerWorkStateRestingManual
	data.StartWorkTime = listenerkey.ListenerRestingStartTime
	data.StopWorkTime = listenerkey.ListenerRestingStopTime
	data.WorkDays = []int64{1, 2, 3, 4, 5, 6, 7}
	data.RestingTimeEnable = listenerkey.ListenerRestingSwitchEnable
	data.Birthday, err = time.Parse(db.DateFormat, in.Birthday)
	if err != nil {
		logx.WithContext(ctx).Errorf("checkListenerProfile %s date format err:%+v", in.Birthday, err)
	}
	data.Constellation = in.Constellation
	// XXX资料
	_, err = l.svcCtx.ListenerProfileModel.InsertTrans(ctx, session, data)
	if err != nil {
		return err
	}
	// XXX钱包
	_, err = l.svcCtx.ListenerWalletModel.InsertTrans(ctx, session, &listenerPgModel2.ListenerWallet{
		ListenerUid: in.ListenerUid,
	})
	if err != nil {
		return err
	}
	// XXX常用语
	wordsData := new(listenerPgModel2.ListenerWords)
	*wordsData = listenerPgModel2.ListenerWords{
		ListenerUid: in.ListenerUid,
		Words1:      listenerkey.DefaultWords1,
		Words2:      listenerkey.DefaultWords2,
		Words3:      listenerkey.DefaultWords3,
		Words4:      listenerkey.DefaultWords4,
		WordsSort:   pq.Int64Array{1, 2, 3, 4},
	}
	_, err = l.svcCtx.ListenerWordsModel.InsertTrans(ctx, session, wordsData)
	if err != nil {
		return err
	}
	// 初始化首页统计数据看板
	err = l.svcCtx.StatRedis.InitListenerDashboard(ctx, in.ListenerUid)
	if err != nil {
		return err
	}
	// 初始化统计数据
	stat := new(listenerPgModel2.ListenerDashboardStat)
	stat.ListenerUid = in.ListenerUid
	stat.Suggestion = []int64{listenerkey.SuggestionNo1}
	_, err = l.svcCtx.ListenerDashboardStatModel.InsertTrans(ctx, session, stat)
	if err != nil {
		return err
	}
	// 初始化今日排名
	ulo := NewUpdateListenerOrderStatLogic(l.ctx, l.svcCtx)
	_, err = ulo.UpdateListenerOrderStat(&pb.UpdateListenerOrderStatReq{ListenerUid: in.ListenerUid})
	if err != nil {
		return err
	}
	// 更新IM资料
	err = l.svcCtx.TimClient.UpdateProfile(in.ListenerUid, strings.Replace(in.NickName, "#", "", -1), in.Avatar)
	if err != nil {
		return err
	}

	// 生成提现合同
	glc := NewGenListenerContractLogic(l.ctx, l.svcCtx)
	_, err = glc.DoGenListenerContract(&pb.GenListenerContractReq{
		ListenerUid:  in.ListenerUid,
		ListenerName: in.ListenerName,
		IdNo:         in.IdNo,
		PhoneNumber:  in.PhoneNumber,
		CheckStatus:  in.CheckStatus,
		ContractType: listenerkey.ListenerContractHFBFSZJJHZHB,
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *AdminCheckListenerProfileLogic) veryDraft(in *pb.CheckListenerProfileReq, draftData *listenerPgModel2.ListenerProfileDraft) (changeField []string) {
	if draftData.NickName != in.NickName {
		changeField = append(changeField, "nickName")
	}

	if draftData.Avatar != in.Avatar {
		changeField = append(changeField, "avatar")
	}

	if draftData.Province != in.Province || draftData.City != in.City {
		changeField = append(changeField, "region")
	}

	if draftData.Gender != in.Gender {
		changeField = append(changeField, "gender")
	}
	var draftDataBirthday string
	if draftData.Birthday.Valid {
		draftDataBirthday = draftData.Birthday.Time.Format(db.DateFormat)
	}
	if draftDataBirthday != in.Birthday || draftData.Constellation != in.Constellation {
		changeField = append(changeField, "birthday")
	}

	if draftData.IdNo != in.IdNo || draftData.IdPhoto1 != in.IdPhoto1 || draftData.IdPhoto2 != in.IdPhoto2 || draftData.IdPhoto3 != in.IdPhoto3 || draftData.ListenerName != in.ListenerName {
		changeField = append(changeField, "id")
		if draftData.IdNo != in.IdNo {
			changeField = append(changeField, "idNo")
		}
		if draftData.IdPhoto1 != in.IdPhoto1 {
			changeField = append(changeField, "idPhoto1")
		}
		if draftData.IdPhoto2 != in.IdPhoto2 {
			changeField = append(changeField, "idPhoto2")
		}
		if draftData.IdPhoto3 != in.IdPhoto3 {
			changeField = append(changeField, "idPhoto3")
		}
		if draftData.ListenerName != in.ListenerName {
			changeField = append(changeField, "listenerName")
		}
	}

	if !tool.IsEqualArrayInt64(draftData.Specialties, in.Specialties) {
		changeField = append(changeField, "specialties")
	}

	if draftData.Introduction != in.Introduction {
		changeField = append(changeField, "introduction")
	}
	if draftData.VoiceFile != in.VoiceFile {
		changeField = append(changeField, "voiceFile")
	}
	if draftData.Experience1 != in.Experience1 {
		changeField = append(changeField, "experience2")
	}
	if draftData.Experience2 != in.Experience2 {
		changeField = append(changeField, "experience2")
	}
	if draftData.CertType != in.CertType || draftData.CertFiles1 != in.CertFiles1 || draftData.CertFiles2 != in.CertFiles2 ||
		draftData.CertFiles3 != in.CertFiles3 || draftData.CertFiles4 != in.CertFiles4 || draftData.CertFiles5 != in.CertFiles5 || draftData.OtherPlatformAccount != in.OtherPlatformAccount {
		changeField = append(changeField, "cert")
		if draftData.CertType != in.CertType {
			changeField = append(changeField, "certType")
		}
		if draftData.CertFiles1 != in.CertFiles1 {
			changeField = append(changeField, "certFiles1")
		}
		if draftData.CertFiles2 != in.CertFiles2 {
			changeField = append(changeField, "certFiles2")
		}
		if draftData.CertFiles3 != in.CertFiles3 {
			changeField = append(changeField, "certFiles3")
		}
		if draftData.CertFiles4 != in.CertFiles4 {
			changeField = append(changeField, "certFiles4")
		}
		if draftData.CertFiles5 != in.CertFiles5 {
			changeField = append(changeField, "certFiles5")
		}
		if draftData.OtherPlatformAccount != in.OtherPlatformAccount {
			changeField = append(changeField, "otherPlatformAccount")
		}
	}
	if draftData.AutoReplyNew != in.AutoReplyNew || draftData.AutoReplyProcessing != in.AutoReplyProcessing || draftData.AutoReplyFinish != in.AutoReplyFinish {
		changeField = append(changeField, "autoReply")
		if draftData.AutoReplyNew != in.AutoReplyNew {
			changeField = append(changeField, "autoReplyNew")
		}
		if draftData.AutoReplyProcessing != in.AutoReplyProcessing {
			changeField = append(changeField, "autoReplyProcessing")
		}
		if draftData.AutoReplyFinish != in.AutoReplyFinish {
			changeField = append(changeField, "autoReplyFinish")
		}
	}
	if draftData.TextChatPrice != in.TextChatPrice {
		changeField = append(changeField, "textChatPrice")
	}
	if draftData.VoiceChatPrice != in.VoiceChatPrice {
		changeField = append(changeField, "voiceChatPrice")
	}
	return
}

func (l *AdminCheckListenerProfileLogic) getUpdateProfile(in *pb.CheckListenerProfileReq, data *listenerPgModel2.ListenerProfile) (nick, avatar string) {
	if tool.IsStringArrayExist("nickName", in.CheckPassField) {
		if in.NickName != "" {
			data.NickName = in.NickName
			nick = in.NickName
		}
	}
	if tool.IsStringArrayExist("avatar", in.CheckPassField) {
		if in.Avatar != "" {
			data.Avatar = in.Avatar
			avatar = in.Avatar
		}
	}
	if tool.IsStringArrayExist("province", in.CheckPassField) {
		if in.Province != "" {
			data.Province = in.Province
		}
	}
	if tool.IsStringArrayExist("city", in.CheckPassField) {
		if in.City != "" {
			data.City = in.City
		}
	}
	if tool.IsStringArrayExist("gender", in.CheckPassField) {
		data.Gender = in.Gender
	}
	if tool.IsStringArrayExist("birthday", in.CheckPassField) {
		var err error
		data.Birthday, err = time.Parse(db.DateFormat, in.Birthday)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("getUpdateProfile %s date format err:%+v", in.Birthday, err)
		}
		data.Constellation = in.Constellation
	}
	if tool.IsStringArrayExist("listenerName", in.CheckPassField) {
		if in.ListenerName != "" {
			data.ListenerName = in.ListenerName
		}
	}
	if tool.IsStringArrayExist("idNo", in.CheckPassField) {
		if in.IdNo != "" {
			data.IdNo = in.IdNo
		}
	}
	if tool.IsStringArrayExist("idPhoto1", in.CheckPassField) {
		if in.IdPhoto1 != "" {
			data.IdPhoto1 = in.IdPhoto1
		}
	}
	if tool.IsStringArrayExist("idPhoto2", in.CheckPassField) {
		if in.IdPhoto2 != "" {
			data.IdPhoto2 = in.IdPhoto2
		}
	}
	if tool.IsStringArrayExist("idPhoto3", in.CheckPassField) {
		if in.IdPhoto3 != "" {
			data.IdPhoto3 = in.IdPhoto3
		}
	}
	if tool.IsStringArrayExist("specialties", in.CheckPassField) {
		if len(in.Specialties) > 0 {
			data.Specialties = in.Specialties
		}
	}
	if tool.IsStringArrayExist("introduction", in.CheckPassField) {
		if in.Introduction != "" {
			data.Introduction = in.Introduction
		}
	}
	if tool.IsStringArrayExist("voiceFile", in.CheckPassField) {
		if in.VoiceFile == listenerkey.EmptyString {
			data.VoiceFile = ""
		} else if in.VoiceFile != "" {
			data.VoiceFile = in.VoiceFile
		}
	}
	if tool.IsStringArrayExist("experience1", in.CheckPassField) {
		if in.Experience1 == listenerkey.EmptyString {
			data.Experience1 = ""
		} else if in.Experience1 != "" {
			data.Experience1 = in.Experience1
		}
	}
	if tool.IsStringArrayExist("experience2", in.CheckPassField) {
		if in.Experience2 == listenerkey.EmptyString {
			data.Experience2 = ""
		} else if in.Experience2 != "" {
			data.Experience2 = in.Experience2
		}
	}
	if tool.IsStringArrayExist("certType", in.CheckPassField) {
		if in.CertType != 0 {
			data.CertType = in.CertType
		}
	}
	if tool.IsStringArrayExist("otherPlatformAccount", in.CheckPassField) {
		if in.OtherPlatformAccount != "" {
			if in.OtherPlatformAccount == listenerkey.EmptyString {
				data.OtherPlatformAccount = ""
			} else {
				data.OtherPlatformAccount = in.OtherPlatformAccount
			}
		}
	}
	if tool.IsStringArrayExist("certFiles1", in.CheckPassField) {
		if in.CertFiles1 != "" {
			if in.CertFiles1 == listenerkey.EmptyString {
				data.CertFiles1 = ""
			} else {
				data.CertFiles1 = in.CertFiles1
			}
		}
	}
	if tool.IsStringArrayExist("certFiles2", in.CheckPassField) {
		if in.CertFiles2 != "" {
			if in.CertFiles2 == listenerkey.EmptyString {
				data.CertFiles2 = ""
			} else {
				data.CertFiles2 = in.CertFiles2
			}
		}
	}
	if tool.IsStringArrayExist("certFiles3", in.CheckPassField) {
		if in.CertFiles3 != "" {
			if in.CertFiles3 == listenerkey.EmptyString {
				data.CertFiles3 = ""
			} else {
				data.CertFiles3 = in.CertFiles3
			}
		}
	}
	if tool.IsStringArrayExist("certFiles4", in.CheckPassField) {
		if in.CertFiles4 != "" {
			if in.CertFiles4 == listenerkey.EmptyString {
				data.CertFiles4 = ""
			} else {
				data.CertFiles4 = in.CertFiles4
			}
		}
	}
	if tool.IsStringArrayExist("certFiles5", in.CheckPassField) {
		if in.CertFiles5 != "" {
			if in.CertFiles5 == listenerkey.EmptyString {
				data.CertFiles5 = ""
			} else {
				data.CertFiles5 = in.CertFiles5
			}
		}
	}
	if tool.IsStringArrayExist("autoReplyNew", in.CheckPassField) {
		if in.AutoReplyNew == listenerkey.EmptyString {
			data.AutoReplyNew = ""
		} else if in.AutoReplyNew != "" {
			data.AutoReplyNew = in.AutoReplyNew
		}
	}
	if tool.IsStringArrayExist("autoReplyProcessing", in.CheckPassField) {
		if in.AutoReplyProcessing == listenerkey.EmptyString {
			data.AutoReplyProcessing = ""
		} else if in.AutoReplyProcessing != "" {
			data.AutoReplyProcessing = in.AutoReplyProcessing
		}
	}
	if tool.IsStringArrayExist("autoReplyFinish", in.CheckPassField) {
		if in.AutoReplyFinish == listenerkey.EmptyString {
			data.AutoReplyFinish = ""
		} else if in.AutoReplyFinish != "" {
			data.AutoReplyFinish = in.AutoReplyFinish
		}
	}
	if tool.IsStringArrayExist("textChatPrice", in.CheckPassField) {
		if in.TextChatPrice == listenerkey.EmptyInt {
			data.TextChatPrice = 0
		} else if in.TextChatPrice != 0 {
			data.TextChatPrice = in.TextChatPrice
		}
	}
	if tool.IsStringArrayExist("voiceChatPrice", in.CheckPassField) {
		if in.VoiceChatPrice == listenerkey.EmptyInt {
			data.VoiceChatPrice = 0
		} else if in.VoiceChatPrice != 0 {
			data.VoiceChatPrice = in.VoiceChatPrice
		}
	}
	return
}
