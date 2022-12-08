package order

import (
	"jakarta/app/admin/api/internal/logic/order"
	"jakarta/app/admin/api/internal/svc"
	"jakarta/app/admin/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRefundOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRefundOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := order.NewGetRefundOrderListLogic(r.Context(), svcCtx)
		resp, err := l.GetRefundOrderList(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
