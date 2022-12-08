package third

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io"
	"jakarta/app/payment/api/internal/logic/third"
	"jakarta/app/payment/api/internal/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"jakarta/app/payment/api/internal/svc"
)

func HfbfCashCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HFBFCashCallbackReq
		req.Params = &types.HFBFCashCallbackData{}

		var err error
		defer func() {
			if err == nil {
				return
			}
			httpx.WriteJson(w, http.StatusBadRequest, &types.HFBFCashCallbackResp{
				Code:    400,
				Message: "fail",
			})
		}()

		//
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logx.WithContext(r.Context()).Errorf("HfbfCashCallbackHandler ReadAll err:%+v query:%s", err, r.URL.RawQuery)
			return
		}

		ul, err := url.QueryUnescape(string(body))
		if err != nil {
			logx.WithContext(r.Context()).Errorf("HfbfCashCallbackHandler QueryUnescape err:%+v query:%s", err, r.URL.RawQuery)
			return
		}
		ul = strings.Replace(ul, "\\u0026", "&", -1)
		arr := strings.Split(ul, "&")
		ma := make(map[string]string)
		for idx := 0; idx < len(arr); idx++ {
			as := strings.Split(arr[idx], "=")
			if len(as) == 2 {
				k := strings.Replace(as[0], "params[", "", -1)
				k = strings.Replace(k, "]", "", -1)
				ma[k] = as[1]
			}
		}

		ts, ok := ma["timeStamp"]
		if ok {
			req.TimeStamp, err = strconv.ParseInt(ts, 10, 64)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("HfbfCashCallbackHandler ParseInt: %+v ts:%s", err, ts)
			}
		}

		ts, ok = ma["type"]
		if ok {
			req.Params.Type, err = strconv.ParseInt(ts, 10, 64)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("HfbfCashCallbackHandler ParseInt: %+v ts:%s", err, ts)
			}
		}

		req.Params.WorkNumber, ok = ma["work_number"]
		ts, ok = ma["company_id"]
		if ok {
			req.Params.CompanyId, err = strconv.ParseInt(ts, 10, 64)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("HfbfCashCallbackHandler ParseInt: %+v ts:%s", err, ts)
			}
		}

		req.Params.UserId = ma["user_id"]
		req.Params.Number = ma["number"]
		req.Params.PayStatus = ma["pay_status"]
		req.Params.CustomNumber = ma["custom_number"]
		req.Params.Msg = ma["msg"]
		req.Params.PayTime = ma["pay_time"]
		req.Sign = ma["sign"]

		l := third.NewHfbfCashCallbackLogic(r.Context(), svcCtx)
		resp, err := l.HfbfCashCallback(&req)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("HfbfCashCallbackHandler: %+v ", err)
			return
		}

		httpx.OkJson(w, resp)
	}
}
