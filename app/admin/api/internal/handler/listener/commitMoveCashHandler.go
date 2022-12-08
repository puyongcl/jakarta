package listener

import (
	"jakarta/app/admin/api/internal/logic/listener"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CommitMoveCashHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommitMoveCashReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := listener.NewCommitMoveCashLogic(r.Context(), svcCtx)
		resp, err := l.CommitMoveCash(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
