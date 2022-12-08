package logic

import (
	"context"
	"database/sql"
	"fmt"
	listenerPgModel2 "jakarta/app/pgModel/listenerPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/listener/rpc/internal/svc"
	"jakarta/app/listener/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateListenerUserStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateListenerUserStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateListenerUserStatLogic {
	return &UpdateListenerUserStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  更新XXX与用户的交互情况
func (l *UpdateListenerUserStatLogic) UpdateListenerUserStat(in *pb.UpdateListenerUserStatReq) (*pb.UpdateListenerUserStatResp, error) {
	tt, err := time.ParseInLocation(db.DateTimeFormat, in.Time, time.Local)
	if err != nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "time is empty")
	}

	for idx := 0; idx < len(in.ListenerUid); idx++ {
		if in.Uid == in.ListenerUid[idx] {
			continue
		}

		switch in.Event {
		case listenerkey.ListenerUserEventRecommend:
			data := listenerPgModel2.ListenerUserRecommendStat{
				Uid: in.Uid,
				RecommendTime: sql.NullTime{
					Time:  tt,
					Valid: true,
				},
				RecommendCnt: 1,
			}

			data.Id = fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid[idx])
			data.ListenerUid = in.ListenerUid[idx]

			_, err = l.svcCtx.ListenerUserRecommendStatModel.InsertOrUpdateRecommend(l.ctx, &data)
			if err != nil {
				return nil, err
			}

		case listenerkey.ListenerUserEventView:
			data := listenerPgModel2.ListenerUserViewStat{
				Uid: in.Uid,
				ViewTime: sql.NullTime{
					Time:  tt,
					Valid: true,
				},
				ViewCnt: 1,
			}

			data.Id = fmt.Sprintf(db.DBUidId, in.Uid, in.ListenerUid[idx])
			data.ListenerUid = in.ListenerUid[idx]
			_, err = l.svcCtx.ListenerUserViewStatModel.InsertOrUpdateView(l.ctx, &data)
			if err != nil {
				return nil, err
			}

		default:
		}
	}

	return &pb.UpdateListenerUserStatResp{}, nil
}
