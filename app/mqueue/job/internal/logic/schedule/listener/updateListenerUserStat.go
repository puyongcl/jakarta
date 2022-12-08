package listener

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"
	pbChat "jakarta/app/chat/rpc/pb"
	pbListener "jakarta/app/listener/rpc/pb"
	"jakarta/app/mqueue/job/internal/svc"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"time"
)

type AutoUpdateListenerUserStatHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAutoUpdateListenerUserStatHandler(svcCtx *svc.ServiceContext) *AutoUpdateListenerUserStatHandler {
	return &AutoUpdateListenerUserStatHandler{
		svcCtx: svcCtx,
	}
}

//
func (l *AutoUpdateListenerUserStatHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	// 获取最近活跃的XXX
	var pageNo int64 = 1
	now := time.Now()
	end := now.Format(db.DateTimeFormat)
	start := now.AddDate(0, 0, -listenerkey.AutoUpdateListenerUserStatRangeDay).Format(db.DateTimeFormat)
	var cnt int
	var rsp *pbListener.FindListenerListRangeByUpdateTimeResp
	var err error
	for ; ; pageNo++ {
		rsp, err = l.svcCtx.ListenerRpc.FindListenerListRangeByUpdateTime(ctx, &pbListener.FindListenerListRangeByUpdateTimeReq{
			PageNo:   pageNo,
			PageSize: 10,
			Start:    start,
			End:      end,
		})
		if err != nil {
			logx.WithContext(ctx).Errorf("AutoUpdateListenerUserStatHandler FindListenerListRangeByUpdateTime err:%+v", err)
			return err
		}

		if len(rsp.Listener) <= 0 {
			if cnt > 0 {
				logx.WithContext(ctx).Infof("AutoUpdateListenerUserStatHandler exit. process listener cnt:%d", cnt)
			}
			break
		}

		// 更新聊天进入用户数
		var in1 pbChat.UpdateTodayEnterChatUserCntReq
		var in2 pbListener.UpdateTodayListenerUserStatReq
		in1.ListenerUid = make([]int64, 0)
		for idx := 0; idx < len(rsp.Listener); idx++ {
			in1.ListenerUid = append(in1.ListenerUid, rsp.Listener[idx].ListenerUid)
			in2.ListenerUid = append(in2.ListenerUid, rsp.Listener[idx].ListenerUid)
		}
		_, err = l.svcCtx.ChatRpc.UpdateTodayEnterChatUserCnt(ctx, &in1)
		if err != nil {
			return err
		}
		// 更新推荐、访问个人资料用户数
		_, err = l.svcCtx.ListenerRpc.UpdateTodayListenerUserStat(ctx, &in2)
		if err != nil {
			return err
		}

		cnt += len(rsp.Listener)
	}

	logx.WithContext(ctx).Infof("AutoUpdateListenerUserStatHandler done")
	return nil
}
