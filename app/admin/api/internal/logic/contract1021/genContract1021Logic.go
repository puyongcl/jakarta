package contract1021

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"jakarta/app/pgModel/adminPgModel"
	"jakarta/common/ctxdata"
	"jakarta/common/key/db"
	"jakarta/common/key/rediskey"
	"jakarta/common/tool"
	"jakarta/common/xerr"
	"time"

	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenContract1021Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenContract1021Logic(ctx context.Context, svcCtx *svc.ServiceContext) *GenContract1021Logic {
	return &GenContract1021Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenContract1021Logic) GenContract1021(req *types.GenContract1021Req) (resp *types.GenContract1021Resp, err error) {
	uid := ctxdata.GetUidFromCtx(l.ctx)
	if uid == 0 {
		logx.WithContext(l.ctx).Errorf("GenContract1021Logic uid is empty")
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "参数错误")
	}

	if l.svcCtx.Config.Mode != service.ProMode {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.RequestParamError, "不支持")
	}

	// 加分布式锁
	rkey := fmt.Sprintf(rediskey.RedisLockUser, uid)
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

	strT := req.UserName + req.Idno + req.StartDate + req.EndDate
	tmd5 := tool.Md5ByString(strT)

	data := adminPgModel.Contract1021{
		ContractId:  tmd5,
		UserName:    req.UserName,
		PhoneNumber: req.PhoneNumber,
		Amount:      req.Amount,
		Idno:        req.Idno,
	}

	data.StartDate, err = time.ParseInLocation(db.DateFormat, req.StartDate, time.Local)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GenContract1021Logic ParseInLocation req:%+v err:%+v", req, err)
		return nil, err
	}
	data.EndDate, err = time.ParseInLocation(db.DateFormat, req.EndDate, time.Local)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("GenContract1021Logic ParseInLocation req:%+v err:%+v", req, err)
		return nil, err
	}
	_, err = l.svcCtx.Contract1021Model.Insert(l.ctx, &data)
	if err != nil {
		return nil, err
	}
	resp = &types.GenContract1021Resp{ContractId: data.ContractId}
	return
}
