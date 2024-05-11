// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: service.proto

package miniresolverproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MiniResolver_Ping_FullMethodName            = "/miniresolverproto.MiniResolver/Ping"
	MiniResolver_AddService_FullMethodName      = "/miniresolverproto.MiniResolver/AddService"
	MiniResolver_RemoveService_FullMethodName   = "/miniresolverproto.MiniResolver/RemoveService"
	MiniResolver_ResolveService_FullMethodName  = "/miniresolverproto.MiniResolver/ResolveService"
	MiniResolver_ResolveServices_FullMethodName = "/miniresolverproto.MiniResolver/ResolveServices"
)

// MiniResolverClient is the client API for MiniResolver service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MiniResolverClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DefaultResponse, error)
	AddService(ctx context.Context, in *ServiceData, opts ...grpc.CallOption) (*DefaultResponse, error)
	RemoveService(ctx context.Context, in *ServiceData, opts ...grpc.CallOption) (*DefaultResponse, error)
	ResolveService(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.StringValue, error)
	ResolveServices(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*ServiceResponse, error)
}

type miniResolverClient struct {
	cc grpc.ClientConnInterface
}

func NewMiniResolverClient(cc grpc.ClientConnInterface) MiniResolverClient {
	return &miniResolverClient{cc}
}

func (c *miniResolverClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DefaultResponse, error) {
	out := new(DefaultResponse)
	err := c.cc.Invoke(ctx, MiniResolver_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *miniResolverClient) AddService(ctx context.Context, in *ServiceData, opts ...grpc.CallOption) (*DefaultResponse, error) {
	out := new(DefaultResponse)
	err := c.cc.Invoke(ctx, MiniResolver_AddService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *miniResolverClient) RemoveService(ctx context.Context, in *ServiceData, opts ...grpc.CallOption) (*DefaultResponse, error) {
	out := new(DefaultResponse)
	err := c.cc.Invoke(ctx, MiniResolver_RemoveService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *miniResolverClient) ResolveService(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*wrapperspb.StringValue, error) {
	out := new(wrapperspb.StringValue)
	err := c.cc.Invoke(ctx, MiniResolver_ResolveService_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *miniResolverClient) ResolveServices(ctx context.Context, in *wrapperspb.StringValue, opts ...grpc.CallOption) (*ServiceResponse, error) {
	out := new(ServiceResponse)
	err := c.cc.Invoke(ctx, MiniResolver_ResolveServices_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MiniResolverServer is the server API for MiniResolver service.
// All implementations must embed UnimplementedMiniResolverServer
// for forward compatibility
type MiniResolverServer interface {
	Ping(context.Context, *emptypb.Empty) (*DefaultResponse, error)
	AddService(context.Context, *ServiceData) (*DefaultResponse, error)
	RemoveService(context.Context, *ServiceData) (*DefaultResponse, error)
	ResolveService(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error)
	ResolveServices(context.Context, *wrapperspb.StringValue) (*ServiceResponse, error)
	mustEmbedUnimplementedMiniResolverServer()
}

// UnimplementedMiniResolverServer must be embedded to have forward compatible implementations.
type UnimplementedMiniResolverServer struct {
}

func (UnimplementedMiniResolverServer) Ping(context.Context, *emptypb.Empty) (*DefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedMiniResolverServer) AddService(context.Context, *ServiceData) (*DefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddService not implemented")
}
func (UnimplementedMiniResolverServer) RemoveService(context.Context, *ServiceData) (*DefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveService not implemented")
}
func (UnimplementedMiniResolverServer) ResolveService(context.Context, *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveService not implemented")
}
func (UnimplementedMiniResolverServer) ResolveServices(context.Context, *wrapperspb.StringValue) (*ServiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveServices not implemented")
}
func (UnimplementedMiniResolverServer) mustEmbedUnimplementedMiniResolverServer() {}

// UnsafeMiniResolverServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MiniResolverServer will
// result in compilation errors.
type UnsafeMiniResolverServer interface {
	mustEmbedUnimplementedMiniResolverServer()
}

func RegisterMiniResolverServer(s grpc.ServiceRegistrar, srv MiniResolverServer) {
	s.RegisterService(&MiniResolver_ServiceDesc, srv)
}

func _MiniResolver_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniResolverServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniResolver_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniResolverServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _MiniResolver_AddService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniResolverServer).AddService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniResolver_AddService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniResolverServer).AddService(ctx, req.(*ServiceData))
	}
	return interceptor(ctx, in, info, handler)
}

func _MiniResolver_RemoveService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniResolverServer).RemoveService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniResolver_RemoveService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniResolverServer).RemoveService(ctx, req.(*ServiceData))
	}
	return interceptor(ctx, in, info, handler)
}

func _MiniResolver_ResolveService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniResolverServer).ResolveService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniResolver_ResolveService_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniResolverServer).ResolveService(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _MiniResolver_ResolveServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrapperspb.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiniResolverServer).ResolveServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MiniResolver_ResolveServices_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiniResolverServer).ResolveServices(ctx, req.(*wrapperspb.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// MiniResolver_ServiceDesc is the grpc.ServiceDesc for MiniResolver service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MiniResolver_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "miniresolverproto.MiniResolver",
	HandlerType: (*MiniResolverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _MiniResolver_Ping_Handler,
		},
		{
			MethodName: "AddService",
			Handler:    _MiniResolver_AddService_Handler,
		},
		{
			MethodName: "RemoveService",
			Handler:    _MiniResolver_RemoveService_Handler,
		},
		{
			MethodName: "ResolveService",
			Handler:    _MiniResolver_ResolveService_Handler,
		},
		{
			MethodName: "ResolveServices",
			Handler:    _MiniResolver_ResolveServices_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}