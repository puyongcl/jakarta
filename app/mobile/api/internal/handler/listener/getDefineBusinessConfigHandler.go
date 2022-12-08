package listener

import (
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/listener"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
)

func GetDefineBusinessConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDefineBusinessConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewGetDefineBusinessConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetBannerAndFilterConfigReq(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
