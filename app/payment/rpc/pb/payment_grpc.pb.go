// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: payment.proto

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

// PaymentClient is the client API for Payment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PaymentClient interface {
	// 创建微信支付预处理订单
	CreatePayment(ctx context.Context, in *CreatePaymentReq, opts ...grpc.CallOption) (*CreatePaymentResp, error)
	// 根据流水号查询流水记录
	GetPaymentByFlowNo(ctx context.Context, in *GetPaymentByFlowNoReq, opts ...grpc.CallOption) (*GetPaymentByFlowNoResp, error)
	// 更新交易状态
	UpdateTradeState(ctx context.Context, in *UpdateTradeStateReq, opts ...grpc.CallOption) (*UpdateTradeStateResp, error)
	// 根据订单id查询流水记录
	GetSuccessPaymentFlowByOrderIdReq(ctx context.Context, in *GetSuccessPaymentFlowByOrderIdReq, opts ...grpc.CallOption) (*GetSuccessPaymentFlowByOrderIdResp, error)
	// 发起退款
	RequestRefund(ctx context.Context, in *RequestRefundReq, opts ...grpc.CallOption) (*RequestRefundResp, error)
	// 更新退款状态
	UpdateRefundState(ctx context.Context, in *UpdateRefundReq, opts ...grpc.CallOption) (*UpdateRefundResp, error)
	// 银行卡转账
	MoveCash(ctx context.Context, in *MoveCashReq, opts ...grpc.CallOption) (*MoveCashResp, error)
	// 更新转账状态
	UpdateMoveCashStatus(ctx context.Context, in *UpdateMoveCashStatusReq, opts ...grpc.CallOption) (*UpdateMoveCashStatusResp, error)
}

type paymentClient struct {
	cc grpc.ClientConnInterface
}

func NewPaymentClient(cc grpc.ClientConnInterface) PaymentClient {
	return &paymentClient{cc}
}

