package user

import (
	"jakarta/app/admin/api/internal/logic/user"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ProcessReportUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminProcessReportUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewProcessReportUserLogic(r.Context(), svcCtx)
		resp, err := l.ProcessReportUser(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
