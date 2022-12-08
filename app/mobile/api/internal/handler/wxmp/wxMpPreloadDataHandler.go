package wxmp

import (
	"jakarta/app/mobile/api/internal/logic/wxmp"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
	"strconv"
)

func WxMpPreloadDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetWxMpPreloadDataReq
		req.Appid = r.URL.Query().Get("appid")
		req.Token = r.URL.Query().Get("token")
		req.Code = r.URL.Query().Get("code")
		val := r.URL.Query().Get("timestamp")
		req.Timestamp, _ = strconv.ParseInt(val, 10, 64)

		req.Path = r.URL.Query().Get("path")
		req.Query = r.URL.Query().Get("query")
		val = r.URL.Query().Get("scene")
		req.Scene, _ = strconv.ParseInt(val, 10, 64)
		l := wxmp.NewWxMpPreloadDataLogic(r.Context(), svcCtx)
		resp, err := l.WxMpPreloadData(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
