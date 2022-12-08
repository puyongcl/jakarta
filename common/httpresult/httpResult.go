package httpresult

import (
	"fmt"
	"net/http"

	"jakarta/common/xerr"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

//http返回
func HttpResult(r *http.Request, w http.ResponseWriter, req interface{}, resp interface{}, err error) {
	if err == nil {
		//成功返回
		logx.WithContext(r.Context()).Infof("【API-OK】 req_path:%s req:%+v", r.URL.Path, req)
		logx.WithContext(r.Context()).Infof("【API-OK】 resp_path:%s resp:%+v", r.URL.Path, resp)
		rsp := Success(resp)
		httpx.WriteJson(w, http.StatusOK, rsp)
		return
	}
	logx.WithContext(r.Context()).Errorf("【API-ERR】 req_path:%s err:%+v req:%+v", r.URL.Path, err, req)
	//错误返回
	errcode := xerr.ServerCommonError
	errmsg := "服务器开小差啦，稍后再来试一试"

	gstatus, ok := status.FromError(err)
	if ok { // grpc err错误
		//区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
		errcode = uint32(gstatus.Code())
		errmsg = xerr.MapErrMsg(errcode)
		if errmsg == "" {
			if errcode > 1000 {
				errmsg = gstatus.Message()
			} else {
				errcode = xerr.ServerCommonError
				errmsg = "服务器开小差啦，稍后再来试一试"
			}
		}
	}

	logx.WithContext(r.Context()).Errorf("【API-ERR】 path:%s err code:%d msg:%s", r.URL.Path, errcode, errmsg)
	httpx.WriteJson(w, http.StatusBadRequest, Error(errcode, errmsg))
	return
}

//授权的http方法
func AuthHttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {

	if err == nil {
		//成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		errcode := xerr.ServerCommonError
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}

		logx.WithContext(r.Context()).Errorf("【GATEWAY-ERR】 : %+v ", err)

		httpx.WriteJson(w, http.StatusUnauthorized, Error(errcode, errmsg))
	}
}

//http 参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", xerr.MapErrMsg(xerr.RequestParamError), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(xerr.RequestParamError, errMsg))
}
