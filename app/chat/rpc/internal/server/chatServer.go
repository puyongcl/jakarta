// Code generated by goctl. DO NOT EDIT!
// Source: chat.proto

package server

import (
	"context"

	"jakarta/app/chat/rpc/internal/logic"
	"jakarta/app/chat/rpc/internal/svc"
	"jakarta/app/chat/rpc/pb"
)

type ChatServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedChatServer
}

func NewChatServer(svcCtx *svc.ServiceContext) *ChatServer {
	return &ChatServer{
		svcCtx: svcCtx,
	}
}

//  同步聊天状态过程
func (s *ChatServer) SyncChatState(ctx context.Context, in *pb.SyncChatStateReq) (*pb.SyncChatStateResp, error) {
	l := logic.NewSyncChatStateLogic(ctx, s.svcCtx)
	return l.SyncChatState(in)
}

//  更新用户聊天可用时间
func (s *ChatServer) UpdateUserChatBalance(ctx context.Context, in *pb.UpdateUserChatBalanceReq) (*pb.UpdateUserChatBalanceResp, error) {
	l := logic.NewUpdateUserChatBalanceLogic(ctx, s.svcCtx)
	return l.UpdateUserChatBalance(in)
}

//  结算当前通话记录
func (s *ChatServer) UpdateVoiceChatStat(ctx context.Context, in *pb.UpdateVoiceChatStatReq) (*pb.UpdateVoiceChatStatResp, error) {
	l := logic.NewUpdateVoiceChatStatLogic(ctx, s.svcCtx)
	return l.UpdateVoiceChatStat(in)
}

//  获取时间快结束的语音通话
func (s *ChatServer) GetUseOutVoiceChat(ctx context.Context, in *pb.GetUseOutVoiceChatReq) (*pb.GetUseOutVoiceChatResp, error) {
	l := logic.NewGetUseOutVoiceChatLogic(ctx, s.svcCtx)
	return l.GetUseOutVoiceChat(in)
}

//  获取时间快结束的文字通话
func (s *ChatServer) GetUseOutTextChat(ctx context.Context, in *pb.GetUseOutTextChatReq) (*pb.GetUseOutTextChatResp, error) {
	l := logic.NewGetUseOutTextChatLogic(ctx, s.svcCtx)
	return l.GetUseOutTextChat(in)
}

//  重置免费聊天次数
func (s *ChatServer) ResetFreeTextChatCnt(ctx context.Context, in *pb.ResetFreeTextChatCntReq) (*pb.ResetFreeTextChatCntResp, error) {
	l := logic.NewResetFreeTextChatCntLogic(ctx, s.svcCtx)
	return l.ResetFreeTextChatCnt(in)
}

//  更新统计进入XXX聊天页面用户数
func (s *ChatServer) UpdateTodayEnterChatUserCnt(ctx context.Context, in *pb.UpdateTodayEnterChatUserCntReq) (*pb.UpdateTodayEnterChatUserCntResp, error) {
	l := logic.NewUpdateTodayEnterChatUserCntLogic(ctx, s.svcCtx)
	return l.UpdateTodayEnterChatUserCnt(in)
}

//  更新统计近几天进入XXX页面用户数
func (s *ChatServer) UpdateLastDaysEnterChatUserCnt(ctx context.Context, in *pb.UpdateLastDaysEnterChatUserCntReq) (*pb.UpdateLastDaysEnterChatUserCntResp, error) {
	l := logic.NewUpdateLastDaysEnterChatUserCntLogic(ctx, s.svcCtx)
	return l.UpdateLastDaysEnterChatUserCnt(in)
}

//  更新文字聊天时间用完
func (s *ChatServer) UpdateTextChatOver(ctx context.Context, in *pb.UpdateTextChatOverReq) (*pb.UpdateTextChatOverResp, error) {
	l := logic.NewUpdateTextChatOverLogic(ctx, s.svcCtx)
	return l.UpdateTextChatOver(in)
}

//  初始化XXX通话状态
func (s *ChatServer) CreateListenerVoiceChatState(ctx context.Context, in *pb.CreateListenerVoiceChatStateReq) (*pb.CreateListenerVoiceChatStateResp, error) {
	l := logic.NewCreateListenerVoiceChatStateLogic(ctx, s.svcCtx)
	return l.CreateListenerVoiceChatState(in)
}

//  用户和XXX交互事件
func (s *ChatServer) SendUserListenerRelationEvent(ctx context.Context, in *pb.SendUserListenerRelationEventReq) (*pb.SendUserListenerRelationEventResp, error) {
	l := logic.NewSendUserListenerRelationEventLogic(ctx, s.svcCtx)
	return l.SendUserListenerRelationEvent(in)
}

//  获取交互最频繁的几位XXX
func (s *ChatServer) GetTopUserAndListenerRelation(ctx context.Context, in *pb.GetTopUserAndListenerRelationReq) (*pb.GetTopUserAndListenerRelationResp, error) {
	l := logic.NewGetTopUserAndListenerRelationLogic(ctx, s.svcCtx)
	return l.GetTopUserAndListenerRelation(in)
}
