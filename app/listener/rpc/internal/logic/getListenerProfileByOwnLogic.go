package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetListenerProfileByOwnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetListenerProfileByOwnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetListenerProfileByOwnLogic {
	return &GetListenerProfileByOwnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  XXX获取自己的资料
func (l *GetListenerProfileByOwnLogic) GetListenerProfileByOwn(in *pb.GetListenerProfileByOwnReq) (*pb.GetListenerProfileByOwnResp, error) {
	rs, err := l.svcCtx.ListenerProfileDraftModel.FindOne(l.ctx, in.ListenerUid)
	if err != nil && err != listenerPgModel.ErrNotFound {
		return nil, err
	}
	var pf pb.ListenerSeeOwnProfile
	pf.Specialties = make([]int64, 0)
	pf.CheckFailField = make([]string, 0)
	pf.CheckingField = make([]string, 0)
	if rs != nil {
		_ = copier.Copy(&pf, rs)
		if rs.Birthday.Valid {
			pf.Birthday = rs.Birthday.Time.Format(db.DateFormat)
		}
	}
	return &pb.GetListenerProfileByOwnResp{Profile: &pf}, nil
}
