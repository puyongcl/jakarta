package svc

import (
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/app/pgModel/statPgModel"
	"jakarta/app/pgModel/userPgModel"
	"jakarta/app/statistic/rpc/internal/config"
	"jakarta/common/key/db"
)

type ServiceContext struct {
	Config config.Config

	UserAuthModel          userPgModel.UserAuthModel
	UserProfileModel       userPgModel.UserProfileModel
	UserLoginStateModel    userPgModel.UserLoginStateModel
	UserStatModel          userPgModel.UserStatModel
	UserBlacklistModel     userPgModel.UserBlacklistModel
	UserReportModel        userPgModel.UserReportModel
	UserNeedHelpModel      userPgModel.UserNeedHelpModel
	UserWechatInfoModel    userPgModel.UserWechatInfoModel
	JoinUserModel          *userPgModel.JoinUserModel
	StatDailyModel         statPgModel.StatDailyModel
	ChatOrderModel         orderPgModel.ChatOrderModel
	UserLoginLogModel      statPgModel.UserLoginLogModel
	UserSelectSpecTagModel statPgModel.UserSelectSpecTagModel
	ListenerProfileModel   listenerPgModel.ListenerProfileModel
	AdultQuizEcrModel      statPgModel.AdultQuizEcrModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConnUserRo := sqlx.NewSqlConn(db.PostgresDriverName, c.DBRO.UserDataSource)
	sqlConnOrderRo := sqlx.NewSqlConn(db.PostgresDriverName, c.DBRO.OrderDataSource)
	sqlListenerRo := sqlx.NewSqlConn(db.PostgresDriverName, c.DBRO.ListenerDataSource)

	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlx.DisableStmtLog()

	return &ServiceContext{
		UserAuthModel:          userPgModel.NewUserAuthModel(sqlConnUserRo, c.Cache),
		UserProfileModel:       userPgModel.NewUserProfileModel(sqlConnUserRo, c.Cache),
		UserLoginStateModel:    userPgModel.NewUserLoginStateModel(sqlConnUserRo, c.Cache),
		UserStatModel:          userPgModel.NewUserStatModel(sqlConnUserRo, c.Cache),
		UserBlacklistModel:     userPgModel.NewUserBlacklistModel(sqlConnUserRo),
		UserReportModel:        userPgModel.NewUserReportModel(sqlConnUserRo),
		UserNeedHelpModel:      userPgModel.NewUserNeedHelpModel(sqlConnUserRo),
		UserWechatInfoModel:    userPgModel.NewUserWechatInfoModel(sqlConnUserRo, c.Cache),
		JoinUserModel:          userPgModel.NewJoinUserModel(sqlConnUserRo),
		StatDailyModel:         statPgModel.NewStatDailyModel(sqlConn),
		ChatOrderModel:         orderPgModel.NewChatOrderModel(sqlConnOrderRo, c.Cache),
		UserLoginLogModel:      statPgModel.NewUserLoginLogModel(sqlConn),
		UserSelectSpecTagModel: statPgModel.NewUserSelectSpecTagModel(sqlConn),
		ListenerProfileModel:   listenerPgModel.NewListenerProfileModel(sqlListenerRo, c.Cache),
		AdultQuizEcrModel:      statPgModel.NewAdultQuizEcrModel(sqlConn),
		Config:                 c,
	}
}
