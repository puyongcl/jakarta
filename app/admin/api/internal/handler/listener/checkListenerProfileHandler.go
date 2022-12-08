package listener

import (
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/admin/api/internal/logic/listener"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
)

func CheckListenerProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckListenerProfileReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewCheckListenerProfileLogic(r.Context(), svcCtx)
		resp, err := l.CheckListenerProfile(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
