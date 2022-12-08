package orderPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	ds := "postgres://jakarta:postgres@api.domain.com:5432/jakarta_order?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewChatOrderModel(sqlConn, c)
	tst, err := time.ParseInLocation(db.DateTimeFormat, "2022-07-20 17:01:41", time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}
	tsp, err := time.ParseInLocation(db.DateTimeFormat, "2022-07-20 17:10:00", time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tpf.UpdateOrderUse(context.Background(), "CV2022072016281291448083", 8, &tst, &tsp, 4, []int64{})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGetRecentGoodComment(t *testing.T) {
	ds := "postgres://jakarta:postgres@api2.domain.com:5432/jakarta_order?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewChatOrderModel(sqlConn, c)
	rsp, err := tpf.FindGoodComment(context.Background(), 1, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(rsp))
}
