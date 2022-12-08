package listener

import (
	"context"
	"jakarta/common/key/listenerkey"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

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

func (l *GetDefineBusinessConfigLogic) GetDefineBusinessConfig(req *types.GetBusinessConfigReq) (resp *types.GetBusinessConfigResp, err error) {
	resp = &types.GetBusinessConfigResp{}
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
	return
}
