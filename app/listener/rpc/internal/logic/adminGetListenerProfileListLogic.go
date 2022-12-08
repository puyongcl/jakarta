package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetListenerProfileListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminGetListenerProfileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetListenerProfileListLogic {
	return &AdminGetListenerProfileListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  管理员接口 获取XXX列表
func (l *AdminGetListenerProfileListLogic) AdminGetListenerProfileList(in *pb.GetListenerProfileListReq) (*pb.GetListenerProfileListResp, error) {
	cnt, err := l.svcCtx.ListenerProfileDraftModel.FindByFilterCount(l.ctx, in.CheckStatus, in.ListenerName, in.ListenerUid, in.CertType)
	if err != nil {
		return nil, err
	}
	if cnt <= 0 {
		return &pb.GetListenerProfileListResp{List: []*pb.CheckListenerProfile{}, Sum: 0}, nil
	}
	datas, err := l.svcCtx.ListenerProfileDraftModel.FindByFilter(l.ctx, in.CheckStatus, in.ListenerName, in.ListenerUid, in.CertType, in.PageNo, in.PageSize)
	if err != nil {
		return nil, err
	}
	if len(datas) <= 0 {
		return &pb.GetListenerProfileListResp{}, nil
	}
	rsp := make([]*pb.CheckListenerProfile, 0)
	for idx := 0; idx < len(datas); idx++ {
		var val pb.AdminSeeListenerProfileDraft
		_ = copier.Copy(&val, datas[idx])
		val.CreateTime = datas[idx].CreateTime.Format(db.DateTimeFormat)
		val.UpdateTime = datas[idx].UpdateTime.Format(db.DateTimeFormat)
		if datas[idx].Birthday.Valid {
			val.Birthday = datas[idx].Birthday.Time.Format(db.DateFormat)
		}

		if tool.IsInt64ArrayExist(datas[idx].CheckStatus, []int64{listenerkey.CheckStatusFirstApplyEdit, listenerkey.CheckStatusFirstApplyChecking, listenerkey.CheckStatusFirstApplyRefuse}) {
			rsp = append(rsp, &pb.CheckListenerProfile{
				Draft: &val,
			})
		} else {
			var pf *listenerPgModel2.ListenerProfile
			pf, err = l.svcCtx.ListenerProfileModel.FindOne(l.ctx, datas[idx].ListenerUid)
			if err != nil && err != listenerPgModel2.ErrNotFound {
				return nil, err
			}
			var val2 pb.AdminSeeListenerProfile
			if pf != nil {
				_ = copier.Copy(&val2, pf)
				val2.CreateTime = pf.CreateTime.Format(db.DateTimeFormat)
				val2.UpdateTime = pf.UpdateTime.Format(db.DateTimeFormat)
				val2.Birthday = pf.Birthday.Format(db.DateFormat)
			}

			rsp = append(rsp, &pb.CheckListenerProfile{
				Draft:   &val,
				Profile: &val2,
			})
		}

	}

	return &pb.GetListenerProfileListResp{List: rsp, Sum: cnt}, nil
}
