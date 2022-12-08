package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/rediskey"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddOrUpdateListenerProfileDraftLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddOrUpdateListenerProfileDraftLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOrUpdateListenerProfileDraftLogic {
	return &AddOrUpdateListenerProfileDraftLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新申请或者编辑
func (l *AddOrUpdateListenerProfileDraftLogic) AddOrUpdateListenerProfileDraft(in *pb.EditListenerProfileDraftReq) (*pb.EditListenerProfileDraftResp, error) {
	if in.ListenerUid == 0 {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "no listener uid")
	}
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

	//
	var draft *listenerPgModel2.ListenerProfileDraft
	draft, err = l.svcCtx.ListenerProfileDraftModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel2.ErrNotFound {
		return nil, err
	}
	if draft == nil { // 新增
		return l.AddListenerProfileDraft(in)
	}

	return l.UpdateListenerProfileDraft(draft, in)
}

func (l *AddOrUpdateListenerProfileDraftLogic) AddListenerProfileDraft(in *pb.EditListenerProfileDraftReq) (*pb.EditListenerProfileDraftResp, error) {
	draftData := new(listenerPgModel2.ListenerProfileDraft)
	draftData.Specialties = make([]int64, 0)
	draftData.CheckFailField = make([]string, 0)
	draftData.CheckingField = make([]string, 0)

	_ = copier.Copy(&draftData, in)

	// 显示初始化为空 否则会导致数据库中写入一条 0001-01-01 00:00:00
	draftData.Birthday = sql.NullTime{Time: time.Time{}, Valid: false}

	var err error
	if in.Birthday != "" {
		draftData.Birthday.Time, err = time.Parse(db.DateFormat, in.Birthday)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("EditUserProfile %s date format err:%+v", in.Birthday, err)
			return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorProfile, "生日格式错误")
		}
		if time.Now().Year()-draftData.Birthday.Time.Year() < 16 {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.ListenerErrorProfile, "年龄不符合要求")
		}
		draftData.Birthday.Valid = true
		draftData.Constellation, _ = tool.GetConstellation(in.Birthday)
	}

	draftData.CheckStatus = listenerkey.CheckStatusFirstApplyEdit
	draftData.DraftVersion = 1
	if in.VoiceChatSwitch == 0 {
		draftData.VoiceChatSwitch = db.Enable
	}
	draftData.VoiceChatPrice = listenerkey.DefaultVoiceChatPrice
	if in.TextChatSwitch == 0 {
		draftData.TextChatSwitch = db.Enable
	}
	draftData.TextChatPrice = listenerkey.DefaultTextChatPrice

	if in.AutoReplyNew == "" {
		draftData.AutoReplyNew = listenerkey.AutoReply1
	}
	if in.AutoReplyProcessing == "" {
		draftData.AutoReplyProcessing = listenerkey.AutoReply2
	}
	if in.AutoReplyFinish == "" {
		draftData.AutoReplyFinish = listenerkey.AutoReply3
	}
	_, err = l.svcCtx.ListenerProfileDraftModel.InsertTrans(l.ctx, nil, draftData)
	if err != nil {
		return nil, err
	}
	return &pb.EditListenerProfileDraftResp{DraftVersion: draftData.DraftVersion}, nil
}

