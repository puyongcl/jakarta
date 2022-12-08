package listener

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/listener"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GetUserTopRelationListenerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserTopRelationListenerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewGetUserTopRelationListenerLogic(r.Context(), svcCtx)
		resp, err := l.GetUserTopRelationListener(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
