package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNeedHelpUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNeedHelpUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNeedHelpUserListLogic {
	return &GetNeedHelpUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取需要帮助用户列表
func (l *GetNeedHelpUserListLogic) GetNeedHelpUserList(in *pb.GetNeedHelpUserListReq) (*pb.GetNeedHelpUserListResp, error) {
	cnt, err := l.svcCtx.UserNeedHelpModel.FindCount(l.ctx, in.Uid, in.ListenerUid, in.Tag, in.State)
	if err != nil {
		return nil, err
	}
	resp := &pb.GetNeedHelpUserListResp{List: make([]*pb.NeedHelpUserData, 0), Sum: cnt}
	if cnt <= 0 {
		return resp, nil
	}
	rs, err := l.svcCtx.UserNeedHelpModel.Find(l.ctx, in.Uid, in.ListenerUid, in.Tag, in.State, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}

	for idx := 0; idx < len(rs); idx++ {
		var val pb.NeedHelpUserData
		_ = copier.Copy(&val, rs[idx])
		val.CreateTime = rs[idx].CreateTime.Format(db.DateTimeFormat)
		val.UpdateTime = rs[idx].UpdateTime.Format(db.DateTimeFormat)
		resp.List = append(resp.List, &val)
	}

	return resp, nil
}
