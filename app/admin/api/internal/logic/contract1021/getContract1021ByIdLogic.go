package contract1021

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/pgModel/adminPgModel"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"jakarta/common/xerr"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContract1021ByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContract1021ByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContract1021ByIdLogic {
	return &GetContract1021ByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetContract1021ByIdLogic) GetContract1021ById(req *types.QueryContract1021ByIdReq) (resp *types.QueryContract1021ByIdResp, err error) {
	if req.ContractId == "" {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}
	// 加分布式锁
	rkey := fmt.Sprintf(rediskey.RedisLockGenContract1021, req.ContractId)
	rl := redis.NewRedisLock(l.svcCtx.RedisClient, rkey)
	rl.SetExpire(2)
	b, err := rl.AcquireCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RedisLockFail, "操作太过频繁")
	}

	defer func() {
		_, err2 := rl.ReleaseCtx(l.ctx)
		if err2 != nil {
			logx.WithContext(l.ctx).Errorf("RedisLock %s release err:%+v", rkey, err2)
			return
		}
	}()

	data, err := l.svcCtx.Contract1021Model.FindOne(l.ctx, req.ContractId)
	if err != nil && err != adminPgModel.ErrNotFound {
		return nil, err
	}
	resp = &types.QueryContract1021ByIdResp{}
	err = nil
	if data != nil {
		_ = copier.Copy(resp, data)
		resp.StartDate = data.StartDate.Format(db.DateFormat)
		resp.EndDate = data.EndDate.Format(db.DateFormat)
		if data.SignTime.Valid {
			resp.SignTime = data.SignTime.Time.Format(db.DateFormat)
		}
	}

	return
}
