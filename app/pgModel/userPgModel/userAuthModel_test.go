package userPgModel

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/usercenter/rpc/pb"
	"jakarta/common/key/db"
	"testing"
)

func Test_defaultUserAuthModel_FindByCheckStatus(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_user?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	ua := newUserAuthModel(sqlConn, c)

	datas, err := ua.FindByUserType(context.Background(), 2, 1, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(datas) <= 0 {
		fmt.Println("no data")
		return
	}
	fmt.Println(len(datas))
	rsp := make([]*pb.UserAuth, 0)
	for idx := 0; idx < len(datas); idx++ {
		var val pb.UserAuth
		_ = copier.Copy(&val, datas[idx])
		val.CreateTime = datas[idx].CreateTime.Format(db.DateTimeFormat)
		val.UpdateTime = datas[idx].UpdateTime.Format(db.DateTimeFormat)
		rsp = append(rsp, &val)
	}
	fmt.Println(rsp)
}