func (c *paymentClient) CreatePayment(ctx context.Context, in *CreatePaymentReq, opts ...grpc.CallOption) (*CreatePaymentResp, error) {
	out := new(CreatePaymentResp)
	err := c.cc.Invoke(ctx, "/pb.payment/createPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetPaymentByFlowNo(ctx context.Context, in *GetPaymentByFlowNoReq, opts ...grpc.CallOption) (*GetPaymentByFlowNoResp, error) {
	out := new(GetPaymentByFlowNoResp)
	err := c.cc.Invoke(ctx, "/pb.payment/getPaymentByFlowNo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) UpdateTradeState(ctx context.Context, in *UpdateTradeStateReq, opts ...grpc.CallOption) (*UpdateTradeStateResp, error) {
	out := new(UpdateTradeStateResp)
	err := c.cc.Invoke(ctx, "/pb.payment/updateTradeState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) GetSuccessPaymentFlowByOrderIdReq(ctx context.Context, in *GetSuccessPaymentFlowByOrderIdReq, opts ...grpc.CallOption) (*GetSuccessPaymentFlowByOrderIdResp, error) {
	out := new(GetSuccessPaymentFlowByOrderIdResp)
	err := c.cc.Invoke(ctx, "/pb.payment/getSuccessPaymentFlowByOrderIdReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) RequestRefund(ctx context.Context, in *RequestRefundReq, opts ...grpc.CallOption) (*RequestRefundResp, error) {
	out := new(RequestRefundResp)
	err := c.cc.Invoke(ctx, "/pb.payment/requestRefund", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) UpdateRefundState(ctx context.Context, in *UpdateRefundReq, opts ...grpc.CallOption) (*UpdateRefundResp, error) {
	out := new(UpdateRefundResp)
	err := c.cc.Invoke(ctx, "/pb.payment/updateRefundState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) MoveCash(ctx context.Context, in *MoveCashReq, opts ...grpc.CallOption) (*MoveCashResp, error) {
	out := new(MoveCashResp)
	err := c.cc.Invoke(ctx, "/pb.payment/moveCash", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *paymentClient) UpdateMoveCashStatus(ctx context.Context, in *UpdateMoveCashStatusReq, opts ...grpc.CallOption) (*UpdateMoveCashStatusResp, error) {
	out := new(UpdateMoveCashStatusResp)
	err := c.cc.Invoke(ctx, "/pb.payment/updateMoveCashStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PaymentServer is the server API for Payment service.
// All implementations must embed UnimplementedPaymentServer
// for forward compatibility
type PaymentServer interface {
	// 创建微信支付预处理订单
	CreatePayment(context.Context, *CreatePaymentReq) (*CreatePaymentResp, error)
	// 根据流水号查询流水记录
	GetPaymentByFlowNo(context.Context, *GetPaymentByFlowNoReq) (*GetPaymentByFlowNoResp, error)
	// 更新交易状态
	UpdateTradeState(context.Context, *UpdateTradeStateReq) (*UpdateTradeStateResp, error)
	// 根据订单id查询流水记录
	GetSuccessPaymentFlowByOrderIdReq(context.Context, *GetSuccessPaymentFlowByOrderIdReq) (*GetSuccessPaymentFlowByOrderIdResp, error)
	// 发起退款
	RequestRefund(context.Context, *RequestRefundReq) (*RequestRefundResp, error)
	// 更新退款状态
	UpdateRefundState(context.Context, *UpdateRefundReq) (*UpdateRefundResp, error)
	// 银行卡转账
	MoveCash(context.Context, *MoveCashReq) (*MoveCashResp, error)
	// 更新转账状态
	UpdateMoveCashStatus(context.Context, *UpdateMoveCashStatusReq) (*UpdateMoveCashStatusResp, error)
	mustEmbedUnimplementedPaymentServer()
}

// UnimplementedPaymentServer must be embedded to have forward compatible implementations.
type UnimplementedPaymentServer struct {
}

func (UnimplementedPaymentServer) CreatePayment(context.Context, *CreatePaymentReq) (*CreatePaymentResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePayment not implemented")
}
func (UnimplementedPaymentServer) GetPaymentByFlowNo(context.Context, *GetPaymentByFlowNoReq) (*GetPaymentByFlowNoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaymentByFlowNo not implemented")
}
func (UnimplementedPaymentServer) UpdateTradeState(context.Context, *UpdateTradeStateReq) (*UpdateTradeStateResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTradeState not implemented")
}
func (UnimplementedPaymentServer) GetSuccessPaymentFlowByOrderIdReq(context.Context, *GetSuccessPaymentFlowByOrderIdReq) (*GetSuccessPaymentFlowByOrderIdResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSuccessPaymentFlowByOrderIdReq not implemented")
}
func (UnimplementedPaymentServer) RequestRefund(context.Context, *RequestRefundReq) (*RequestRefundResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestRefund not implemented")
}
func (UnimplementedPaymentServer) UpdateRefundState(context.Context, *UpdateRefundReq) (*UpdateRefundResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRefundState not implemented")
}
func (UnimplementedPaymentServer) MoveCash(context.Context, *MoveCashReq) (*MoveCashResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveCash not implemented")
}
func (UnimplementedPaymentServer) UpdateMoveCashStatus(context.Context, *UpdateMoveCashStatusReq) (*UpdateMoveCashStatusResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMoveCashStatus not implemented")
}
func (UnimplementedPaymentServer) mustEmbedUnimplementedPaymentServer() {}

// UnsafePaymentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PaymentServer will
// result in compilation errors.
type UnsafePaymentServer interface {
	mustEmbedUnimplementedPaymentServer()
}

func RegisterPaymentServer(s grpc.ServiceRegistrar, srv PaymentServer) {
	s.RegisterService(&Payment_ServiceDesc, srv)
}

func _Payment_CreatePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePaymentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).CreatePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/createPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).CreatePayment(ctx, req.(*CreatePaymentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetPaymentByFlowNo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPaymentByFlowNoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetPaymentByFlowNo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/getPaymentByFlowNo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetPaymentByFlowNo(ctx, req.(*GetPaymentByFlowNoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_UpdateTradeState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTradeStateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).UpdateTradeState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/updateTradeState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).UpdateTradeState(ctx, req.(*UpdateTradeStateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_GetSuccessPaymentFlowByOrderIdReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSuccessPaymentFlowByOrderIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).GetSuccessPaymentFlowByOrderIdReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/getSuccessPaymentFlowByOrderIdReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).GetSuccessPaymentFlowByOrderIdReq(ctx, req.(*GetSuccessPaymentFlowByOrderIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_RequestRefund_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestRefundReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).RequestRefund(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/requestRefund",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).RequestRefund(ctx, req.(*RequestRefundReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_UpdateRefundState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRefundReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).UpdateRefundState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/updateRefundState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).UpdateRefundState(ctx, req.(*UpdateRefundReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_MoveCash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveCashReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).MoveCash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/moveCash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).MoveCash(ctx, req.(*MoveCashReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Payment_UpdateMoveCashStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMoveCashStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PaymentServer).UpdateMoveCashStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.payment/updateMoveCashStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PaymentServer).UpdateMoveCashStatus(ctx, req.(*UpdateMoveCashStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Payment_ServiceDesc is the grpc.ServiceDesc for Payment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Payment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.payment",
	HandlerType: (*PaymentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createPayment",
			Handler:    _Payment_CreatePayment_Handler,
		},
		{
			MethodName: "getPaymentByFlowNo",
			Handler:    _Payment_GetPaymentByFlowNo_Handler,
		},
		{
			MethodName: "updateTradeState",
			Handler:    _Payment_UpdateTradeState_Handler,
		},
		{
			MethodName: "getSuccessPaymentFlowByOrderIdReq",
			Handler:    _Payment_GetSuccessPaymentFlowByOrderIdReq_Handler,
		},
		{
			MethodName: "requestRefund",
			Handler:    _Payment_RequestRefund_Handler,
		},
		{
			MethodName: "updateRefundState",
			Handler:    _Payment_UpdateRefundState_Handler,
		},
		{
			MethodName: "moveCash",
			Handler:    _Payment_MoveCash_Handler,
		},
		{
			MethodName: "updateMoveCashStatus",
			Handler:    _Payment_UpdateMoveCashStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "payment.proto",
}
