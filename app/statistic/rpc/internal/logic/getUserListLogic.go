package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *pb.GetUserListReq) (*pb.GetUserListResp, error) {
	// 只能查普通用户
	in.UserType = userkey.UserTypeNormalUser
	// 查询指定uid
	if in.Uid != 0 {
		var val pb.UserDetail
		auth, err := l.svcCtx.UserAuthModel.FindOne(l.ctx, in.Uid)
		if err != nil && err != userPgModel.ErrNotFound {
			return nil, err
		}
		if auth == nil {
			resp := pb.GetUserListResp{
				Sum:  0,
				List: []*pb.UserDetail{&val},
			}
			return &resp, nil
		}
		stat, err := l.svcCtx.UserStatModel.FindOne(l.ctx, in.Uid)
		if err != nil {
			return nil, err
		}
		pf, err := l.svcCtx.UserProfileModel.FindOne(l.ctx, in.Uid)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&val, auth)
		_ = copier.Copy(&val, stat)
		_ = copier.Copy(&val, pf)

		val.OpenId = auth.AuthKey

		resp := pb.GetUserListResp{
			Sum:  1,
			List: []*pb.UserDetail{&val},
		}

		return &resp, nil
	}
	// 查询指定openid
	if in.AuthKey != "" {
		if in.AuthType == "" {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "auth type为空")
		}
		var val pb.UserDetail
		auth, err := l.svcCtx.UserAuthModel.FindOneByAuthKeyAuthType2(l.ctx, in.AuthKey, in.AuthType)
		if err != nil && err != userPgModel.ErrNotFound {
			return nil, err
		}
		if auth == nil {
			resp := pb.GetUserListResp{
				Sum:  0,
				List: []*pb.UserDetail{&val},
			}
			return &resp, nil
		}
		stat, err := l.svcCtx.UserStatModel.FindOne(l.ctx, auth.Uid)
		if err != nil {
			return nil, err
		}
		pf, err := l.svcCtx.UserProfileModel.FindOne(l.ctx, auth.Uid)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&val, auth)
		_ = copier.Copy(&val, stat)
		_ = copier.Copy(&val, pf)

		val.OpenId = auth.AuthKey

		resp := pb.GetUserListResp{
			Sum:  1,
			List: []*pb.UserDetail{&val},
		}

		return &resp, nil
	}
	var start, end *time.Time
	if in.CreateTimeStart != "" {
		st, err := time.ParseInLocation(db.DateTimeFormat, in.CreateTimeStart, time.Local)
		if err != nil {
			return nil, err
		}
		start = &st
	}
	if in.CreateTimeEnd != "" {
		st, err := time.ParseInLocation(db.DateTimeFormat, in.CreateTimeEnd, time.Local)
		if err != nil {
			return nil, err
		}
		end = &st
	}

	// 是否付费
	if in.IsPaidUser != 0 {
		var cnt int64
		var err error
		cnt, err = l.svcCtx.JoinUserModel.FindUserListCount(l.ctx, start, end, in.IsPaidUser == db.Enable, in.UserType, in.Channel)
		if err != nil {
			return nil, err
		}

		resp := &pb.GetUserListResp{List: make([]*pb.UserDetail, 0), Sum: cnt}
		rsp, err := l.svcCtx.JoinUserModel.FindUserList(l.ctx, start, end, in.IsPaidUser == db.Enable, in.UserType, in.Channel, in.PageNo, in.PageSize)
		if err != nil {
			return nil, err
		}
		if len(rsp) <= 0 {
			return resp, nil
		}
		for idx := 0; idx < len(rsp); idx++ {
			var val pb.UserDetail
			_ = copier.Copy(&val, rsp[idx])
			var pf *userPgModel.UserProfile
			pf, err = l.svcCtx.UserProfileModel.FindOne(l.ctx, rsp[idx].Uid)
			if err != nil {
				return nil, err
			}
			_ = copier.Copy(&val, pf)
			val.CreateTime = pf.CreateTime.Format(db.DateTimeFormat)
			if rsp[idx].FreeTime.Valid {
				val.FreeTime = rsp[idx].FreeTime.Time.Format(db.DateTimeFormat)
			}
			val.OpenId = rsp[idx].AuthKey

			resp.List = append(resp.List, &val)
		}

		return resp, nil
	}

	// 没有付费条件
	var cnt int64
	var err error
	cnt, err = l.svcCtx.UserAuthModel.FindUserByCreateTimeCount(l.ctx, start, end, in.UserType, in.Channel)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetUserListResp{List: make([]*pb.UserDetail, 0), Sum: cnt}

	rsp, err := l.svcCtx.UserAuthModel.FindUserByCreateTime(l.ctx, start, end, in.UserType, in.Channel, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(rsp) <= 0 {
		return resp, nil
	}
	for idx := 0; idx < len(rsp); idx++ {
		var val pb.UserDetail
		_ = copier.Copy(&val, rsp[idx])
		//
		var pf *userPgModel.UserProfile
		pf, err = l.svcCtx.UserProfileModel.FindOne(l.ctx, rsp[idx].Uid)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&val, pf)
		val.CreateTime = pf.CreateTime.Format(db.DateTimeFormat)

		//
		var stat *userPgModel.UserStat
		stat, err = l.svcCtx.UserStatModel.FindOne(l.ctx, rsp[idx].Uid)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&val, stat)

		if rsp[idx].FreeTime.Valid {
			val.FreeTime = rsp[idx].FreeTime.Time.Format(db.DateTimeFormat)
		}
		val.OpenId = rsp[idx].AuthKey

		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
