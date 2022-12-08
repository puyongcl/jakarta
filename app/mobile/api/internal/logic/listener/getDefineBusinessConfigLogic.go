package listener

import (
	"context"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
	"jakarta/common/key/userkey"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDefineBusinessConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDefineBusinessConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDefineBusinessConfigLogic {
	return &GetDefineBusinessConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDefineBusinessConfigLogic) GetBannerAndFilterConfigReq(req *types.GetDefineBusinessConfigReq) (resp *types.GetDefineBusinessConfigResp, err error) {
	return GetBanner(), nil
}

func GetBanner() *types.GetDefineBusinessConfigResp {
	//
	resp := &types.GetDefineBusinessConfigResp{}

	// banner大图
	resp.UserBanner = []*types.Banner{
		{
			Id:   userkey.TopBannerId1,
			Name: userkey.TopBannerId1Name,
			Pic:  userkey.TopBannerId1Pic,
			Url:  userkey.TopBannerId1Url,
		},
		{
			Id:   userkey.TopBannerId2,
			Name: userkey.TopBannerId2Name,
			Pic:  userkey.TopBannerId2Pic,
			Url:  userkey.TopBannerId2Url,
		},
		{
			Id:   userkey.TopBannerId3,
			Name: userkey.TopBannerId3Name,
			Pic:  userkey.TopBannerId3Pic,
			Url:  userkey.TopBannerId3Url,
		},
	}

	resp.ListenerBanner = []*types.Banner{
		{
			Id:   listenerkey.TopBannerId1,
			Name: listenerkey.TopBannerId1Name,
			Pic:  listenerkey.TopBannerId1Pic,
			Url:  listenerkey.TopBannerId1Url,
		},
		//{
		//	Id:   listenerkey.TopBannerId2,
		//	Name: listenerkey.TopBannerId2Name,
		//	Pic:  listenerkey.TopBannerId2Pic,
		//	Url:  listenerkey.TopBannerId2Url,
		//},
	}

	// 专业领域
	resp.Specialties = make([]*types.Banner, 0)
	for idx1 := 0; idx1 < len(listenerkey.SpecialtiesLevelOneId); idx1++ {
		var val types.Banner
		val.Id = listenerkey.SpecialtiesLevelOneId[idx1]
		val.Pic = listenerkey.SpecialtiesPic[idx1]
		val.Name = listenerkey.Specialties[val.Id]
		val.Child = make([]*types.Pair, 0)
		for idx2 := 0; idx2 < len(listenerkey.SpecialtiesLevelTwoId[idx1]); idx2++ {
			var val2 types.Pair
			val2.Id = listenerkey.SpecialtiesLevelTwoId[idx1][idx2]
			val2.Name = listenerkey.Specialties[val2.Id]
			val.Child = append(val.Child, &val2)
		}
		resp.Specialties = append(resp.Specialties, &val)
	}

	// 筛选条件
	// 聊天类型
	resp.ChatTypeFilter = []*types.Pair{
		{
			orderkey.ListenerOrderTypeTextChat,
			"文字",
		},
		{
			orderkey.ListenerOrderTypeVoiceChat,
			"通话",
		},
	}

	// 性别
	resp.GenderFileter = []*types.Pair{
		{
			listenerkey.GenderMale,
			"男",
		},
		{
			listenerkey.GenderFemale,
			"女",
		},
	}

	// 年龄
	resp.AgeFilter = []*types.Pair{
		{
			listenerkey.AgeRange1,
			"60后",
		},
		{
			listenerkey.AgeRange2,
			"70后",
		},
		{
			listenerkey.AgeRange3,
			"80后",
		},
		{
			listenerkey.AgeRange4,
			"90后",
		},
	}
	// 工作状态
	resp.WorkStateFilter = []*types.Pair{
		{
			listenerkey.ListenerWorkStateWorking,
			"接单中",
		},
		{
			listenerkey.ListenerWorkStateRestingManual,
			"休息中",
		},
	}
	// 排序字段 默认正序
	resp.SortOrderFilter = []*types.Pair{
		{
			listenerkey.ListenerSortOrderDefault,
			"综合排序",
		},
		{
			listenerkey.ListenerSortOrderRatingStar,
			"服务满意率",
		},
		{
			listenerkey.ListenerSortOrderRepeatCustomer,
			"回头客人数",
		},
		{
			listenerkey.ListenerSortOrderChatMinute,
			"服务时长",
		},
	}

	// 评价标签
	resp.CommentTag = make([]*types.Banner, orderkey.CommentTagCnt)
	for idx := 0; idx < len(orderkey.CommentStar); idx++ {
		id := orderkey.CommentStar[idx]
		name, ok := orderkey.CommentTag[id]
		if ok {
			resp.CommentTag[idx] = &types.Banner{
				Id:    int64(id),
				Name:  name,
				Pic:   "",
				Child: make([]*types.Pair, 0),
				Url:   "",
			}
			for idx2 := orderkey.CommentStar[idx]*100 + 1; ; idx2++ {
				cname, ok2 := orderkey.CommentTag[idx2]
				if !ok2 {
					break
				}
				resp.CommentTag[idx].Child = append(resp.CommentTag[idx].Child, &types.Pair{
					Id:   int64(idx2),
					Name: cname,
				})
			}
		}
	}

	// 退款原因
	resp.RefundReasonTag = make([]*types.Pair, orderkey.RefundReasonTagCnt)
	for idx := 1; idx < orderkey.RefundReasonTagCnt+1; idx++ {
		name, ok := orderkey.RefundReasonTag[idx]
		if !ok {
			break
		}
		resp.RefundReasonTag[idx-1] = &types.Pair{
			Id:   int64(idx),
			Name: name,
		}
	}

	// 举报标签
	resp.ReportTag = make([]*types.Pair, len(userkey.ReportTag))
	for idx := 1; idx < len(userkey.ReportTag)+1; idx++ {
		name, ok := userkey.ReportTagText[idx]
		if !ok {
			break
		}
		resp.ReportTag[idx-1] = &types.Pair{
			Id:   int64(idx),
			Name: name,
		}
	}

	// XX顾问常用语
	resp.AdviserWords = listenerkey.AdviserWords
	return resp
}
