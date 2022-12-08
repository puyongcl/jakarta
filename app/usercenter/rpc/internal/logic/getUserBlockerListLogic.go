package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/common/key/db"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserBlockerListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserBlockerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserBlockerListLogic {
	return &GetUserBlockerListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取拉黑用户列表
func (l *GetUserBlockerListLogic) GetUserBlockerList(in *pb.GetUserBlockListReq) (*pb.GetUserBlockListResp, error) {
	resp := &pb.GetUserBlockListResp{List: make([]*pb.BlockUserInfo, 0)}
	rs, err := l.svcCtx.UserBlacklistModel.Find(l.ctx, in.Uid, 0, db.Enable, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(rs) > 0 {
		for idx := 0; idx < len(rs); idx++ {
			var val pb.BlockUserInfo
			_ = copier.Copy(&val, rs[idx])
			resp.List = append(resp.List, &val)
		}
	}

	return resp, nil
}
