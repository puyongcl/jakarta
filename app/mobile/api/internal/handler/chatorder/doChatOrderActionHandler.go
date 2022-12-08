package chatorder

import (
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/chatorder"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
)

func DoChatOrderActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DoChatOrderActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := chatorder.NewDoChatOrderActionLogic(r.Context(), svcCtx)
		resp, err := l.DoChatOrderAction(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
