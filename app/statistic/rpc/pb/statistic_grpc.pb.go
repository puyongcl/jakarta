// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: statistic.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// StatisticClient is the client API for Statistic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StatisticClient interface {
	// 获取用户列表
	GetUserList(ctx context.Context, in *GetUserListReq, opts ...grpc.CallOption) (*GetUserListResp, error)
	// 更新每日统计数据
	UpdateStatisticDailyData(ctx context.Context, in *UpdateStatisticDailyDataReq, opts ...grpc.CallOption) (*UpdateStatisticDailyDataResp, error)
	// 获取每日统计数据
	GetDailyStatList(ctx context.Context, in *GetDailyStatListReq, opts ...grpc.CallOption) (*GetDailyStatListResp, error)
	// 更新每日登陆时间
	UpdateLoginLog(ctx context.Context, in *UpdateLoginLogReq, opts ...grpc.CallOption) (*UpdateLoginLogResp, error)
	// 获取统计近多少日的用户在昨日累计数据
	GetLifeTimeValueStat(ctx context.Context, in *GetLifeTimeValueStatReq, opts ...grpc.CallOption) (*GetLifeTimeValueStatResp, error)
	// 获取用户渠道列表
	GetUserChannelList(ctx context.Context, in *GetUserChannelListReq, opts ...grpc.CallOption) (*GetUserChannelListResp, error)
	// 新用户选择的XX标签
	SaveNewUserSelectSpec(ctx context.Context, in *SaveNewUserSelectSpecReq, opts ...grpc.CallOption) (*SaveNewUserSelectSpecResp, error)
	// 定时统计用户和XXX状态数据
	UpdateUserStateStat(ctx context.Context, in *UpdateUserStateStatReq, opts ...grpc.CallOption) (*UpdateUserStateStatResp, error)
	// 保存成人依恋量表测试结果
	SaveAdultQuizECR(ctx context.Context, in *SaveAdultQuizEcrReq, opts ...grpc.CallOption) (*SaveAdultQuizEcrResp, error)
	// 获取最新成人依恋量表测试结果
	GetAdultQuizEcr(ctx context.Context, in *GetAdultQuizEcrReq, opts ...grpc.CallOption) (*GetAdultQuizEcrResp, error)
}

type statisticClient struct {
	cc grpc.ClientConnInterface
}

func NewStatisticClient(cc grpc.ClientConnInterface) StatisticClient {
	return &statisticClient{cc}
}

