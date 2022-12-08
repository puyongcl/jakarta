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

func Test_defaultListenerProfileDraftModel_FindByFilter(t *testing.T) {
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
	tpf := newListenerProfileDraftModel(sqlConn, c)
	rsp, err := tpf.FindByFilter(context.Background(), []int64{8}, "", 0, 1, 20)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}

func Test_defaultListenerProfileDraftModel_FindByFilterCount(t *testing.T) {
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
	tpf := newListenerProfileDraftModel(sqlConn, c)
	rsp, err := tpf.FindByFilterCount(context.Background(), []int64{8}, "", 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}
