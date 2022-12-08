package adminlog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	adminPgModel2 "jakarta/app/pgModel/adminPgModel"
	"jakarta/common/xerr"
)

func SaveAdminLog(ctx context.Context, pgM adminPgModel2.AdminLogModel, routePath string, adminUid int64, errIn error, req, resp interface{}) {
	sreq, err := json.Marshal(req)
	if err != nil {
		logx.WithContext(ctx).Errorf("saveAdminLog req:%+v json marshal err:%+v", req, err)
		return
	}
	var sresp string
	if errIn != nil {
		resp = xerr.ErrMsg{
			Msg: fmt.Sprintf("%+v", errIn),
		}
	}
	var jresp []byte
	jresp, err = json.Marshal(resp)
	if err != nil {
		logx.WithContext(ctx).Errorf("saveAdminLog resp:%+v json marshal err:%+v", resp, err)
		return
	}
	sresp = string(jresp)

	data := adminPgModel2.AdminLog{
		AdminUid:  adminUid,
		Request:   string(sreq),
		Response:  sresp,
		RoutePath: routePath,
	}

	_, err = pgM.Insert(ctx, &data)
	if err != nil {
		logx.WithContext(ctx).Errorf("saveAdminLog data:%+v Insert err:%+v", &data, err)
		return
	}
	return
}
