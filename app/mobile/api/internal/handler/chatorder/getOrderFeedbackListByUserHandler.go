package chatorder

import (
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/chatorder"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
)

func GetOrderFeedbackListByUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetChatOrderFeedbackListByUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := chatorder.NewGetOrderFeedbackListByUserLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderFeedbackListByUser(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
