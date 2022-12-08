package listenerPgModel

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"jakarta/common/key/db"
	"reflect"
	"testing"
	"time"
)

func TestPg(t *testing.T) {
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
	data := ListenerProfileDraft{
		ListenerUid:   1000,
		NickName:      "test",
		Avatar:        "test1",
		SmallAvatar:   "test1",
		MaritalStatus: 1,
		PhoneNumber:   "128193",
		Constellation: 2,
		Province:      "121",
		City:          "1212",
		Job:           "212",
		Education:     2,
		Gender:        3,
		Birthday: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		IdNo:                "1212",
		IdPhoto1:            "2112",
		IdPhoto2:            "1212",
		IdPhoto3:            "21",
		Specialties:         []int64{109, 208},
		Introduction:        "21",
		VoiceFile:           "212",
		Experience1:         "12",
		Experience2:         "21",
		CertType:            0,
		CertFiles1:          "212",
		CertFiles2:          "21",
		CertFiles3:          "21",
		CertFiles4:          "21",
		CertFiles5:          "21",
		AutoReplyNew:        "21",
		AutoReplyProcessing: "21",
		AutoReplyFinish:     "21",
		TextChatPrice:       110,
		VoiceChatPrice:      110,
		CheckFailField:      []string{"AutoReplyProcessing", "CertFiles3"},
		CheckStatus:         2,
	}
	rs, err := tpf.Insert(context.Background(), &data)
	fmt.Println(rs)
	fmt.Println(err)
	return
}

func TestPgQuery(t *testing.T) {
	ds := "postgres://jakarta:postgres@api.domain.com:5432/jakarta_listener?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewListenerProfileDraftModel(sqlConn, c)

	resp, err := tpf.FindOne(context.Background(), 100002)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	resp.NickName = "TestName2"
	err = tpf.UpdateTrans(context.Background(), nil, resp)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestUpdate(t *testing.T) {
	ds := "postgres://jakarta:postgres@api.domain.com:5432/jakarta_listener?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "192.168.1.12:36379",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpf := NewListenerProfileModel(sqlConn, c)

	err := tpf.UpdateListenerStat(context.Background(), &AddListenerStat{
		ListenerUid:       100003,
		AddUserCount:      0,
		AddChatDuration:   1,
		AddRatingSum:      0,
		AddFiveStar:       0,
		AddRefundOrderCnt: 0,
		AddFinishOrderCnt: 0,
		AddPaidOrderCnt:   0,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTrans(t *testing.T) {
	ds := "postgres://jakarta:postgres@api.domain.com:5432/jakarta_listener?sslmode=disable"
	c := []cache.NodeConf{
		{RedisConf: redis.RedisConf{
			Host: "api.domain.com:36381",
			Pass: "G62m50oigInC30sf",
		},
			Weight: 100,
		},
	}
	sqlConn := sqlx.NewSqlConn(db.PostgresDriverName, ds)
	tpd := newListenerProfileDraftModel(sqlConn, c)
	tp := NewListenerProfileModel(sqlConn, c)
	//tl := newListenerWalletModel(sqlConn, c)
	err := tpd.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		//_, err2 := tl.InsertOrUpdateMPTrans(ctx, session, &ListenerWallet{
		//	ListenerUid: 199,
		//})
		//if err2 != nil {
		//	return err2
		//}
		_, err2 := tp.InsertTrans(ctx, session, &ListenerProfile{
			ListenerUid: 199,
			Specialties: pq.Int64Array{1, 2},
		})
		if err2 != nil {
			return err2
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestCopy2(t *testing.T) {
	type birthday struct {
		Uid      string `json:"uid"`
		Birthday string `json:"birthday"`
	}
	type NullTime struct {
		Time  time.Time
		Valid bool // Valid is true if Time is not NULL
	}
	type birthday2 struct {
		Uid      string   `json:"uid"`
		Birthday NullTime `json:"birthday"`
	}
	in := &birthday{
		Uid:      "221",
		Birthday: "",
	}
	draftData := new(birthday2)
	_ = copier.Copy(draftData, in)
	fmt.Println(draftData)
}

func TestReflect(t *testing.T) {
	// 获取平均值
	var data ListenerStatAverage
	t1 := reflect.TypeOf(ListenerDashboardStat{})
	t2 := reflect.TypeOf(data)
	tv2 := reflect.ValueOf(&data)
	mv := make(map[string]int64, 0)
	var cnt int64
	var dbFieldName, fieldName string
	for k := 0; k < t1.NumField(); k++ {
		if t1.Field(k).Type.Kind() == reflect.Int64 {
			fieldName = t1.Field(k).Name
			if fieldName == "ListenerUid" {
				continue
			}
			dbFieldName = t1.Field(k).Tag.Get("db")
			fmt.Println(dbFieldName)
			cnt++
			mv[fieldName] = cnt
		}
	}

	for k := 0; k < t2.NumField(); k++ {
		if t2.Field(k).Type.Kind() == reflect.Int64 {
			fieldName = t2.Field(k).Name
			v, ok := mv[fieldName]
			if ok {
				tv2.Elem().FieldByName(fieldName).SetInt(v)
			}
		}
	}
	fmt.Println(&data)
}
