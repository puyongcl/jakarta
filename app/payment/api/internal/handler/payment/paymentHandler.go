package payment

import (
	"jakarta/common/httpresult"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/payment/api/internal/logic/payment"
	"jakarta/app/payment/api/internal/svc"
	"jakarta/app/payment/api/internal/types"
)

func PaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := payment.NewPaymentLogic(r.Context(), svcCtx)
		resp, err := l.Payment(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
