package logic

import (
	"context"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zeromicro/go-zero/core/service"

	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFwhTemplateMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendFwhTemplateMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFwhTemplateMsgLogic {
	return &SendFwhTemplateMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  发送微信服务号消息
func (l *SendFwhTemplateMsgLogic) SendFwhTemplateMsg(in *pb.SendFwhTemplateMsgReq) (*pb.SendFwhTemplateMsgResp, error) {
	if l.svcCtx.Config.Mode != service.ProMode {
		return &pb.SendFwhTemplateMsgResp{}, nil
	}
	//
	msg := message.TemplateMessage{
		ToUser:     in.OpenId,
		TemplateID: in.TemplateId,
		URL:        "", // 模板跳转链接（海外帐号没有跳转能力）
		Color:      "", // 模板内容字体颜色，不填默认为黑色
		Data: map[string]*message.TemplateDataItem{
			"first":    {Value: in.First},
			"keyword1": {Value: in.Keyword1},
			"keyword2": {Value: in.Keyword2},
			"keyword3": {Value: in.Keyword3},
			"keyword4": {Value: in.Keyword4},
			"remark":   {Value: in.Remark, Color: in.Color},
		},
		MiniProgram: struct {
			AppID    string `json:"appid"`
			PagePath string `json:"pagepath"`
		}{AppID: l.svcCtx.Config.WxMiniConf.AppId,
			PagePath: in.Path},
	}
	_, err := l.svcCtx.Wxfwh.GetTemplate().Send(&msg)
	if err != nil {
		return nil, err
	}

	return &pb.SendFwhTemplateMsgResp{}, nil
}
