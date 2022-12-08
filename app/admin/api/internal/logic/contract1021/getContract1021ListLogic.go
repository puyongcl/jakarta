package contract1021

import (
	"context"
	"github.com/jinzhu/copier"
	"jakarta/app/pgModel/adminPgModel"
	"jakarta/common/key/db"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContract1021ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContract1021ListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContract1021ListLogic {
	return &GetContract1021ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContract1021ListLogic) GetContract1021List(req *types.ListContract1021Req) (resp *types.ListContract1021Resp, err error) {
	resp = &types.ListContract1021Resp{List: make([]*types.Contract1021Data, 0)}
	resp.Sum, err = l.svcCtx.Contract1021Model.Count(l.ctx)
	if err != nil {
		return nil, err
	}
	var data []*adminPgModel.Contract1021
	data, err = l.svcCtx.Contract1021Model.Find(l.ctx, req.PageNo, req.PageSize)
	if err != nil {
		return nil, err
	}

	if len(data) <= 0 {
		return resp, err
	}

	for idx := 0; idx < len(data); idx++ {
		val := types.Contract1021Data{}
		_ = copier.Copy(&val, data[idx])
		val.StartDate = data[idx].StartDate.Format(db.DateFormat)
		val.EndDate = data[idx].EndDate.Format(db.DateFormat)
		if data[idx].SignTime.Valid {
			val.SignTime = data[idx].SignTime.Time.Format(db.DateFormat)
		}

		resp.List = append(resp.List, &val)
	}
	return
}
