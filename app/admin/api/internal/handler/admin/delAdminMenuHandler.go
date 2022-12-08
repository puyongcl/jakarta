package admin

import (
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/admin/api/internal/logic/admin"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
)

func DelAdminMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelAdminMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := admin.NewDelAdminMenuLogic(r.Context(), svcCtx)
		resp, err := l.DelAdminMenu(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
