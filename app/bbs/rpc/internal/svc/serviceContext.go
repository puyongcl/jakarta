package svc

import (
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/bbs/rpc/internal/config"
	"jakarta/app/pgModel/bbsPgModel"
	"jakarta/app/redisModel/listenerRedisModel"
	"jakarta/app/redisModel/userRedisModel"
	"jakarta/common/key/db"
)

type ServiceContext struct {
	Config config.Config

	StoryModel             bbsPgModel.StoryModel
	StoryReplyModel        bbsPgModel.StoryReplyModel
	StoryReplyLikeLogModel bbsPgModel.StoryReplyLikeLogModel
	StoryViewLogModel      bbsPgModel.StoryViewLogModel

	UserRedis     *userRedisModel.UserRedis
	ListenerRedis *listenerRedisModel.ListenerRedis

	KqueueSendSubscribeNotifyMsgClient *kq.Pusher
	KqueueSendWxFwhMsgClient           *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlx.DisableStmtLog()
	rc := c.Redis.NewRedis()
	return &ServiceContext{
		Config: c,

		StoryModel:             bbsPgModel.NewStoryModel(sqlConn, c.Cache),
		StoryReplyModel:        bbsPgModel.NewStoryReplyModel(sqlConn, c.Cache),
		StoryReplyLikeLogModel: bbsPgModel.NewStoryReplyLikeLogModel(sqlConn, c.Cache),
		StoryViewLogModel:      bbsPgModel.NewStoryViewLogModel(sqlConn),

		UserRedis:     userRedisModel.NewUserRedis(rc),
		ListenerRedis: listenerRedisModel.NewListenerRedis(rc),

		KqueueSendSubscribeNotifyMsgClient: kq.NewPusher(c.SendSubscribeNotifyMsgClientConf.Brokers, c.SendSubscribeNotifyMsgClientConf.Topic),
		KqueueSendWxFwhMsgClient:           kq.NewPusher(c.SendWxFwhProgramMsgClientConf.Brokers, c.SendWxFwhProgramMsgClientConf.Topic),
	}
}
