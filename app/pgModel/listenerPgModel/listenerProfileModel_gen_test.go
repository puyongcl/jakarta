package listenerPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
)

func Test_defaultListenerProfileModel(t *testing.T) {
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
	tpf := newListenerProfileModel(sqlConn, c)
	rs, err := tpf.FindRecommendListenerList(context.Background(), 1, 10, 101, 0, 0, 0, 0, 0, 0, []int64{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rs)
}
