package svc

import (
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/listener/rpc/internal/config"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/app/redisModel/imRedisModel"
	"jakarta/app/redisModel/listenerRedisModel"
	"jakarta/app/redisModel/statisticRedisModel"
	"jakarta/app/redisModel/userRedisModel"
	"jakarta/common/key/db"
	"jakarta/common/third_party/hfbfcash"
	"jakarta/common/third_party/tencentcloud"
	"jakarta/common/third_party/tim"
)

type ServiceContext struct {
	Config                         config.Config
	ListenerProfileModel           listenerPgModel.ListenerProfileModel
	ListenerProfileDraftModel      listenerPgModel.ListenerProfileDraftModel
	ListenerWalletModel            listenerPgModel.ListenerWalletModel
	ListenerWalletFlowModel        listenerPgModel.ListenerWalletFlowModel
	ListenerRemarkUserModel        listenerPgModel.ListenerRemarkUserModel
	ListenerBankCardModel          listenerPgModel.ListenerBankCardModel
	ListenerWordsModel             listenerPgModel.ListenerWordsModel
	ListenerUserRecommendStatModel listenerPgModel.ListenerUserRecommendStatModel
	ListenerUserViewStatModel      listenerPgModel.ListenerUserViewStatModel
	ListenerDashboardStatModel     listenerPgModel.ListenerDashboardStatModel
	ListenerStatAverageModel       listenerPgModel.ListenerStatAverageModel
	ListenerContractModel          listenerPgModel.ListenerContractModel

	RedisClient   *redis.Redis
	ListenerRedis *listenerRedisModel.ListenerRedis
	StatRedis     *statisticRedisModel.StatisticRedis
	UserRedis     *userRedisModel.UserRedis
	ImRedis       *imRedisModel.IMRedis

	KqueueSendDefineMsgClient      *kq.Pusher
	KqSubscribeNotifyMsgClient     *kq.Pusher
	KqUpdateListenerUserStatClient *kq.Pusher
	KqSendHelloWhenUserLoginClient *kq.Pusher

	TimClient *tim.RestApiClient

	TencentCosClient *tencentcloud.TencentCosClient

	HfbfCashClient *hfbfcash.RestHfbfCashClient
	//AsynqClient *asynq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlx.DisableStmtLog()
	rc := c.Redis.NewRedis()
	tc := tim.InitTimRestApiClient(c.TimConf.SDKAPPID, c.TimConf.ADMINID, c.TimConf.ADMINSIGN)

	return &ServiceContext{
		RedisClient: rc,
		Config:      c,
		TimClient:   tc,

		ListenerProfileModel:           listenerPgModel.NewListenerProfileModel(sqlConn, c.Cache),
		ListenerProfileDraftModel:      listenerPgModel.NewListenerProfileDraftModel(sqlConn, c.Cache),
		ListenerWalletModel:            listenerPgModel.NewListenerWalletModel(sqlConn, c.Cache),
		ListenerWalletFlowModel:        listenerPgModel.NewListenerWalletFlowModel(sqlConn),
		ListenerRemarkUserModel:        listenerPgModel.NewListenerRemarkUserModel(sqlConn, c.Cache),
		ListenerBankCardModel:          listenerPgModel.NewListenerBankCardModel(sqlConn, c.Cache),
		ListenerWordsModel:             listenerPgModel.NewListenerWordsModel(sqlConn, c.Cache),
		ListenerUserRecommendStatModel: listenerPgModel.NewListenerUserRecommendStatModel(sqlConn),
		ListenerUserViewStatModel:      listenerPgModel.NewListenerUserViewStatModel(sqlConn),
		ListenerDashboardStatModel:     listenerPgModel.NewListenerDashboardStatModel(sqlConn, c.Cache),
		ListenerStatAverageModel:       listenerPgModel.NewListenerStatAverageModel(sqlConn),
		ListenerContractModel:          listenerPgModel.NewListenerContractModel(sqlConn),

		ListenerRedis: listenerRedisModel.NewListenerRedis(rc),
		UserRedis:     userRedisModel.NewUserRedis(rc),
		StatRedis:     statisticRedisModel.NewStatisticRedis(rc),
		ImRedis:       imRedisModel.NewIMRedis(rc),

		KqueueSendDefineMsgClient:      kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
		KqSubscribeNotifyMsgClient:     kq.NewPusher(c.KqSubscribeNotifyMsgConf.Brokers, c.KqSubscribeNotifyMsgConf.Topic),
		KqUpdateListenerUserStatClient: kq.NewPusher(c.KqUpdateListenerUserStatConf.Brokers, c.KqUpdateListenerUserStatConf.Topic),
		KqSendHelloWhenUserLoginClient: kq.NewPusher(c.KqSendHelloWhenUserLoginConf.Brokers, c.KqSendHelloWhenUserLoginConf.Topic),

		TencentCosClient: tencentcloud.InitCosClient(c.TencentConf.SecretId, c.TencentConf.SecretKey, c.TencentConf.Bucket, c.TencentConf.Region),

		HfbfCashClient: hfbfcash.InitCqCashClient(c.HfbfCashConf.AppId, c.HfbfCashConf.AppSecret),
		//AsynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAsynq.Host, Password: c.RedisAsynq.Pass}),
	}
}
