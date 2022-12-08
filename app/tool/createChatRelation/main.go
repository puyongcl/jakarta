package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/app/pgModel/chatPgModel"
	"jakarta/common/key/db"
)

type CreateNewData struct {
	UserListenerRelationModel chatPgModel.UserListenerRelationModel
	ChatBalanceModel          chatPgModel.ChatBalanceModel
	ctx                       context.Context
}

func main() {
	var chatDataSource string
	chatDataSource = "postgres://jakarta:h7su+92Tgscx@10.0.0.9:5432/jakarta_chat?sslmode=disable&timezone=Asia/Shanghai"
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

	start(chatDataSource, c, rdc)
}

func main2() {
	var chatDataSource string
	chatDataSource = "postgres://jakarta:postgres@127.0.0.1:5432/jakarta_chat?sslmode=disable&timezone=Asia/Shanghai"
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

	start(chatDataSource, c, rdc)
}

func start(chatDataSource string, c cache.CacheConf, rdc redis.RedisKeyConf) {
	sqlConnOrder := sqlx.NewSqlConn(db.PostgresDriverName, chatDataSource)

	ggc := &CreateNewData{
		UserListenerRelationModel: chatPgModel.NewUserListenerRelationModel(sqlConnOrder),
		ChatBalanceModel:          chatPgModel.NewChatBalanceModel(sqlConnOrder, c),
		ctx:                       context.Background(),
	}
	ggc.control()
}

func (g *CreateNewData) control() {
	var pageNo int64 = 1
	var cnt, sum int
	//
	var err error
	for ; ; pageNo++ {
		cnt, err = g.loopCreate(pageNo, 10)
		if err != nil {
			return
		}
		if cnt <= 0 {
			fmt.Println("create UserListenerRelationModel sum:", sum)
			return
		}
		sum += cnt
	}
}

func (g *CreateNewData) loopCreate(pageNo, pageSize int64) (int, error) {
	rsp, err := g.ChatBalanceModel.FindByUidAndState(g.ctx, 0, 0, 0, pageNo, pageSize)
	if err != nil {
		fmt.Println("FindByUidAndState error:", err)
		return 0, err
	}

	for idx := 0; idx < len(rsp); idx++ {
		err = g.doCreate(rsp[idx].Uid, rsp[idx].ListenerUid)
		if err != nil {
			return 0, err
		}
	}
	return len(rsp), nil
}

func (g *CreateNewData) doCreate(uid, listenerUid int64) error {
	id := fmt.Sprintf(db.DBUidId, uid, listenerUid)
	_, err := g.UserListenerRelationModel.FindOne(g.ctx, id)
	if err != nil && err != chatPgModel.ErrNotFound {
		fmt.Println("doCreate FindOne id err", id, err)
		return err
	}
	if err == chatPgModel.ErrNotFound {
		ulr := chatPgModel.UserListenerRelation{
			Id:          id,
			Uid:         uid,
			ListenerUid: listenerUid,
			TotalScore:  0,
		}
		_, err = g.UserListenerRelationModel.Insert(g.ctx, &ulr)
		if err != nil {
			fmt.Println("doCreate Insert id err", id, err)
			return err
		}
	}
	return nil
}
