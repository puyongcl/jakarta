package handler

import (
	"jakarta/app/im/api/internal/handler/fwh"
	"jakarta/app/im/api/internal/handler/mp"
	"net/http"

	"jakarta/app/im/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers2(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/wx/fwh",
				Handler: fwh.CallbackVerifyHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/wx/fwh",
				Handler: fwh.CallbackHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/wx/mp",
				Handler: mp.CallbackVerifyHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/wx/mp",
				Handler: mp.CallbackHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/wx/mp/url/baidu/gen",
				Handler: mp.GenBaiduMpUrlHandler(serverCtx),
			},
		},
		rest.WithPrefix("/im/v1"),
	)
}
