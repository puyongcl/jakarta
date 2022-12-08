package svc

import (
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/redisModel/userRedisModel"
	"jakarta/app/usercenter/rpc/internal/config"
	"jakarta/common/key/db"
	"jakarta/common/third_party/tim"
)

type ServiceContext struct {
	Config      config.Config
	RedisClient *redis.Redis

	UserAuthModel                userPgModel.UserAuthModel
	UserProfileModel             userPgModel.UserProfileModel
	UserLoginStateModel          userPgModel.UserLoginStateModel
	UserStatModel                userPgModel.UserStatModel
	UserBlacklistModel           userPgModel.UserBlacklistModel
	UserReportModel              userPgModel.UserReportModel
	UserNeedHelpModel            userPgModel.UserNeedHelpModel
	UserWechatInfoModel          userPgModel.UserWechatInfoModel
	UserChannelCallbackModel     userPgModel.UserChannelCallbackModel
	UserAdviserConversationModel userPgModel.UserAdviserConversationModel

	UserRedis *userRedisModel.UserRedis

	TimClient *tim.RestApiClient

	AsynqClient *asynq.Client

	KqueueUploadUserEventClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlx.DisableStmtLog()
	rc := c.Redis.NewRedis()
	tc := tim.InitTimRestApiClient(c.TimConf.SDKAPPID, c.TimConf.ADMINID, c.TimConf.ADMINSIGN)
	return &ServiceContext{
		Config:                       c,
		RedisClient:                  rc,
		TimClient:                    tc,
		UserAuthModel:                userPgModel.NewUserAuthModel(sqlConn, c.Cache),
		UserProfileModel:             userPgModel.NewUserProfileModel(sqlConn, c.Cache),
		UserLoginStateModel:          userPgModel.NewUserLoginStateModel(sqlConn, c.Cache),
		UserStatModel:                userPgModel.NewUserStatModel(sqlConn, c.Cache),
		UserBlacklistModel:           userPgModel.NewUserBlacklistModel(sqlConn),
		UserReportModel:              userPgModel.NewUserReportModel(sqlConn),
		UserNeedHelpModel:            userPgModel.NewUserNeedHelpModel(sqlConn),
		UserWechatInfoModel:          userPgModel.NewUserWechatInfoModel(sqlConn, c.Cache),
		UserChannelCallbackModel:     userPgModel.NewUserChannelCallbackModel(sqlConn, c.Cache),
		UserAdviserConversationModel: userPgModel.NewUserAdviserConversationModel(sqlConn, c.Cache),

		UserRedis: userRedisModel.NewUserRedis(rc),

		AsynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAsynq.Host, Password: c.RedisAsynq.Pass}),

		KqueueUploadUserEventClient: kq.NewPusher(c.UploadUserEventClientConf.Brokers, c.UploadUserEventClientConf.Topic),
	}
}
