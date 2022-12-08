package contract1021

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/admin/api/internal/logic/contract1021"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GetContract1021ByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryContract1021ByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := contract1021.NewGetContract1021ByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetContract1021ById(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
