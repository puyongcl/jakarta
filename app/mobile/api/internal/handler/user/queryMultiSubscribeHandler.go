package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/user"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func QueryMultiSubscribeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryMultiSubscribeNotifyMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewQueryMultiSubscribeLogic(r.Context(), svcCtx)
		resp, err := l.QueryMultiSubscribe(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
