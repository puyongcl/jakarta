package bbs

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/bbs"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func GetStoryReplyListByOwnHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetStoryReplyListByOwnReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := bbs.NewGetStoryReplyListByOwnLogic(r.Context(), svcCtx)
		resp, err := l.GetStoryReplyListByOwn(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
