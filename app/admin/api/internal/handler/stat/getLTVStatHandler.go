package stat

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/admin/api/internal/logic/stat"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GetLTVStatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLifeTimeValueStatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := stat.NewGetLTVStatLogic(r.Context(), svcCtx)
		resp, err := l.GetLTVStat(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
