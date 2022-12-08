package svc

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/imPgModel"
	"jakarta/common/key/db"
	"jakarta/common/third_party/tim"

	_ "github.com/lib/pq"
	"jakarta/app/im/rpc/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	TimClient   *tim.RestApiClient
	MiniProgram *miniprogram.MiniProgram
	Wx          *wechat.Wechat
	Wxfwh       *officialaccount.OfficialAccount

	ImMsgLogModel imPgModel.ImMsgLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	tc := tim.InitTimRestApiClient(c.TimConf.SDKAPPID, c.TimConf.ADMINID, c.TimConf.ADMINSIGN)
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, c.DB.DataSource)
	sqlx.DisableStmtLog()
	wx, wf, mp := InitWechat(&c)
	return &ServiceContext{
		Config:      c,
		TimClient:   tc,
		Wx:          wx,
		Wxfwh:       wf,
		MiniProgram: mp,

		ImMsgLogModel: imPgModel.NewImMsgLogModel(sqlConn),
	}
}
