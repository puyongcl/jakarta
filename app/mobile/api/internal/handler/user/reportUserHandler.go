package user

import (
	"jakarta/app/mobile/api/internal/logic/user"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ReportUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReportUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewReportUserLogic(r.Context(), svcCtx)
		resp, err := l.ReportUser(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
