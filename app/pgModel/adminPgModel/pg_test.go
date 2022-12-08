package adminPgModel

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
)

type metaData struct {
	KeepAlive bool     `json:"keepAlive"`
	Title     string   `json:"title"`
	Role      []string `json:"role"`
	Icon      string   `json:"icon"`
}
type adminMenu struct {
	Path      string   `json:"path"`
	Name      string   `json:"name"`
	Redirect  string   `json:"redirect"`
	Meta      metaData `json:"meta"`
	Component string   `json:"component"`
}
type listAdminMenuReq struct {
	AdminUid int64 `json:"adminUid"` // 管理员id
	Uid      int64 `json:"uid"`      // 要查询的管理员id
	PageNo   int64 `json:"pageNo"`
	PageSize int64 `json:"pageSize"`
	Menu1Id  int64 `json:"menu1Id"` // 一级菜单id 当菜单id 都为0 是查询1级菜单
	Menu2Id  int64 `json:"menu2Id"` // 二级菜单id 当只有2级菜单为0 是查询1级菜单对应的2级菜单
}

type listAdminMenuResp struct {
	List []adminMenu `json:"list"`
	Sum  int64       `json:"sum"`
}

func TestPg(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_admin?sslmode=disable"
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := newAdminLogModel(sqlConn)
	req := listAdminMenuReq{
		AdminUid: 10000,
		Uid:      1003,
		PageNo:   1,
		PageSize: 10,
		Menu1Id:  1,
		Menu2Id:  2,
	}
	resp := listAdminMenuResp{
		List: []adminMenu{
			{
				Path:     "test",
				Name:     "11",
				Redirect: "11",
				Meta: metaData{
					KeepAlive: false,
					Title:     "test",
					Role:      []string{"t1", "t2"},
					Icon:      "ss",
				},
				Component: "asdad",
			},
			{
				Path:     "test",
				Name:     "11",
				Redirect: "11",
				Meta: metaData{
					KeepAlive: false,
					Title:     "test",
					Role:      []string{"t1", "t2"},
					Icon:      "ss",
				},
				Component: "asdad",
			},
		},
		Sum: 0,
	}
	buf, err := json.Marshal(&req)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf2, err := json.Marshal(&resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := AdminLog{
		AdminUid:  10000,
		Request:   string(buf),
		Response:  string(buf2),
		RoutePath: "test",
	}
	rs, err := tpf.Insert(context.Background(), &data)
	fmt.Println(rs)
	fmt.Println(err)

	data = AdminLog{
		AdminUid:  10000,
		Request:   "",
		Response:  "",
		RoutePath: "test",
	}
	rs, err = tpf.Insert(context.Background(), &data)
	fmt.Println(rs)
	fmt.Println(err)
	return
}

func TestPgQuery(t *testing.T) {
	ds := "postgres://jakarta:postgres@192.168.1.12:5432/jakarta_listener?sslmode=disable"
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := newAdminLogModel(sqlConn)

	resp, err := tpf.FindOne(context.Background(), 1)
	fmt.Println(resp)
	fmt.Println(err)
}

func TestPgCount(t *testing.T) {
	ds := "postgres://jakarta:postgres@api2.domain.com:5432/jakarta_admin?sslmode=disable&timezone=Asia/Shanghai"
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := newContract1021Model(sqlConn)

	resp, err := tpf.Count(context.Background())
	fmt.Println(resp)
	fmt.Println(err)
}
