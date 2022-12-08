package logic

import (
	"context"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/userkey"
	"time"

	"jakarta/app/statistic/rpc/internal/svc"
	"jakarta/app/statistic/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLifeTimeValueStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLifeTimeValueStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLifeTimeValueStatLogic {
	return &GetLifeTimeValueStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取统计近多少日的用户在昨日累计数据
func (l *GetLifeTimeValueStatLogic) GetLifeTimeValueStat(in *pb.GetLifeTimeValueStatReq) (*pb.GetLifeTimeValueStatResp, error) {
	var startDate, endDate time.Time
	var err error
	startDate, err = time.ParseInLocation(db.DateFormat2, in.StartDate, time.Local)
	if err != nil {
		return nil, err
	}
	endDate, err = time.ParseInLocation(db.DateFormat2, in.StartDate, time.Local)
	if err != nil {
		return nil, err
	}
	var curStartDate, curEndDate time.Time
	curStartDate = startDate
	curEndDate = startDate.AddDate(0, 0, 1)

	resp := pb.GetLifeTimeValueStatResp{List: make([]*pb.LifeTimeValueStat, 0)}
	for {
		if curStartDate.Day() == endDate.Day() {
			break
		}
		var val pb.LifeTimeValueStat
		if in.UserFlag == 1 {
			err = l.queryNewUserData(&curStartDate, &curEndDate, in.Channel, &val)
		} else {
			err = l.queryUserData(&curStartDate, &curEndDate, in.Channel, &val)
		}
		if err != nil {
			return nil, err
		}

		curEndDate = curStartDate.AddDate(0, 0, 1)
		curStartDate = curEndDate
		resp.List = append(resp.List, &val)
	}
	return &pb.GetLifeTimeValueStatResp{}, nil
}

func (l *GetLifeTimeValueStatLogic) queryUserData(startDate, endDate *time.Time, channel string, stat *pb.LifeTimeValueStat) (err error) {
	stat.Date = startDate.Format(db.DateFormat2)
	//
	stat.UserCnt, err = l.svcCtx.UserLoginLogModel.Count(l.ctx, startDate, endDate, channel, userkey.UserTypeNormalUser)
	if err != nil {
		return
	}
	stat.PaidUserCnt, err = l.svcCtx.ChatOrderModel.CountPaidUserCntRangeCreateTime(l.ctx, startDate, endDate, channel)
	if err != nil {
		return
	}
	stat.RefundOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserOrderRangeCreateTime(l.ctx, startDate, endDate, channel, orderkey.ChatOrderRefundState)
	if err != nil {
		return
	}
	stat.PaidOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserPaidOrderCntRangeCreateTime(l.ctx, startDate, endDate, channel)
	if err != nil {
		return
	}
	stat.RepeatPaidUserCnt, err = l.svcCtx.ChatOrderModel.CountRepeatPaidUserRangeCreateTime(l.ctx, startDate, endDate, channel)
	if err != nil {
		return
	}
	stat.CommentOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserCommentOrderRangeCreateTime(l.ctx, startDate, endDate, channel)
	if err != nil {
		return
	}
	stat.FiveStarOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserCommentOrderByStarRangeCreateTime(l.ctx, startDate, endDate, listenerkey.Rating5Star, channel)
	if err != nil {
		return
	}
	stat.ThreeStarOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserCommentOrderByStarRangeCreateTime(l.ctx, startDate, endDate, listenerkey.Rating3Star, channel)
	if err != nil {
		return
	}
	stat.OneStarOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserCommentOrderByStarRangeCreateTime(l.ctx, startDate, endDate, listenerkey.Rating1Star, channel)
	if err != nil {
		return
	}
	stat.TextChatOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserOrderCntByTypeRangeCreateTime(l.ctx, startDate, endDate, orderkey.ListenerOrderTypeTextChat, channel)
	if err != nil {
		return
	}
	stat.VoiceChatOrderCnt, err = l.svcCtx.ChatOrderModel.CountUserOrderCntByTypeRangeCreateTime(l.ctx, startDate, endDate, orderkey.ListenerOrderTypeVoiceChat, channel)
	if err != nil {
		return
	}
	stat.PaidAmountSum, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTime(l.ctx, startDate, endDate, channel)
	if err != nil {
		return
	}

	err = l.queryUserLtv(startDate, channel, stat)
	if err != nil {
		return
	}
	return
}

func (l *GetLifeTimeValueStatLogic) queryUserLtv(startDate *time.Time, channel string, stat *pb.LifeTimeValueStat) (err error) {
	// 查询当天登陆的用户
	var uids []int64
	var am int64
	var ltvEndTime time.Time
	var pageNo int64
	for pageNo = 1; ; pageNo++ {
		uids, err = l.svcCtx.UserLoginLogModel.FindUid(l.ctx, startDate, &ltvEndTime, channel, userkey.UserTypeNormalUser, pageNo, 10)
		if err != nil {
			return
		}
		if len(uids) <= 0 {
			break
		}

		ltvEndTime = startDate.AddDate(0, 0, 2)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv1Day += am
		}
		ltvEndTime = startDate.AddDate(0, 0, 4)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv3Day += am
		}
		ltvEndTime = startDate.AddDate(0, 0, 8)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv7Day += am
		}
		ltvEndTime = startDate.AddDate(0, 0, 15)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv14Day += am
		}
		ltvEndTime = startDate.AddDate(0, 0, 22)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv21Day += am
		}
		ltvEndTime = startDate.AddDate(0, 0, 31)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv30Day += am
		}
		ltvEndTime = startDate.AddDate(0, 0, 61)
		if ltvEndTime.Unix() < time.Now().Unix() {
			am, err = l.svcCtx.ChatOrderModel.SumUserPaidAmountRangeCreateTimeLtv(l.ctx, startDate, &ltvEndTime, channel, uids)
			if err != nil {
				return
			}
			stat.Ltv60Day += am
		}
	}

	return
}

