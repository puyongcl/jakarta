package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/usercenter/rpc/internal/svc"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/db"
	"jakarta/common/key/userkey"
	"jakarta/common/xerr"
)

type BlockUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockUserLogic {
	return &BlockUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  拉黑操作
func (l *BlockUserLogic) BlockUser(in *pb.BlockUserReq) (*pb.BlockUserResp, error) {
	if in.Uid == in.TargetUid {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "same uid")
	}
	if in.Action == db.Enable {
		cnt, err := l.svcCtx.UserRedis.CountBlacklist(l.ctx, in.Uid)
		if err != nil {
			return nil, err
		}
		if cnt > userkey.MaxBlacklistCnt {
			return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "超出最大限制")
		}
	}
	data := new(userPgModel.UserBlacklist)
	data = &userPgModel.UserBlacklist{
		Uid:            in.Uid,
		Id:             fmt.Sprintf(db.DBUidId, in.Uid, in.TargetUid),
		TargetUid:      in.TargetUid,
		TargetNickName: in.TargetNickName,
		TargetAvatar:   in.TargetAvatar,
		State:          in.Action,
	}

	//
	err := l.svcCtx.UserBlacklistModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		var err2 error
		_, err2 = l.svcCtx.UserBlacklistModel.InsertOrUpdate(ctx, session, data)
		if err2 != nil {
			return err2
		}

		switch in.Action {
		case db.Enable:
			err2 = l.svcCtx.TimClient.AddBlacklist(in.Uid, in.TargetUid)
			if err2 != nil {
				return err2
			}
			err2 = l.svcCtx.UserRedis.AddBlacklist(l.ctx, in.Uid, in.TargetUid)
			if err2 != nil {
				return err2
			}
		case db.Disable:
			err2 = l.svcCtx.TimClient.DeleteBlacklist(in.Uid, in.TargetUid)
			if err2 != nil {
				return err2
			}
			err2 = l.svcCtx.UserRedis.DelBlacklist(l.ctx, in.Uid, in.TargetUid)
			if err2 != nil {
				return err2
			}
		default:
			return xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.BlockUserResp{}, nil
}
