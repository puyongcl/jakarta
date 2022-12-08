package order

import (
	"jakarta/app/admin/api/internal/logic/order"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckRefundOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckRefundOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := order.NewCheckRefundOrderLogic(r.Context(), svcCtx)
		resp, err := l.CheckRefundOrder(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
