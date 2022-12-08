package listener

import (
	"jakarta/app/mobile/api/internal/logic/listener"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetListenerHomePageDashboardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetListenerHomePageDashboardReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewGetListenerHomePageDashboardLogic(r.Context(), svcCtx)
		resp, err := l.GetListenerHomePageDashboard(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
