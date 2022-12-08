package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/listenerPgModel"
	"jakarta/app/pgModel/orderPgModel"
	"jakarta/app/redisModel/listenerRedisModel"
	"jakarta/common/key/db"
	"jakarta/common/key/listenerkey"
	"jakarta/common/key/orderkey"
)

type GenGoodComment struct {
	listenerRedis             *listenerRedisModel.ListenerRedis
	listenerProfileModel      listenerPgModel.ListenerProfileModel
	chatOrderModel            orderPgModel.ChatOrderModel
	ChatOrderPricingPlanModel orderPgModel.ChatOrderPricingPlanModel
	ctx                       context.Context
}

func main2() {
	var orderDataSource, listenerDataSource string
	orderDataSource = "postgres://jakarta:h7su+92Tgscx@10.0.0.9:5432/jakarta_order?sslmode=disable&timezone=Asia/Shanghai"
	listenerDataSource = "postgres://jakarta:h7su+92Tgscx@10.0.0.9:5432/jakarta_listener?sslmode=disable&timezone=Asia/Shanghai"
	var rdc redis.RedisKeyConf
	rdc = redis.RedisKeyConf{RedisConf: redis.RedisConf{
		Host: "10.0.0.5:6379",
		Pass: "h18js0iwoxafaws",
	},
	}
	var c cache.CacheConf
	c = []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "10.0.0.5:6379",
			Pass: "h18js0iwoxafaws",
		},
			Weight: 100,
		},
	}

	start(orderDataSource, listenerDataSource, c, rdc)
}

func main() {
	var orderDataSource, listenerDataSource string
	orderDataSource = "postgres://jakarta:postgres@127.0.0.1:5432/jakarta_order?sslmode=disable&timezone=Asia/Shanghai"
	listenerDataSource = "postgres://jakarta:postgres@127.0.0.1:5432/jakarta_listener?sslmode=disable&timezone=Asia/Shanghai"
	var rdc redis.RedisKeyConf
	rdc = redis.RedisKeyConf{RedisConf: redis.RedisConf{
		Host: "127.0.0.1:36379",
		Pass: "G62m50oigInC30sf",
	},
	}
	var c cache.CacheConf
	c = []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "127.0.0.1:36381",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}

	start(orderDataSource, listenerDataSource, c, rdc)
}

func start(orderDataSource, listenerDataSource string, c cache.CacheConf, rdc redis.RedisKeyConf) {
	sqlConnOrder := sqlx.NewSqlConn(db.PostgresDriverName, orderDataSource)
	sqlConnListener := sqlx.NewSqlConn(db.PostgresDriverName, listenerDataSource)
	rc := rdc.NewRedis()
	ggc := &GenGoodComment{
		listenerRedis:             listenerRedisModel.NewListenerRedis(rc),
		listenerProfileModel:      listenerPgModel.NewListenerProfileModel(sqlConnListener, c),
		chatOrderModel:            orderPgModel.NewChatOrderModel(sqlConnOrder, c),
		ChatOrderPricingPlanModel: orderPgModel.NewChatOrderPricingPlanModel(sqlConnOrder),
		ctx:                       context.Background(),
	}

	ggc.control()
}

func (g *GenGoodComment) control() {
	var pageNo int64 = 1
	var cnt, sum int
	//
	var err error
	for ; ; pageNo++ {
		cnt, err = g.doFix(pageNo, 10)
		if err != nil {
			return
		}
		if cnt <= 0 {
			fmt.Println("StartGenGoodComment sum:", sum)
			return
		}
		sum += cnt
	}
}

func (g *GenGoodComment) doFix(pageNo, pageSize int64) (int, error) {
	rsp, err := g.listenerProfileModel.FindListenerUidRangeUpdateTime(g.ctx, pageNo, pageSize, 0, nil, nil)
	if err != nil {
		fmt.Println("FindRecentActive error:", err)
		return 0, err
	}

	for idx := 0; idx < len(rsp); idx++ {
		err = g.fix(rsp[idx].ListenerUid)
		if err != nil {
			return 0, err
		}
	}
	return len(rsp), nil
}

func (g *GenGoodComment) fix(listenerUid int64) error {
	up := listenerPgModel.AddListenerStat{
		ListenerUid:          listenerUid,
		AddUserCount:         0,
		AddChatDuration:      0,
		AddRatingSum:         0,
		AddFiveStar:          0,
		AddThreeStar:         0,
		AddOneStar:           0,
		AddRefundOrderCnt:    0,
		AddFinishOrderCnt:    0,
		AddPaidOrderCnt:      0,
		AddRepeatPaidUserCnt: 0,
	}
	var err error
	up.AddUserCount, err = g.chatOrderModel.CountPaidUser2(g.ctx, listenerUid)
	if err != nil {
		return err
	}

	up.AddChatDuration, err = g.chatOrderModel.SumOrderBuyUnitCnt(g.ctx, listenerUid)
	if err != nil {
		return err
	}
	up.AddChatDuration = up.AddChatDuration * listenerkey.PerUnitMinute

	up.AddRatingSum, err = g.chatOrderModel.CountOrderRatingCnt(g.ctx, listenerUid, 0)
	if err != nil {
		return err
	}
	up.AddFiveStar, err = g.chatOrderModel.CountOrderRatingCnt(g.ctx, listenerUid, listenerkey.Rating5Star)
	if err != nil {
		return err
	}
	up.AddThreeStar, err = g.chatOrderModel.CountOrderRatingCnt(g.ctx, listenerUid, listenerkey.Rating3Star)
	if err != nil {
		return err
	}
	up.AddOneStar, err = g.chatOrderModel.CountOrderRatingCnt(g.ctx, listenerUid, listenerkey.Rating1Star)
	if err != nil {
		return err
	}
	up.AddRefundOrderCnt, err = g.chatOrderModel.CountOrder(g.ctx, listenerUid, []int64{orderkey.ChatOrderStateFinishRefund8})
	if err != nil {
		return err
	}
	up.AddFinishOrderCnt, err = g.chatOrderModel.CountOrder(g.ctx, listenerUid, []int64{orderkey.ChatOrderStateSettle21})
	if err != nil {
		return err
	}
	up.AddPaidOrderCnt, err = g.chatOrderModel.CountPaidOrderCnt2(g.ctx, listenerUid)
	if err != nil {
		return err
	}
	up.AddRepeatPaidUserCnt, err = g.chatOrderModel.CountRepeatPaidUser2(g.ctx, listenerUid)
	if err != nil {
		return err
	}

	err = g.listenerProfileModel.ResetListenerStat(g.ctx, &up)
	if err != nil {
		return err
	}
	return nil
}
