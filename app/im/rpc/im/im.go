// Code generated by goctl. DO NOT EDIT!
// Source: im.proto

package im

import (
	"context"

	"jakarta/app/im/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AddImMsgLogReq                  = pb.AddImMsgLogReq
	AddImMsgLogResp                 = pb.AddImMsgLogResp
	GenWxMpUrlReq                   = pb.GenWxMpUrlReq
	GenWxMpUrlResp                  = pb.GenWxMpUrlResp
	GetUserUnionIdByFwhOpenIdReq    = pb.GetUserUnionIdByFwhOpenIdReq
	GetUserUnionIdByFwhOpenIdResp   = pb.GetUserUnionIdByFwhOpenIdResp
	SendDefineMsgReq                = pb.SendDefineMsgReq
	SendDefineMsgResp               = pb.SendDefineMsgResp
	SendFwhTemplateMsgReq           = pb.SendFwhTemplateMsgReq
	SendFwhTemplateMsgResp          = pb.SendFwhTemplateMsgResp
	SendMiniProgramSubscribeMsgReq  = pb.SendMiniProgramSubscribeMsgReq
	SendMiniProgramSubscribeMsgResp = pb.SendMiniProgramSubscribeMsgResp
	SendTextMsgReq                  = pb.SendTextMsgReq
	SendTextMsgResp                 = pb.SendTextMsgResp

	Im interface {
		// 发送普通文字消息
		SendTextMsg(ctx context.Context, in *SendTextMsgReq, opts ...grpc.CallOption) (*SendTextMsgResp, error)
		// 发送自定义消息
		SendDefineMsg(ctx context.Context, in *SendDefineMsgReq, opts ...grpc.CallOption) (*SendDefineMsgResp, error)
		// 发送小程序订阅消息
		SendMiniProgramSubscribeMsg(ctx context.Context, in *SendMiniProgramSubscribeMsgReq, opts ...grpc.CallOption) (*SendMiniProgramSubscribeMsgResp, error)
		// 发送微信服务号消息
		SendFwhTemplateMsg(ctx context.Context, in *SendFwhTemplateMsgReq, opts ...grpc.CallOption) (*SendFwhTemplateMsgResp, error)
		// 根据用户的服务号openid获取用户的unionId
		GetUserUnionIdByFwhOpenId(ctx context.Context, in *GetUserUnionIdByFwhOpenIdReq, opts ...grpc.CallOption) (*GetUserUnionIdByFwhOpenIdResp, error)
		// 生成小程序url link/schema
		GenWxMpUrl(ctx context.Context, in *GenWxMpUrlReq, opts ...grpc.CallOption) (*GenWxMpUrlResp, error)
		// im msg log
		AddImMsgLog(ctx context.Context, in *AddImMsgLogReq, opts ...grpc.CallOption) (*AddImMsgLogResp, error)
	}

	defaultIm struct {
		cli zrpc.Client
	}
)

func NewIm(cli zrpc.Client) Im {
	return &defaultIm{
		cli: cli,
	}
}

// 发送普通文字消息
func (m *defaultIm) SendTextMsg(ctx context.Context, in *SendTextMsgReq, opts ...grpc.CallOption) (*SendTextMsgResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.SendTextMsg(ctx, in, opts...)
}

// 发送自定义消息
func (m *defaultIm) SendDefineMsg(ctx context.Context, in *SendDefineMsgReq, opts ...grpc.CallOption) (*SendDefineMsgResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.SendDefineMsg(ctx, in, opts...)
}

// 发送小程序订阅消息
func (m *defaultIm) SendMiniProgramSubscribeMsg(ctx context.Context, in *SendMiniProgramSubscribeMsgReq, opts ...grpc.CallOption) (*SendMiniProgramSubscribeMsgResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.SendMiniProgramSubscribeMsg(ctx, in, opts...)
}

// 发送微信服务号消息
func (m *defaultIm) SendFwhTemplateMsg(ctx context.Context, in *SendFwhTemplateMsgReq, opts ...grpc.CallOption) (*SendFwhTemplateMsgResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.SendFwhTemplateMsg(ctx, in, opts...)
}

// 根据用户的服务号openid获取用户的unionId
func (m *defaultIm) GetUserUnionIdByFwhOpenId(ctx context.Context, in *GetUserUnionIdByFwhOpenIdReq, opts ...grpc.CallOption) (*GetUserUnionIdByFwhOpenIdResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.GetUserUnionIdByFwhOpenId(ctx, in, opts...)
}

// 生成小程序url link/schema
func (m *defaultIm) GenWxMpUrl(ctx context.Context, in *GenWxMpUrlReq, opts ...grpc.CallOption) (*GenWxMpUrlResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.GenWxMpUrl(ctx, in, opts...)
}

// im msg log
func (m *defaultIm) AddImMsgLog(ctx context.Context, in *AddImMsgLogReq, opts ...grpc.CallOption) (*AddImMsgLogResp, error) {
	client := pb.NewImClient(m.cli.Conn())
	return client.AddImMsgLog(ctx, in, opts...)
}
