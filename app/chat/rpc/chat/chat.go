// Code generated by goctl. DO NOT EDIT!
// Source: chat.proto

package chat

import (
	"context"

	"jakarta/app/chat/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateListenerVoiceChatStateReq    = pb.CreateListenerVoiceChatStateReq
	CreateListenerVoiceChatStateResp   = pb.CreateListenerVoiceChatStateResp
	GetTopUserAndListenerRelationReq   = pb.GetTopUserAndListenerRelationReq
	GetTopUserAndListenerRelationResp  = pb.GetTopUserAndListenerRelationResp
	GetUseOutTextChatReq               = pb.GetUseOutTextChatReq
	GetUseOutTextChatResp              = pb.GetUseOutTextChatResp
	GetUseOutVoiceChatReq              = pb.GetUseOutVoiceChatReq
	GetUseOutVoiceChatResp             = pb.GetUseOutVoiceChatResp
	ResetFreeTextChatCntReq            = pb.ResetFreeTextChatCntReq
	ResetFreeTextChatCntResp           = pb.ResetFreeTextChatCntResp
	SendUserListenerRelationEventReq   = pb.SendUserListenerRelationEventReq
	SendUserListenerRelationEventResp  = pb.SendUserListenerRelationEventResp
	SyncChatStateReq                   = pb.SyncChatStateReq
	SyncChatStateResp                  = pb.SyncChatStateResp
	SyncListenerFreeChatCntReq         = pb.SyncListenerFreeChatCntReq
	SyncListenerFreeChatCntResp        = pb.SyncListenerFreeChatCntResp
	TextChatUser                       = pb.TextChatUser
	UpdateLastDaysEnterChatUserCntReq  = pb.UpdateLastDaysEnterChatUserCntReq
	UpdateLastDaysEnterChatUserCntResp = pb.UpdateLastDaysEnterChatUserCntResp
	UpdateTextChatOverReq              = pb.UpdateTextChatOverReq
	UpdateTextChatOverResp             = pb.UpdateTextChatOverResp
	UpdateTodayEnterChatUserCntReq     = pb.UpdateTodayEnterChatUserCntReq
	UpdateTodayEnterChatUserCntResp    = pb.UpdateTodayEnterChatUserCntResp
	UpdateUserChatBalanceReq           = pb.UpdateUserChatBalanceReq
	UpdateUserChatBalanceResp          = pb.UpdateUserChatBalanceResp
	UpdateVoiceChatStatReq             = pb.UpdateVoiceChatStatReq
	UpdateVoiceChatStatResp            = pb.UpdateVoiceChatStatResp
	UserAndListenerRelation            = pb.UserAndListenerRelation
	VoiceChatUser                      = pb.VoiceChatUser

	Chat interface {
		//  同步聊天状态过程
		SyncChatState(ctx context.Context, in *SyncChatStateReq, opts ...grpc.CallOption) (*SyncChatStateResp, error)
		//  更新用户聊天可用时间
		UpdateUserChatBalance(ctx context.Context, in *UpdateUserChatBalanceReq, opts ...grpc.CallOption) (*UpdateUserChatBalanceResp, error)
		//  结算当前通话记录
		UpdateVoiceChatStat(ctx context.Context, in *UpdateVoiceChatStatReq, opts ...grpc.CallOption) (*UpdateVoiceChatStatResp, error)
		//  获取时间快结束的语音通话
		GetUseOutVoiceChat(ctx context.Context, in *GetUseOutVoiceChatReq, opts ...grpc.CallOption) (*GetUseOutVoiceChatResp, error)
		//  获取时间快结束的文字通话
		GetUseOutTextChat(ctx context.Context, in *GetUseOutTextChatReq, opts ...grpc.CallOption) (*GetUseOutTextChatResp, error)
		//  重置免费聊天次数
		ResetFreeTextChatCnt(ctx context.Context, in *ResetFreeTextChatCntReq, opts ...grpc.CallOption) (*ResetFreeTextChatCntResp, error)
		//  更新统计进入XXX聊天页面用户数
		UpdateTodayEnterChatUserCnt(ctx context.Context, in *UpdateTodayEnterChatUserCntReq, opts ...grpc.CallOption) (*UpdateTodayEnterChatUserCntResp, error)
		//  更新统计近几天进入XXX页面用户数
		UpdateLastDaysEnterChatUserCnt(ctx context.Context, in *UpdateLastDaysEnterChatUserCntReq, opts ...grpc.CallOption) (*UpdateLastDaysEnterChatUserCntResp, error)
		//  更新文字聊天时间用完
		UpdateTextChatOver(ctx context.Context, in *UpdateTextChatOverReq, opts ...grpc.CallOption) (*UpdateTextChatOverResp, error)
		//  初始化XXX通话状态
		CreateListenerVoiceChatState(ctx context.Context, in *CreateListenerVoiceChatStateReq, opts ...grpc.CallOption) (*CreateListenerVoiceChatStateResp, error)
		//  用户和XXX交互事件
		SendUserListenerRelationEvent(ctx context.Context, in *SendUserListenerRelationEventReq, opts ...grpc.CallOption) (*SendUserListenerRelationEventResp, error)
		//  获取交互最频繁的几位XXX
		GetTopUserAndListenerRelation(ctx context.Context, in *GetTopUserAndListenerRelationReq, opts ...grpc.CallOption) (*GetTopUserAndListenerRelationResp, error)
	}

	defaultChat struct {
		cli zrpc.Client
	}
)

