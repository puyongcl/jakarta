package svc

import (
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"jakarta/app/admin/api/internal/config"
	"jakarta/app/chat/rpc/chat"
	"jakarta/app/im/rpc/im"
	"jakarta/app/listener/rpc/listener"
	"jakarta/app/order/rpc/order"
	"jakarta/app/payment/rpc/payment"
	"jakarta/app/pgModel/adminPgModel"
	"jakarta/app/statistic/rpc/statistic"
	"jakarta/app/usercenter/rpc/usercenter"
	"jakarta/common/key/db"
)

type ServiceContext struct {
	Config config.Config

	AdminLogModel             adminPgModel.AdminLogModel
	AdminMenuPermissionsModel adminPgModel.AdminMenuPermissionsModel
	Contract1021Model         adminPgModel.Contract1021Model

	ListenerRpc   listener.Listener
	UsercenterRpc usercenter.Usercenter
	OrderRpc      order.Order
	ChatRpc       chat.Chat
	StatRpc       statistic.Statistic
	ImRpc         im.Im
	PaymentRpc    payment.Payment

	RedisClient *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.AdminDataSource)
	rc := c.Redis.NewRedis()
	return &ServiceContext{
		Config:                    c,
		AdminLogModel:             adminPgModel.NewAdminLogModel(sqlConn),
		AdminMenuPermissionsModel: adminPgModel.NewAdminMenuPermissionsModel(sqlConn),
		Contract1021Model:         adminPgModel.NewContract1021Model(sqlConn),

		ListenerRpc:   listener.NewListener(zrpc.MustNewClient(c.ListenerRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		ChatRpc:       chat.NewChat(zrpc.MustNewClient(c.ChatRpcConf)),
		StatRpc:       statistic.NewStatistic(zrpc.MustNewClient(c.StatRpcConf)),
		ImRpc:         im.NewIm(zrpc.MustNewClient(c.ImRpcConf)),
		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		RedisClient:   rc,
	}
}
