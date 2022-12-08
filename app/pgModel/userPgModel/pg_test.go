package userPgModel

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"testing"
	"time"
)

func TestFindUserList(t *testing.T) {
	ds := "postgres://jakarta:postgres@api.domain.com:5432/jakarta_user?sslmode=disable"

	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewJoinUserModel(sqlConn)
	start, err := time.ParseInLocation(db.DateTimeFormat, "2022-08-01 17:01:41", time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}
	end, err := time.ParseInLocation(db.DateTimeFormat, "2022-08-22 17:10:00", time.Local)
	if err != nil {
		fmt.Println(err)
		return
	}
	var rsp []*UserListDetail
	rsp, err = tpf.FindUserList(context.Background(), &start, &end, true, 2, "", 1, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(rsp))
	for k, _ := range rsp {
		fmt.Println(rsp[k])
	}
	var cnt int64
	cnt, err = tpf.FindUserListCount(context.Background(), &start, &end, true, 2, "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cnt)
}
