package user

import (
	"jakarta/app/mobile/api/internal/logic/user"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func BlockUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BlockUserReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewBlockUserLogic(r.Context(), svcCtx)
		resp, err := l.BlockUser(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