func NewChat(cli zrpc.Client) Chat {
	return &defaultChat{
		cli: cli,
	}
}

//  同步聊天状态过程
func (m *defaultChat) SyncChatState(ctx context.Context, in *SyncChatStateReq, opts ...grpc.CallOption) (*SyncChatStateResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.SyncChatState(ctx, in, opts...)
}

//  更新用户聊天可用时间
func (m *defaultChat) UpdateUserChatBalance(ctx context.Context, in *UpdateUserChatBalanceReq, opts ...grpc.CallOption) (*UpdateUserChatBalanceResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.UpdateUserChatBalance(ctx, in, opts...)
}

//  结算当前通话记录
func (m *defaultChat) UpdateVoiceChatStat(ctx context.Context, in *UpdateVoiceChatStatReq, opts ...grpc.CallOption) (*UpdateVoiceChatStatResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.UpdateVoiceChatStat(ctx, in, opts...)
}

//  获取时间快结束的语音通话
func (m *defaultChat) GetUseOutVoiceChat(ctx context.Context, in *GetUseOutVoiceChatReq, opts ...grpc.CallOption) (*GetUseOutVoiceChatResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.GetUseOutVoiceChat(ctx, in, opts...)
}

//  获取时间快结束的文字通话
func (m *defaultChat) GetUseOutTextChat(ctx context.Context, in *GetUseOutTextChatReq, opts ...grpc.CallOption) (*GetUseOutTextChatResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.GetUseOutTextChat(ctx, in, opts...)
}

//  重置免费聊天次数
func (m *defaultChat) ResetFreeTextChatCnt(ctx context.Context, in *ResetFreeTextChatCntReq, opts ...grpc.CallOption) (*ResetFreeTextChatCntResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.ResetFreeTextChatCnt(ctx, in, opts...)
}

//  更新统计进入XXX聊天页面用户数
func (m *defaultChat) UpdateTodayEnterChatUserCnt(ctx context.Context, in *UpdateTodayEnterChatUserCntReq, opts ...grpc.CallOption) (*UpdateTodayEnterChatUserCntResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.UpdateTodayEnterChatUserCnt(ctx, in, opts...)
}

//  更新统计近几天进入XXX页面用户数
func (m *defaultChat) UpdateLastDaysEnterChatUserCnt(ctx context.Context, in *UpdateLastDaysEnterChatUserCntReq, opts ...grpc.CallOption) (*UpdateLastDaysEnterChatUserCntResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.UpdateLastDaysEnterChatUserCnt(ctx, in, opts...)
}

//  更新文字聊天时间用完
func (m *defaultChat) UpdateTextChatOver(ctx context.Context, in *UpdateTextChatOverReq, opts ...grpc.CallOption) (*UpdateTextChatOverResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.UpdateTextChatOver(ctx, in, opts...)
}

//  初始化XXX通话状态
func (m *defaultChat) CreateListenerVoiceChatState(ctx context.Context, in *CreateListenerVoiceChatStateReq, opts ...grpc.CallOption) (*CreateListenerVoiceChatStateResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.CreateListenerVoiceChatState(ctx, in, opts...)
}

//  用户和XXX交互事件
func (m *defaultChat) SendUserListenerRelationEvent(ctx context.Context, in *SendUserListenerRelationEventReq, opts ...grpc.CallOption) (*SendUserListenerRelationEventResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.SendUserListenerRelationEvent(ctx, in, opts...)
}

//  获取交互最频繁的几位XXX
func (m *defaultChat) GetTopUserAndListenerRelation(ctx context.Context, in *GetTopUserAndListenerRelationReq, opts ...grpc.CallOption) (*GetTopUserAndListenerRelationResp, error) {
	client := pb.NewChatClient(m.cli.Conn())
	return client.GetTopUserAndListenerRelation(ctx, in, opts...)
}
