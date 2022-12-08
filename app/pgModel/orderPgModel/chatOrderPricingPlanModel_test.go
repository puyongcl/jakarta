package orderPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
)

func Test_defaultChatOrderPricingPlanModel_FindPriceConfig(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_listener?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	ua := newChatOrderPricingPlanModel(sqlConn, c)
	rsp, err := ua.FindPriceConfig(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}
