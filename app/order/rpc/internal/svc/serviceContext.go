package svc

import (
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/order/rpc/internal/config"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/app/redisModel/imRedisModel"
	"jakarta/app/redisModel/statisticRedisModel"
	"jakarta/common/key/db"
)

type ServiceContext struct {
	Config config.Config

	ChatOrderModel            orderPgModel.ChatOrderModel
	ChatOrderStatusLogModel   orderPgModel.ChatOrderStatusLogModel
	ChatOrderPricingPlanModel orderPgModel.ChatOrderPricingPlanModel
	ChatOrderStatModel        orderPgModel.ChatOrderStatModel

	AsynqClient *asynq.Client

	KqueueUpdateOrderActionClient  *kq.Pusher
	KqueueSendDefineMsgClient      *kq.Pusher
	KqSendSubscribeNotifyMsgClient *kq.Pusher

	RedisClient *redis.Redis
	StatRedis   *statisticRedisModel.StatisticRedis
	IMRedis     *imRedisModel.IMRedis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlx.DisableStmtLog()
	rc := c.Redis.NewRedis()
	return &ServiceContext{
		ChatOrderModel:            orderPgModel.NewChatOrderModel(sqlConn, c.Cache),
		ChatOrderStatusLogModel:   orderPgModel.NewChatOrderStatusLogModel(sqlConn),
		ChatOrderPricingPlanModel: orderPgModel.NewChatOrderPricingPlanModel(sqlConn),
		ChatOrderStatModel:        orderPgModel.NewChatOrderStatModel(sqlConn, c.Cache),
		Config:                    c,

		RedisClient: rc,
		StatRedis:   statisticRedisModel.NewStatisticRedis(rc),
		IMRedis:     imRedisModel.NewIMRedis(rc),

		AsynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisAsynq.Host, Password: c.RedisAsynq.Pass}),

		KqueueUpdateOrderActionClient:  kq.NewPusher(c.KqUpdateOrderActionConf.Brokers, c.KqUpdateOrderActionConf.Topic),
		KqueueSendDefineMsgClient:      kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
		KqSendSubscribeNotifyMsgClient: kq.NewPusher(c.KqSendSubscribeNotifyMsgClientConf.Brokers, c.KqSendSubscribeNotifyMsgClientConf.Topic),
	}
}