func (l *GetLifeTimeValueStatLogic) queryNewUserData(startDate, endDate *time.Time, channel string, stat *pb.LifeTimeValueStat) (err error) {
	var startUid, endUid int64
	stat.Date = startDate.Format(db.DateFormat2)
	startUid, err = l.svcCtx.UserAuthModel.FindFirstUid(l.ctx, startDate, endDate)
	if err != nil {
		return err
	}
	endUid, err = l.svcCtx.UserAuthModel.FindFirstUid(l.ctx, startDate, endDate)
	if err != nil {
		return err
	}
	stat.UserCnt, err = l.svcCtx.UserAuthModel.CountNewUser(l.ctx, startDate, endDate, channel)
	if err != nil {
		return
	}
	stat.PaidUserCnt, err = l.svcCtx.ChatOrderModel.CountPaidNewUserCntRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, channel)
	if err != nil {
		return
	}
	stat.RefundOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserOrderRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, channel, orderkey.ChatOrderRefundState)
	if err != nil {
		return
	}
	stat.PaidOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserPaidOrderCntRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, channel)
	if err != nil {
		return
	}
	stat.RepeatPaidUserCnt, err = l.svcCtx.ChatOrderModel.CountRepeatPaidNewUserRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, channel)
	if err != nil {
		return
	}
	stat.CommentOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserCommentOrderRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, channel)
	if err != nil {
		return
	}
	stat.FiveStarOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserCommentOrderByStarRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, listenerkey.Rating5Star, channel)
	if err != nil {
		return
	}
	stat.ThreeStarOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserCommentOrderByStarRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, listenerkey.Rating3Star, channel)
	if err != nil {
		return
	}
	stat.OneStarOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserCommentOrderByStarRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, listenerkey.Rating1Star, channel)
	if err != nil {
		return
	}
	stat.TextChatOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserOrderCntByTypeRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, orderkey.ListenerOrderTypeTextChat, channel)
	if err != nil {
		return
	}
	stat.VoiceChatOrderCnt, err = l.svcCtx.ChatOrderModel.CountNewUserOrderCntByTypeRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, orderkey.ListenerOrderTypeVoiceChat, channel)
	if err != nil {
		return
	}
	stat.PaidAmountSum, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, endDate, startUid, endUid, channel)
	if err != nil {
		return
	}

	ltvEndTime := startDate.AddDate(0, 0, 2)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv1Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	ltvEndTime = startDate.AddDate(0, 0, 4)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv3Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	ltvEndTime = startDate.AddDate(0, 0, 8)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv7Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	ltvEndTime = startDate.AddDate(0, 0, 15)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv14Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	ltvEndTime = startDate.AddDate(0, 0, 22)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv21Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	ltvEndTime = startDate.AddDate(0, 0, 31)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv30Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	ltvEndTime = startDate.AddDate(0, 0, 61)
	if ltvEndTime.Unix() < time.Now().Unix() {
		stat.Ltv60Day, err = l.svcCtx.ChatOrderModel.SumNewUserPaidAmountRangeCreateTime(l.ctx, startDate, &ltvEndTime, startUid, endUid, channel)
		if err != nil {
			return
		}
	}
	return
}
