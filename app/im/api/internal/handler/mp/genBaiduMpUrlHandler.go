package mp

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/im/api/internal/logic/mp"
	"jakarta/app/im/api/internal/svc"
	"jakarta/app/im/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GenBaiduMpUrlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenBaiduWxMpUrlReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := mp.NewGenBaiduMpUrlLogic(r.Context(), svcCtx)
		resp := l.GenBaiduMpUrl(&req)
		httpx.WriteJson(w, http.StatusOK, resp)
	}
}
