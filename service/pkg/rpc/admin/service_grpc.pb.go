// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package admin

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

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminServiceClient interface {
	GetAuthChallenge(ctx context.Context, in *GetAuthChallenge_Request, opts ...grpc.CallOption) (*GetAuthChallenge_Response, error)
	Auth(ctx context.Context, in *Auth_Request, opts ...grpc.CallOption) (*Auth_Response, error)
	SendMessageToRoom(ctx context.Context, in *SendMessageToRoom_Request, opts ...grpc.CallOption) (*SendMessageToRoom_Response, error)
	DumpDirectMessages(ctx context.Context, in *DumpDirectMessages_Request, opts ...grpc.CallOption) (AdminService_DumpDirectMessagesClient, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) GetAuthChallenge(ctx context.Context, in *GetAuthChallenge_Request, opts ...grpc.CallOption) (*GetAuthChallenge_Response, error) {
	out := new(GetAuthChallenge_Response)
	err := c.cc.Invoke(ctx, "/AdminService/GetAuthChallenge", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) Auth(ctx context.Context, in *Auth_Request, opts ...grpc.CallOption) (*Auth_Response, error) {
	out := new(Auth_Response)
	err := c.cc.Invoke(ctx, "/AdminService/Auth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) SendMessageToRoom(ctx context.Context, in *SendMessageToRoom_Request, opts ...grpc.CallOption) (*SendMessageToRoom_Response, error) {
	out := new(SendMessageToRoom_Response)
	err := c.cc.Invoke(ctx, "/AdminService/SendMessageToRoom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) DumpDirectMessages(ctx context.Context, in *DumpDirectMessages_Request, opts ...grpc.CallOption) (AdminService_DumpDirectMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &AdminService_ServiceDesc.Streams[0], "/AdminService/DumpDirectMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &adminServiceDumpDirectMessagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AdminService_DumpDirectMessagesClient interface {
	Recv() (*DumpDirectMessages_Response, error)
	grpc.ClientStream
}

type adminServiceDumpDirectMessagesClient struct {
	grpc.ClientStream
}

func (x *adminServiceDumpDirectMessagesClient) Recv() (*DumpDirectMessages_Response, error) {
	m := new(DumpDirectMessages_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AdminServiceServer is the server API for AdminService service.
// All implementations must embed UnimplementedAdminServiceServer
// for forward compatibility
type AdminServiceServer interface {
	GetAuthChallenge(context.Context, *GetAuthChallenge_Request) (*GetAuthChallenge_Response, error)
	Auth(context.Context, *Auth_Request) (*Auth_Response, error)
	SendMessageToRoom(context.Context, *SendMessageToRoom_Request) (*SendMessageToRoom_Response, error)
	DumpDirectMessages(*DumpDirectMessages_Request, AdminService_DumpDirectMessagesServer) error
	mustEmbedUnimplementedAdminServiceServer()
}

// UnimplementedAdminServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (UnimplementedAdminServiceServer) GetAuthChallenge(context.Context, *GetAuthChallenge_Request) (*GetAuthChallenge_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthChallenge not implemented")
}
func (UnimplementedAdminServiceServer) Auth(context.Context, *Auth_Request) (*Auth_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}
func (UnimplementedAdminServiceServer) SendMessageToRoom(context.Context, *SendMessageToRoom_Request) (*SendMessageToRoom_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessageToRoom not implemented")
}
func (UnimplementedAdminServiceServer) DumpDirectMessages(*DumpDirectMessages_Request, AdminService_DumpDirectMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method DumpDirectMessages not implemented")
}
func (UnimplementedAdminServiceServer) mustEmbedUnimplementedAdminServiceServer() {}

// UnsafeAdminServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServiceServer will
// result in compilation errors.
type UnsafeAdminServiceServer interface {
	mustEmbedUnimplementedAdminServiceServer()
}

func RegisterAdminServiceServer(s grpc.ServiceRegistrar, srv AdminServiceServer) {
	s.RegisterService(&AdminService_ServiceDesc, srv)
}

func _AdminService_GetAuthChallenge_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthChallenge_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetAuthChallenge(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/GetAuthChallenge",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetAuthChallenge(ctx, req.(*GetAuthChallenge_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Auth_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/Auth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).Auth(ctx, req.(*Auth_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_SendMessageToRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessageToRoom_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).SendMessageToRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminService/SendMessageToRoom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).SendMessageToRoom(ctx, req.(*SendMessageToRoom_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_DumpDirectMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DumpDirectMessages_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AdminServiceServer).DumpDirectMessages(m, &adminServiceDumpDirectMessagesServer{stream})
}

type AdminService_DumpDirectMessagesServer interface {
	Send(*DumpDirectMessages_Response) error
	grpc.ServerStream
}

type adminServiceDumpDirectMessagesServer struct {
	grpc.ServerStream
}

func (x *adminServiceDumpDirectMessagesServer) Send(m *DumpDirectMessages_Response) error {
	return x.ServerStream.SendMsg(m)
}

// AdminService_ServiceDesc is the grpc.ServiceDesc for AdminService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthChallenge",
			Handler:    _AdminService_GetAuthChallenge_Handler,
		},
		{
			MethodName: "Auth",
			Handler:    _AdminService_Auth_Handler,
		},
		{
			MethodName: "SendMessageToRoom",
			Handler:    _AdminService_SendMessageToRoom_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "DumpDirectMessages",
			Handler:       _AdminService_DumpDirectMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/rpc/admin/service.proto",
}
