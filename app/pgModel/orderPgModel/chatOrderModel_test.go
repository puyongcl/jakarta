package orderPgModel

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
)

func Test_defaultChatOrderModel_FindCount(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_order?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := newChatOrderModel(sqlConn, c)
	rs, err := tpf.FindCount(context.Background(), 1002, 0, 0, []int64{}, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rs)
}

func Test_defaultChatOrderModel_Find(t *testing.T) {
}
