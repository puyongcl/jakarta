package bbs

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/mobile/api/internal/logic/bbs"
	"jakarta/app/mobile/api/internal/svc"
	"jakarta/app/mobile/api/internal/types"
	"jakarta/common/httpresult"
	"net/http"
)

func AddStoryReplyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddStoryReplyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpresult.ParamErrorResult(r, w, err)
			return
		}

		l := bbs.NewAddStoryReplyLogic(r.Context(), svcCtx)
		resp, err := l.AddStoryReply(&req)
		httpresult.HttpResult(r, w, &req, resp, err)
	}
}
