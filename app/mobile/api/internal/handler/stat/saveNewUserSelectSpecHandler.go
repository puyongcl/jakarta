package stat

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/stat"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func SaveNewUserSelectSpecHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveNewUserSelectSpecReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := stat.NewSaveNewUserSelectSpecLogic(r.Context(), svcCtx)
		resp, err := l.SaveNewUserSelectSpec(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
