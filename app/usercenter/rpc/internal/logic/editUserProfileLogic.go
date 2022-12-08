package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/tool"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditUserProfileLogic {
	return &EditUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EditUserProfileLogic) EditUserProfile(in *pb.EditUserProfileReq) (*pb.EditUserProfileResp, error) {
	var data *userPgModel.UserProfile
	var err error
	data, err = l.svcCtx.UserProfileModel.FindOne(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	var updateCnt int
	var imNick string
	var imAvatar string
	if in.Nickname != "" {
		data.Nickname = in.Nickname
		imNick = in.Nickname
		updateCnt++
	}
	if in.Avatar != "" {
		data.Avatar = in.Avatar
		imAvatar = in.Avatar
		updateCnt++
	}
	if in.Birthday != "" {
		data.Birthday.Time, err = time.Parse(db.DateFormat, in.Birthday)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("EditUserProfile %s date format err:%+v", in.Birthday, err)
		} else {
			data.Birthday.Valid = true
			data.Constellation, _ = tool.GetConstellation(in.Birthday)
		}
		updateCnt++
	}
	if in.Gender != 0 {
		data.Gender = in.Gender
		updateCnt++
	}
	if in.Introduction != "" {
		data.Introduction = in.Introduction
		updateCnt++
	}
	if in.PhoneNumber != "" {
		data.PhoneNumber = in.PhoneNumber
		updateCnt++
	}

	if updateCnt > 0 {
		err = l.svcCtx.UserProfileModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
			err2 := l.svcCtx.UserProfileModel.Update2(l.ctx, session, data)
			if err2 != nil {
				return err2
			}

			// 更新IM
			err2 = l.svcCtx.TimClient.UpdateProfile(in.Uid, imNick, imAvatar)
			if err2 != nil {
				return err2
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	var up pb.UserProfile
	_ = copier.Copy(&up, &data)
	if data.Birthday.Valid {
		up.Birthday = data.Birthday.Time.Format(db.DateFormat)
		up.Age = tool.GetAge2(data.Birthday.Time)
	}

	return &pb.EditUserProfileResp{User: &up}, nil
}
