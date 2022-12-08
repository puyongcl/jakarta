package paymentPgModel

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
	"time"
)

func TestPg(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_payment?sslmode=disable"
	//c := []cache.NodeConf{
	//	{RedisConf: redis.RedisConf{
	//		Host: "192.168.1.12:36379",
	//		Pass: "G62m50oigInC30sf",
	//	},
	//		Weight: 100,
	//	},
	//}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewThirdPaymentFlowModel(sqlConn)
	data := ThirdPaymentFlow{
		FlowNo:         "f0001",
		Uid:            10000,
		PayMode:        "paymode",
		TradeType:      "tradetype",
		TradeState:     "FAIL",
		PayAmount:      10000,
		TransactionId:  "t0001",
		TradeStateDesc: "fail",
		OrderId:        "s0001",
		OrderType:      2,
		PayStatus:      1,
		PayTime: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}
	rs, err := tpf.Insert(context.Background(), &data)
	fmt.Println(rs)
	fmt.Println(err)
	return
}

func TestPgQuery(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_payment?sslmode=disable"
	//c := []cache.NodeConf{
	//	{RedisConf: redis.RedisConf{
	//		Host: "192.168.1.12:36379",
	//		Pass: "G62m50oigInC30sf",
	//	},
	//		Weight: 100,
	//	},
	//}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewThirdPaymentFlowModel(sqlConn)

	resp, err := tpf.FindOneByOrderId(context.Background(), "s0001", []int64{1, 2, 3})
	fmt.Println(resp)
	fmt.Println(err)
}

func TestPgUpdate(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_payment?sslmode=disable"
	//c := []cache.NodeConf{
	//	{RedisConf: redis.RedisConf{
	//		Host: "192.168.1.12:36379",
	//		Pass: "G62m50oigInC30sf",
	//	},
	//		Weight: 100,
	//	},
	//}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewThirdPaymentFlowModel(sqlConn)
	data := ThirdPaymentFlow{
		FlowNo:         "f0001",
		Uid:            10000,
		PayMode:        "paymode",
		TradeType:      "tradetype",
		TradeState:     "PASS",
		PayAmount:      10000,
		TransactionId:  "t0001",
		TradeStateDesc: "PASS",
		OrderId:        "s0001",
		OrderType:      2,
		PayStatus:      1,
		PayTime: sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		},
	}

	err := tpf.Update(context.Background(), &data)
	fmt.Println(err)
}
