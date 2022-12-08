package third

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/payment/api/internal/logic/third"
	"net/http"

	"jakarta/app/payment/api/internal/svc"
)

func WxRefundCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := third.NewWxRefundCallbackLogic(r.Context(), svcCtx)
		resp, err := l.WxRefundCallback(w, r)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("【API-ERR】WxRefundCallbackHandler: %+v ", err)
			httpx.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		httpx.OkJson(w, resp)
	}
}
