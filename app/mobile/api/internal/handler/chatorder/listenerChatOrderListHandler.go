package chatorder

import (
	"jakarta/app/mobile/api/internal/logic/chatorder"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListenerChatOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetListenerSeeChatOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := chatorder.NewListenerChatOrderListLogic(r.Context(), svcCtx)
		resp, err := l.ListenerChatOrderList(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
