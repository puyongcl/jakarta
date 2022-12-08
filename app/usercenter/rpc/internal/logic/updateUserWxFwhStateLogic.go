package logic

import (
	"context"
	"fmt"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/common/key/db"
	"jakarta/common/xerr"

	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserWxFwhStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserWxFwhStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserWxFwhStateLogic {
	return &UpdateUserWxFwhStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新用户服务号的关注情况
func (l *UpdateUserWxFwhStateLogic) UpdateUserWxFwhState(in *pb.UpdateUserWxFwhStateReq) (*pb.UpdateUserWxFwhStateResp, error) {
	if in.OpenId == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数为空")
	}

	switch in.State {
	case db.Enable:
		uwi := new(userPgModel.UserWechatInfo)
		*uwi = userPgModel.UserWechatInfo{
			FwhOpenid: in.OpenId,
			UnionId:   in.UnionId,
			FwhState:  in.State,
		}
		_, err := l.svcCtx.UserWechatInfoModel.InsertOrUpdateFwhTrans(l.ctx, uwi)
		if err != nil {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.DbError, fmt.Sprintf("FWH subscribe db user_wechat_info Insert err:%+v,unionid:%s", err, in.UnionId))
		}

	case db.Disable:
		err := l.svcCtx.UserWechatInfoModel.UpdateFwhUnsubscribe(l.ctx, in.OpenId, in.State)
		if err != nil {
			return nil, xerr.NewErrCodeMsg(xerr.DbError, err.Error())
		}
	default:
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, fmt.Sprintf("参数错误：%d", in.State))
	}

	return &pb.UpdateUserWxFwhStateResp{}, nil
}