func (l *AddOrUpdateListenerProfileDraftLogic) UpdateListenerProfileDraft(draftData *listenerPgModel2.ListenerProfileDraft, in *pb.EditListenerProfileDraftReq) (*pb.EditListenerProfileDraftResp, error) {
	if draftData.DraftVersion != in.DraftVersion {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "draft version not match")
	}
	// 获取需要更新的数据
	newCheckField, isChange := l.getUpdatedData(draftData, in)
	if !isChange {
		return &pb.EditListenerProfileDraftResp{DraftVersion: draftData.DraftVersion}, nil
	}
	// 更新
	err := l.svcCtx.ListenerProfileDraftModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if len(newCheckField) > 0 {
			switch draftData.CheckStatus {
			case listenerkey.CheckStatusFirstApplyEdit, listenerkey.CheckStatusFirstApplyRefuse, listenerkey.CheckStatusFirstApplyChecking: // 首次申请
				draftData.CheckStatus = listenerkey.CheckStatusFirstApplyEdit
			case listenerkey.CheckStatusFirstApplyPass, listenerkey.CheckStatusEditWaitChecking, listenerkey.CheckStatusEditRefuse,
				listenerkey.CheckStatusEditPass: // 成为XXX之后编辑修改
				draftData.CheckStatus = listenerkey.CheckStatusEditWaitChecking
			default:
				return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "wrong check status")
			}
			draftData.DraftVersion += 1
			// 把当前更新的字段 从上次审核失败的字段中移除
			draftData.CheckFailField = stringx.Remove(draftData.CheckFailField, newCheckField...)
			// 更新审核中的字段
			if len(draftData.CheckingField) <= 0 {
				draftData.CheckingField = newCheckField
			} else {
				draftData.CheckingField = tool.CombineStringArray(draftData.CheckingField, newCheckField)
			}
		}

		// 更新pg
		err2 := l.svcCtx.ListenerProfileDraftModel.UpdateTrans(ctx, session, draftData)
		if err2 != nil {
			return err2
		}

		// 不需要审核的资料 当成为XXX以后才能更新
		if tool.IsInt64ArrayExist(draftData.CheckStatus, listenerkey.ListenerCheckStatus) && (in.MaritalStatus != 0 || in.Education != 0 || in.Job != "" || in.TextChatSwitch != 0 || in.VoiceChatSwitch != 0) {
			var pf *listenerPgModel2.ListenerProfile
			pf, err2 = l.svcCtx.ListenerProfileModel.FindOne(ctx, in.ListenerUid)
			if err2 != nil {
				return err2
			}
			var update int
			if in.Education != 0 && in.Education != pf.Education {
				pf.Education = in.Education
				update++
			}
			if in.MaritalStatus != 0 && in.MaritalStatus != pf.MaritalStatus {
				pf.MaritalStatus = in.MaritalStatus
				update++
			}
			if in.Job != "" && in.Job != pf.Job {
				pf.Job = in.Job
				update++
			}
			if in.TextChatSwitch != 0 && in.TextChatSwitch != pf.TextChatSwitch {
				pf.TextChatSwitch = in.TextChatSwitch
				update++
			}
			if in.VoiceChatSwitch != 0 && in.VoiceChatSwitch != pf.VoiceChatSwitch {
				pf.VoiceChatSwitch = in.VoiceChatSwitch
				update++
			}
			if update > 0 {
				err2 = l.svcCtx.ListenerProfileModel.UpdateNoNeedCheckProfile(ctx, session, pf)
				if err2 != nil {
					return err2
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.EditListenerProfileDraftResp{DraftVersion: draftData.DraftVersion}, nil
}

func (l *AddOrUpdateListenerProfileDraftLogic) getUpdatedData(draftData *listenerPgModel2.ListenerProfileDraft, in *pb.EditListenerProfileDraftReq) (newCheckField []string, isChange bool) {
	// 需要审核的资料 先判断是否修改 可以修改 修改完则加入审核中
	if in.NickName != "" && in.NickName != draftData.NickName {
		draftData.NickName = in.NickName
		newCheckField = append(newCheckField, "nickName")
	}
	if in.Avatar != "" && in.Avatar != draftData.Avatar {
		draftData.Avatar = in.Avatar
		newCheckField = append(newCheckField, "avatar")
	}

	if in.Province != "" && in.Province != draftData.Province {
		draftData.Province = in.Province
		newCheckField = append(newCheckField, "province")
	}
	if in.City != "" && in.City != draftData.City {
		draftData.City = in.City
		newCheckField = append(newCheckField, "city")
	}

	if in.Gender != 0 && in.Gender != draftData.Gender {
		draftData.Gender = in.Gender
		newCheckField = append(newCheckField, "gender")
	}
	var draftDataBirthday string
	if draftData.Birthday.Valid {
		draftDataBirthday = draftData.Birthday.Time.Format(db.DateFormat)
	}
	if in.Birthday != "" && in.Birthday != draftDataBirthday {
		var err error
		draftData.Birthday.Time, err = time.Parse(db.DateFormat, in.Birthday)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("UpdateListenerProfileDraft %s date format err:%+v", in.Birthday, err)
		} else {
			draftData.Birthday.Valid = true
			newCheckField = append(newCheckField, "birthday")
			draftData.Constellation, _ = tool.GetConstellation(in.Birthday)
		}
	}

	if in.PhoneNumber != "" && in.PhoneNumber != draftData.PhoneNumber {
		draftData.PhoneNumber = in.PhoneNumber
		newCheckField = append(newCheckField, "phoneNumber")
	}
	if in.ListenerName != "" && in.ListenerName != draftData.ListenerName {
		draftData.ListenerName = in.ListenerName
		newCheckField = append(newCheckField, "listenerName")
	}
	if in.IdNo != "" && in.IdNo != draftData.IdNo {
		draftData.IdNo = in.IdNo
		newCheckField = append(newCheckField, "idNo")
	}
	if in.IdPhoto1 != "" && in.IdPhoto1 != draftData.IdPhoto1 {
		draftData.IdPhoto1 = in.IdPhoto1
		newCheckField = append(newCheckField, "idPhoto1")
	}
	if in.IdPhoto2 != "" && in.IdPhoto2 != draftData.IdPhoto2 {
		draftData.IdPhoto2 = in.IdPhoto2
		newCheckField = append(newCheckField, "idPhoto2")
	}
	if in.IdPhoto3 != "" && in.IdPhoto3 != draftData.IdPhoto3 {
		draftData.IdPhoto3 = in.IdPhoto3
		newCheckField = append(newCheckField, "idPhoto3")
	}

	if len(in.Specialties) > 0 && !tool.IsEqualArrayInt64(in.Specialties, draftData.Specialties) {
		draftData.Specialties = in.Specialties
		newCheckField = append(newCheckField, "specialties")
	}
	if in.Introduction != "" && in.Introduction != draftData.Introduction {
		draftData.Introduction = in.Introduction
		newCheckField = append(newCheckField, "introduction")
	}
	if in.VoiceFile != "" && in.VoiceFile != draftData.VoiceFile {
		draftData.VoiceFile = in.VoiceFile
		newCheckField = append(newCheckField, "voiceFile")
	}
	if in.Experience1 != "" && in.Experience1 != draftData.Experience1 {
		draftData.Experience1 = in.Experience1
		newCheckField = append(newCheckField, "experience1")
	}
	if in.Experience2 != "" && in.Experience2 != draftData.Experience2 {
		draftData.Experience2 = in.Experience2
		newCheckField = append(newCheckField, "experience2")
	}

	if in.OtherPlatformAccount != "" && in.OtherPlatformAccount != draftData.OtherPlatformAccount && in.CertType == listenerkey.CertTypeExperienced {
		draftData.OtherPlatformAccount = in.OtherPlatformAccount
		newCheckField = append(newCheckField, "otherPlatformAccount")
	}
	if in.CertType != 0 && in.CertType != draftData.CertType {
		draftData.CertType = in.CertType
		newCheckField = append(newCheckField, "certType")
	}
	if in.CertFiles1 != "" && in.CertFiles1 != draftData.CertFiles1 {
		draftData.CertFiles1 = in.CertFiles1
		newCheckField = append(newCheckField, "certFiles1")
	}
	if in.CertFiles2 != "" && in.CertFiles2 != draftData.CertFiles2 {
		draftData.CertFiles2 = in.CertFiles2
		newCheckField = append(newCheckField, "certFiles2")
	}
	if in.CertFiles3 != "" && in.CertFiles3 != draftData.CertFiles3 {
		draftData.CertFiles3 = in.CertFiles3
		newCheckField = append(newCheckField, "certFiles3")
	}
	if in.CertFiles4 != "" && in.CertFiles4 != draftData.CertFiles4 {
		draftData.CertFiles4 = in.CertFiles4
		newCheckField = append(newCheckField, "certFiles4")
	}
	if in.CertFiles5 != "" && in.CertFiles5 != draftData.CertFiles5 {
		draftData.CertFiles5 = in.CertFiles5
		newCheckField = append(newCheckField, "certFiles5")
	}

	if in.AutoReplyNew != "" && in.AutoReplyNew != draftData.AutoReplyNew {
		draftData.AutoReplyNew = in.AutoReplyNew
		newCheckField = append(newCheckField, "autoReplyNew")
	}
	if in.AutoReplyProcessing != "" && in.AutoReplyProcessing != draftData.AutoReplyProcessing {
		draftData.AutoReplyProcessing = in.AutoReplyProcessing
		newCheckField = append(newCheckField, "autoReplyProcessing")
	}
	if in.AutoReplyFinish != "" && in.AutoReplyFinish != draftData.AutoReplyFinish {
		draftData.AutoReplyFinish = in.AutoReplyFinish
		newCheckField = append(newCheckField, "autoReplyFinish")
	}

	if in.TextChatPrice != 0 && in.TextChatPrice != draftData.TextChatPrice {
		draftData.TextChatPrice = in.TextChatPrice
		newCheckField = append(newCheckField, "textChatPrice")
	}

	if in.VoiceChatPrice != 0 && in.VoiceChatPrice != draftData.VoiceChatPrice {
		draftData.VoiceChatPrice = in.VoiceChatPrice
		newCheckField = append(newCheckField, "voiceChatPrice")
	}
	if len(newCheckField) > 0 {
		isChange = true
	}
	// 不需要审核的资料
	if in.MaritalStatus != 0 && in.MaritalStatus != draftData.MaritalStatus {
		draftData.MaritalStatus = in.MaritalStatus
		isChange = true
	}
	if in.Education != 0 && in.Education != draftData.Education {
		draftData.Education = in.Education
		isChange = true
	}
	if in.Job != "" && in.Job != draftData.Job {
		draftData.Job = in.Job
		isChange = true
	}
	if in.TextChatSwitch != 0 && in.TextChatSwitch != draftData.TextChatSwitch {
		draftData.TextChatSwitch = in.TextChatSwitch
		isChange = true
	}
	if in.VoiceChatSwitch != 0 && in.VoiceChatSwitch != draftData.VoiceChatSwitch {
		draftData.VoiceChatSwitch = in.VoiceChatSwitch
		isChange = true
	}
	return
}