func (c *statisticClient) GetUserList(ctx context.Context, in *GetUserListReq, opts ...grpc.CallOption) (*GetUserListResp, error) {
	out := new(GetUserListResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/getUserList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) UpdateStatisticDailyData(ctx context.Context, in *UpdateStatisticDailyDataReq, opts ...grpc.CallOption) (*UpdateStatisticDailyDataResp, error) {
	out := new(UpdateStatisticDailyDataResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/updateStatisticDailyData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) GetDailyStatList(ctx context.Context, in *GetDailyStatListReq, opts ...grpc.CallOption) (*GetDailyStatListResp, error) {
	out := new(GetDailyStatListResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/getDailyStatList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) UpdateLoginLog(ctx context.Context, in *UpdateLoginLogReq, opts ...grpc.CallOption) (*UpdateLoginLogResp, error) {
	out := new(UpdateLoginLogResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/updateLoginLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) GetLifeTimeValueStat(ctx context.Context, in *GetLifeTimeValueStatReq, opts ...grpc.CallOption) (*GetLifeTimeValueStatResp, error) {
	out := new(GetLifeTimeValueStatResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/getLifeTimeValueStat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) GetUserChannelList(ctx context.Context, in *GetUserChannelListReq, opts ...grpc.CallOption) (*GetUserChannelListResp, error) {
	out := new(GetUserChannelListResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/getUserChannelList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) SaveNewUserSelectSpec(ctx context.Context, in *SaveNewUserSelectSpecReq, opts ...grpc.CallOption) (*SaveNewUserSelectSpecResp, error) {
	out := new(SaveNewUserSelectSpecResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/saveNewUserSelectSpec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) UpdateUserStateStat(ctx context.Context, in *UpdateUserStateStatReq, opts ...grpc.CallOption) (*UpdateUserStateStatResp, error) {
	out := new(UpdateUserStateStatResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/updateUserStateStat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) SaveAdultQuizECR(ctx context.Context, in *SaveAdultQuizEcrReq, opts ...grpc.CallOption) (*SaveAdultQuizEcrResp, error) {
	out := new(SaveAdultQuizEcrResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/saveAdultQuizECR", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *statisticClient) GetAdultQuizEcr(ctx context.Context, in *GetAdultQuizEcrReq, opts ...grpc.CallOption) (*GetAdultQuizEcrResp, error) {
	out := new(GetAdultQuizEcrResp)
	err := c.cc.Invoke(ctx, "/pb.statistic/getAdultQuizEcr", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StatisticServer is the server API for Statistic service.
// All implementations must embed UnimplementedStatisticServer
// for forward compatibility
type StatisticServer interface {
	// 获取用户列表
	GetUserList(context.Context, *GetUserListReq) (*GetUserListResp, error)
	// 更新每日统计数据
	UpdateStatisticDailyData(context.Context, *UpdateStatisticDailyDataReq) (*UpdateStatisticDailyDataResp, error)
	// 获取每日统计数据
	GetDailyStatList(context.Context, *GetDailyStatListReq) (*GetDailyStatListResp, error)
	// 更新每日登陆时间
	UpdateLoginLog(context.Context, *UpdateLoginLogReq) (*UpdateLoginLogResp, error)
	// 获取统计近多少日的用户在昨日累计数据
	GetLifeTimeValueStat(context.Context, *GetLifeTimeValueStatReq) (*GetLifeTimeValueStatResp, error)
	// 获取用户渠道列表
	GetUserChannelList(context.Context, *GetUserChannelListReq) (*GetUserChannelListResp, error)
	// 新用户选择的XX标签
	SaveNewUserSelectSpec(context.Context, *SaveNewUserSelectSpecReq) (*SaveNewUserSelectSpecResp, error)
	// 定时统计用户和XXX状态数据
	UpdateUserStateStat(context.Context, *UpdateUserStateStatReq) (*UpdateUserStateStatResp, error)
	// 保存成人依恋量表测试结果
	SaveAdultQuizECR(context.Context, *SaveAdultQuizEcrReq) (*SaveAdultQuizEcrResp, error)
	// 获取最新成人依恋量表测试结果
	GetAdultQuizEcr(context.Context, *GetAdultQuizEcrReq) (*GetAdultQuizEcrResp, error)
	mustEmbedUnimplementedStatisticServer()
}

// UnimplementedStatisticServer must be embedded to have forward compatible implementations.
type UnimplementedStatisticServer struct {
}

func (UnimplementedStatisticServer) GetUserList(context.Context, *GetUserListReq) (*GetUserListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
func (UnimplementedStatisticServer) UpdateStatisticDailyData(context.Context, *UpdateStatisticDailyDataReq) (*UpdateStatisticDailyDataResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatisticDailyData not implemented")
}
func (UnimplementedStatisticServer) GetDailyStatList(context.Context, *GetDailyStatListReq) (*GetDailyStatListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDailyStatList not implemented")
}
func (UnimplementedStatisticServer) UpdateLoginLog(context.Context, *UpdateLoginLogReq) (*UpdateLoginLogResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLoginLog not implemented")
}
func (UnimplementedStatisticServer) GetLifeTimeValueStat(context.Context, *GetLifeTimeValueStatReq) (*GetLifeTimeValueStatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLifeTimeValueStat not implemented")
}
func (UnimplementedStatisticServer) GetUserChannelList(context.Context, *GetUserChannelListReq) (*GetUserChannelListResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserChannelList not implemented")
}
func (UnimplementedStatisticServer) SaveNewUserSelectSpec(context.Context, *SaveNewUserSelectSpecReq) (*SaveNewUserSelectSpecResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveNewUserSelectSpec not implemented")
}
func (UnimplementedStatisticServer) UpdateUserStateStat(context.Context, *UpdateUserStateStatReq) (*UpdateUserStateStatResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserStateStat not implemented")
}
func (UnimplementedStatisticServer) SaveAdultQuizECR(context.Context, *SaveAdultQuizEcrReq) (*SaveAdultQuizEcrResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveAdultQuizECR not implemented")
}
func (UnimplementedStatisticServer) GetAdultQuizEcr(context.Context, *GetAdultQuizEcrReq) (*GetAdultQuizEcrResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdultQuizEcr not implemented")
}
func (UnimplementedStatisticServer) mustEmbedUnimplementedStatisticServer() {}

// UnsafeStatisticServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StatisticServer will
// result in compilation errors.
type UnsafeStatisticServer interface {
	mustEmbedUnimplementedStatisticServer()
}

func RegisterStatisticServer(s grpc.ServiceRegistrar, srv StatisticServer) {
	s.RegisterService(&Statistic_ServiceDesc, srv)
}

func _Statistic_GetUserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetUserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/getUserList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetUserList(ctx, req.(*GetUserListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_UpdateStatisticDailyData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatisticDailyDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).UpdateStatisticDailyData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/updateStatisticDailyData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).UpdateStatisticDailyData(ctx, req.(*UpdateStatisticDailyDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_GetDailyStatList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDailyStatListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetDailyStatList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/getDailyStatList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetDailyStatList(ctx, req.(*GetDailyStatListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_UpdateLoginLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLoginLogReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).UpdateLoginLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/updateLoginLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).UpdateLoginLog(ctx, req.(*UpdateLoginLogReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_GetLifeTimeValueStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLifeTimeValueStatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetLifeTimeValueStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/getLifeTimeValueStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetLifeTimeValueStat(ctx, req.(*GetLifeTimeValueStatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_GetUserChannelList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserChannelListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetUserChannelList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/getUserChannelList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetUserChannelList(ctx, req.(*GetUserChannelListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_SaveNewUserSelectSpec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveNewUserSelectSpecReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).SaveNewUserSelectSpec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/saveNewUserSelectSpec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).SaveNewUserSelectSpec(ctx, req.(*SaveNewUserSelectSpecReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_UpdateUserStateStat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserStateStatReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).UpdateUserStateStat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/updateUserStateStat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).UpdateUserStateStat(ctx, req.(*UpdateUserStateStatReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_SaveAdultQuizECR_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveAdultQuizEcrReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).SaveAdultQuizECR(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/saveAdultQuizECR",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).SaveAdultQuizECR(ctx, req.(*SaveAdultQuizEcrReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Statistic_GetAdultQuizEcr_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdultQuizEcrReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StatisticServer).GetAdultQuizEcr(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.statistic/getAdultQuizEcr",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StatisticServer).GetAdultQuizEcr(ctx, req.(*GetAdultQuizEcrReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Statistic_ServiceDesc is the grpc.ServiceDesc for Statistic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Statistic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.statistic",
	HandlerType: (*StatisticServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getUserList",
			Handler:    _Statistic_GetUserList_Handler,
		},
		{
			MethodName: "updateStatisticDailyData",
			Handler:    _Statistic_UpdateStatisticDailyData_Handler,
		},
		{
			MethodName: "getDailyStatList",
			Handler:    _Statistic_GetDailyStatList_Handler,
		},
		{
			MethodName: "updateLoginLog",
			Handler:    _Statistic_UpdateLoginLog_Handler,
		},
		{
			MethodName: "getLifeTimeValueStat",
			Handler:    _Statistic_GetLifeTimeValueStat_Handler,
		},
		{
			MethodName: "getUserChannelList",
			Handler:    _Statistic_GetUserChannelList_Handler,
		},
		{
			MethodName: "saveNewUserSelectSpec",
			Handler:    _Statistic_SaveNewUserSelectSpec_Handler,
		},
		{
			MethodName: "updateUserStateStat",
			Handler:    _Statistic_UpdateUserStateStat_Handler,
		},
		{
			MethodName: "saveAdultQuizECR",
			Handler:    _Statistic_SaveAdultQuizECR_Handler,
		},
		{
			MethodName: "getAdultQuizEcr",
			Handler:    _Statistic_GetAdultQuizEcr_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statistic.proto",
}
