package logic

import (
	"context"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/zeromicro/go-zero/core/service"
	"jakarta/common/tool"

	"jakarta/app/im/rpc/internal/svc"
	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMiniProgramSubscribeMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMiniProgramSubscribeMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMiniProgramSubscribeMsgLogic {
	return &SendMiniProgramSubscribeMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  发送小程序订阅消息
func (l *SendMiniProgramSubscribeMsgLogic) SendMiniProgramSubscribeMsg(in *pb.SendMiniProgramSubscribeMsgReq) (*pb.SendMiniProgramSubscribeMsgResp, error) {
	if l.svcCtx.Config.Mode != service.ProMode {
		return &pb.SendMiniProgramSubscribeMsgResp{}, nil
	}

	msg := subscribe.Message{
		ToUser:     in.OpenId,
		TemplateID: in.TemplateId,
		Page:       in.Page,
		Data:       make(map[string]*subscribe.DataItem, 0),
	}
	// thing 限制20个字符 超过需要做处理
	if in.Thing5 != "" {
		text := tool.CutText(in.Thing5, 20, "...")
		msg.Data["thing5"] = &subscribe.DataItem{
			Value: text,
		}
	}
	if in.Thing4 != "" {
		text := tool.CutText(in.Thing4, 20, "...")
		msg.Data["thing4"] = &subscribe.DataItem{
			Value: text,
		}
	}
	if in.Time3 != "" {
		msg.Data["time3"] = &subscribe.DataItem{Value: in.Time3}
	}
	if in.Time2 != "" {
		msg.Data["time2"] = &subscribe.DataItem{Value: in.Time2}
	}
	if in.Thing1 != "" {
		text := tool.CutText(in.Thing1, 20, "...")
		msg.Data["thing1"] = &subscribe.DataItem{Value: text}
	}
	if in.Thing3 != "" {
		text := tool.CutText(in.Thing3, 20, "...")
		msg.Data["thing3"] = &subscribe.DataItem{Value: text}
	}
	if in.Date4 != "" {
		msg.Data["date4"] = &subscribe.DataItem{Value: in.Date4}
	}

	if l.svcCtx.Config.Mode != service.ProMode {
		msg.MiniprogramState = "developer"
	} else {
		msg.MiniprogramState = "formal"
	}

	err := l.svcCtx.MiniProgram.GetSubscribe().Send(&msg)
	if err != nil {
		return nil, err
	}

	return &pb.SendMiniProgramSubscribeMsgResp{}, nil
}
