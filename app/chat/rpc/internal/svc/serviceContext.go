package svc

import (
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/chat/rpc/internal/config"
	"jakarta/app/pgModel/chatPgModel"
	"jakarta/app/redisModel/chatRedisModel"
	"jakarta/common/key/db"
)

type ServiceContext struct {
	Config config.Config

	RedisClient                 *redis.Redis
	ChatRedis                   *chatRedisModel.ChatRedis
	VoiceChatLogModel           chatPgModel.VoiceChatLogModel
	ChatBalanceModel            chatPgModel.ChatBalanceModel
	ChatBalanceLogModel         chatPgModel.ChatBalanceLogModel
	ChatBalanceModelRO          chatPgModel.ChatBalanceModel
	ListenerVoiceChatStateModel chatPgModel.ListenerVoiceChatStateModel
	UserVoiceChatStateModel     chatPgModel.UserVoiceChatStateModel
	UserListenerRelationModel   chatPgModel.UserListenerRelationModel

	KqueueUpdateChatStatClient *kq.Pusher
	KqueueSendDefineMsgClient  *kq.Pusher
	KqueueFirstEnterChatClient *kq.Pusher
}

func NewServiceContext(c config.Config) *ServiceContext {
	rc := c.Redis.NewRedis()
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlConnRo := sqlx.NewSqlConn(db.PostgresDriverName, c.DBRO.DataSource)
	sqlx.DisableStmtLog()
	return &ServiceContext{
		Config:                      c,
		RedisClient:                 rc,
		ChatRedis:                   chatRedisModel.NewChatRedis(rc),
		VoiceChatLogModel:           chatPgModel.NewVoiceChatLogModel(sqlConn, c.Cache),
		ChatBalanceModel:            chatPgModel.NewChatBalanceModel(sqlConn, c.Cache),
		ChatBalanceModelRO:          chatPgModel.NewChatBalanceModel(sqlConnRo, c.Cache),
		ChatBalanceLogModel:         chatPgModel.NewChatBalanceLogModel(sqlConn),
		ListenerVoiceChatStateModel: chatPgModel.NewListenerVoiceChatStateModel(sqlConn, c.Cache),
		UserVoiceChatStateModel:     chatPgModel.NewUserVoiceChatStateModel(sqlConn, c.Cache),
		UserListenerRelationModel:   chatPgModel.NewUserListenerRelationModel(sqlConn),

		KqueueUpdateChatStatClient: kq.NewPusher(c.KqUpdateChatStatConf.Brokers, c.KqUpdateChatStatConf.Topic),
		KqueueSendDefineMsgClient:  kq.NewPusher(c.KqSendDefineMsgConf.Brokers, c.KqSendDefineMsgConf.Topic),
		KqueueFirstEnterChatClient: kq.NewPusher(c.KqFirstEnterChatConf.Brokers, c.KqFirstEnterChatConf.Topic),
	}
}
