package mp

import (
	"context"
	"github.com/silenceper/wechat/v2/util"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"jakarta/app/im/api/internal/svc"
	"net/http"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type WxMpCallbackReq struct {
	ToUserName string `json:"ToUserName"`
	Encrypt    string `json:"Encrypt"`
}

func (l *CallbackLogic) WxMpCallback(w http.ResponseWriter, r *http.Request) (err error) {
	var req WxMpCallbackReq
	err = httpx.Parse(r, &req)
	if err != nil {
		return
	}

	//
	var jsonStr []byte
	_, jsonStr, err = util.DecryptMsg(l.svcCtx.Config.WxMiniConf.AppId, req.Encrypt, l.svcCtx.Config.WxMiniConf.EncodingAESKey)
	if err != nil {
		return err
	}
	logx.WithContext(l.ctx).Infof("WxMpCallback:%s", jsonStr) //在返回页面中显示内容。
	return nil
}
