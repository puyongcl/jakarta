package listener

import (
	"jakarta/app/mobile/api/internal/logic/listener"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetListenerWordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetListenerWordsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewGetListenerWordsLogic(r.Context(), svcCtx)
		resp, err := l.GetListenerWords(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
