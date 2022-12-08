package stat

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/admin/api/internal/logic/stat"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GetDailyStatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDailyStatListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := stat.NewGetDailyStatLogic(r.Context(), svcCtx)
		resp, err := l.GetDailyStat(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
