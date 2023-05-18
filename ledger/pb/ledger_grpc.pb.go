// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pb/ledger.proto

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

// LedgerServiceClient is the client API for LedgerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LedgerServiceClient interface {
	CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*CreatePaymentResponse, error)
	ReadPayment(ctx context.Context, in *ReadPaymentRequest, opts ...grpc.CallOption) (*ReadPaymentResponse, error)
	UpdatePaymentToPending(ctx context.Context, in *UpdatePaymentToPendingRequest, opts ...grpc.CallOption) (*UpdatePaymentToPendingResponse, error)
	UpdatePaymentToSuccess(ctx context.Context, in *UpdatePaymentToSuccessRequest, opts ...grpc.CallOption) (*UpdatePaymentToSuccessResponse, error)
	UpdatePaymentToFail(ctx context.Context, in *UpdatePaymentToFailRequest, opts ...grpc.CallOption) (*UpdatePaymentToFailResponse, error)
}

type ledgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerServiceClient(cc grpc.ClientConnInterface) LedgerServiceClient {
	return &ledgerServiceClient{cc}
}

func (c *ledgerServiceClient) CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...grpc.CallOption) (*CreatePaymentResponse, error) {
	out := new(CreatePaymentResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/CreatePayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) ReadPayment(ctx context.Context, in *ReadPaymentRequest, opts ...grpc.CallOption) (*ReadPaymentResponse, error) {
	out := new(ReadPaymentResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/ReadPayment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) UpdatePaymentToPending(ctx context.Context, in *UpdatePaymentToPendingRequest, opts ...grpc.CallOption) (*UpdatePaymentToPendingResponse, error) {
	out := new(UpdatePaymentToPendingResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/UpdatePaymentToPending", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) UpdatePaymentToSuccess(ctx context.Context, in *UpdatePaymentToSuccessRequest, opts ...grpc.CallOption) (*UpdatePaymentToSuccessResponse, error) {
	out := new(UpdatePaymentToSuccessResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/UpdatePaymentToSuccess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerServiceClient) UpdatePaymentToFail(ctx context.Context, in *UpdatePaymentToFailRequest, opts ...grpc.CallOption) (*UpdatePaymentToFailResponse, error) {
	out := new(UpdatePaymentToFailResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/UpdatePaymentToFail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LedgerServiceServer is the server API for LedgerService service.
// All implementations must embed UnimplementedLedgerServiceServer
// for forward compatibility
type LedgerServiceServer interface {
	CreatePayment(context.Context, *CreatePaymentRequest) (*CreatePaymentResponse, error)
	ReadPayment(context.Context, *ReadPaymentRequest) (*ReadPaymentResponse, error)
	UpdatePaymentToPending(context.Context, *UpdatePaymentToPendingRequest) (*UpdatePaymentToPendingResponse, error)
	UpdatePaymentToSuccess(context.Context, *UpdatePaymentToSuccessRequest) (*UpdatePaymentToSuccessResponse, error)
	UpdatePaymentToFail(context.Context, *UpdatePaymentToFailRequest) (*UpdatePaymentToFailResponse, error)
	mustEmbedUnimplementedLedgerServiceServer()
}

// UnimplementedLedgerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLedgerServiceServer struct {
}

func (UnimplementedLedgerServiceServer) CreatePayment(context.Context, *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePayment not implemented")
}
func (UnimplementedLedgerServiceServer) ReadPayment(context.Context, *ReadPaymentRequest) (*ReadPaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadPayment not implemented")
}
func (UnimplementedLedgerServiceServer) UpdatePaymentToPending(context.Context, *UpdatePaymentToPendingRequest) (*UpdatePaymentToPendingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePaymentToPending not implemented")
}
func (UnimplementedLedgerServiceServer) UpdatePaymentToSuccess(context.Context, *UpdatePaymentToSuccessRequest) (*UpdatePaymentToSuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePaymentToSuccess not implemented")
}
func (UnimplementedLedgerServiceServer) UpdatePaymentToFail(context.Context, *UpdatePaymentToFailRequest) (*UpdatePaymentToFailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePaymentToFail not implemented")
}
func (UnimplementedLedgerServiceServer) mustEmbedUnimplementedLedgerServiceServer() {}

// UnsafeLedgerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LedgerServiceServer will
// result in compilation errors.
type UnsafeLedgerServiceServer interface {
	mustEmbedUnimplementedLedgerServiceServer()
}

func RegisterLedgerServiceServer(s grpc.ServiceRegistrar, srv LedgerServiceServer) {
	s.RegisterService(&LedgerService_ServiceDesc, srv)
}

func _LedgerService_CreatePayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).CreatePayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/CreatePayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).CreatePayment(ctx, req.(*CreatePaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_ReadPayment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadPaymentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).ReadPayment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/ReadPayment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).ReadPayment(ctx, req.(*ReadPaymentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_UpdatePaymentToPending_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePaymentToPendingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).UpdatePaymentToPending(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/UpdatePaymentToPending",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).UpdatePaymentToPending(ctx, req.(*UpdatePaymentToPendingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_UpdatePaymentToSuccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePaymentToSuccessRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).UpdatePaymentToSuccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/UpdatePaymentToSuccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).UpdatePaymentToSuccess(ctx, req.(*UpdatePaymentToSuccessRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LedgerService_UpdatePaymentToFail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePaymentToFailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).UpdatePaymentToFail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/UpdatePaymentToFail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).UpdatePaymentToFail(ctx, req.(*UpdatePaymentToFailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LedgerService_ServiceDesc is the grpc.ServiceDesc for LedgerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LedgerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ledger.LedgerService",
	HandlerType: (*LedgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePayment",
			Handler:    _LedgerService_CreatePayment_Handler,
		},
		{
			MethodName: "ReadPayment",
			Handler:    _LedgerService_ReadPayment_Handler,
		},
		{
			MethodName: "UpdatePaymentToPending",
			Handler:    _LedgerService_UpdatePaymentToPending_Handler,
		},
		{
			MethodName: "UpdatePaymentToSuccess",
			Handler:    _LedgerService_UpdatePaymentToSuccess_Handler,
		},
		{
			MethodName: "UpdatePaymentToFail",
			Handler:    _LedgerService_UpdatePaymentToFail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/ledger.proto",
}