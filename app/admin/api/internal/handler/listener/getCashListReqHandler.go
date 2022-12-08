package listener

import (
	"jakarta/app/admin/api/internal/logic/listener"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCashListReqHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCashListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewGetCashListReqLogic(r.Context(), svcCtx)
		resp, err := l.GetCashListReq(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
