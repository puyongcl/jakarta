package mp

import (
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/app/im/api/internal/logic/mp"
	"net/http"

	"jakarta/app/im/api/internal/svc"
)

// 接入验证
func CallbackVerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		echo := r.URL.Query().Get("echostr")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(echo))
		if err != nil {
			logx.WithContext(r.Context()).Errorf("CallbackVerifyHandler err:%+v", err)
			return
		}
		return
	}
}

// 事件、消息推送
func CallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ca := mp.NewCallbackLogic(r.Context(), svcCtx)
		err := ca.WxMpCallback(w, r)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("WxMpCallback err:%+v", err)
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("success"))
		if err != nil {
			logx.WithContext(r.Context()).Errorf("WxMpCallback err:%+v", err)
			return
		}
		return
	}
}
