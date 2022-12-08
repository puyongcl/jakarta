package user

import (
	"context"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/key/tencentcloudkey"
	"jakarta/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCloudConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCloudConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCloudConfigLogic {
	return &GetCloudConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCloudConfigLogic) GetCloudConfig(req *types.GetCloudConfigReq) (resp *types.GetCloudConfigResp, err error) {
	rs, err := l.svcCtx.TencentAPIClient.GetTempSecret()
	if err != nil {
		return nil, err
	}
	if rs.Credentials == nil {
		return nil, xerr.NewGrpcErrCodeMsg(xerr.ThirdPartRequestError, "cos temp key error")
	}
	resp = &types.GetCloudConfigResp{
		Credentials: &types.Credentials{
			TmpSecretID:  rs.Credentials.TmpSecretID,
			TmpSecretKey: rs.Credentials.TmpSecretKey,
			SessionToken: rs.Credentials.SessionToken,
		},
		ExpiredTime:    rs.ExpiredTime,
		Expiration:     rs.Expiration,
		StartTime:      rs.StartTime,
		RequestId:      rs.RequestId,
		Bucket:         l.svcCtx.Config.TencentConf.Bucket,
		Region:         l.svcCtx.Config.TencentConf.Region,
		CdnBasePath:    tencentcloudkey.CDNBasePath,
		BucketBasePath: tencentcloudkey.BucketBasePath,
	}
	return
}
