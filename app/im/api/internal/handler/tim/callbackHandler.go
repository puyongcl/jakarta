package tim

import (
	"github.com/zeromicro/go-zero/core/logx"
	"jakarta/common/key/timkey"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/im/api/internal/logic/tim"
	"jakarta/app/im/api/internal/svc"
	"jakarta/app/im/api/internal/types"
)

func CallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			resp := &types.TIMCallbackResp{
				ActionStatus: "OK",
				ErrorCode:    0,
				ErrorInfo:    "",
			}
			httpx.OkJson(w, resp)
		}()

		var err error
		callbackCommand := r.URL.Query().Get("CallbackCommand")

		//logx.WithContext(r.Context()).Infof("Url val:%+v", urlArg)
		// 根据command处理不同请求
		switch callbackCommand {
		case timkey.TimCommandStateChange: // 1、用户状态变更
			var req types.TIMCallbackStateChangeReq
			err = httpx.Parse(r, &req)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("Parse bod ERR:%+v ", err)
				return
			}

			l := tim.NewCallbackLogic(r.Context(), svcCtx)
			_, err = l.StateChangeCallback(&req)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("StateChangeCallback req:%+v ERR:%+v ", req, err)
			}

		case timkey.TimCommandAfterSendMsg: // 发送单聊消息之后
			var req types.TIMCallbackAfterSendMsgReq
			err = httpx.Parse(r, &req)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("Parse bod ERR:%+v ", err)
				return
			}

			l := tim.NewCallbackLogic(r.Context(), svcCtx)
			_, err = l.AfterSendMsgCallback(&req)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("AfterSendMsgCallback req:%+v ERR:%+v ", req, err)
			}

		default:

		}
		return
	}
}
