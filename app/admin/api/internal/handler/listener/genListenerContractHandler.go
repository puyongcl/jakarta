package listener

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/admin/api/internal/logic/listener"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GenListenerContractHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenListenerContractReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewGenListenerContractLogic(r.Context(), svcCtx)
		resp, err := l.GenListenerContract(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
