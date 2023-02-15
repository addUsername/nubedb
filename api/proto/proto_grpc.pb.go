// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/proto/proto.proto

package proto

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	ExecuteOnLeader(ctx context.Context, in *ExecuteOnLeaderRequest, opts ...grpc.CallOption) (*Empty, error)
	IsLeader(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*IsLeaderResponse, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) ExecuteOnLeader(ctx context.Context, in *ExecuteOnLeaderRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Service/ExecuteOnLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) IsLeader(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*IsLeaderResponse, error) {
	out := new(IsLeaderResponse)
	err := c.cc.Invoke(ctx, "/proto.Service/IsLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	ExecuteOnLeader(context.Context, *ExecuteOnLeaderRequest) (*Empty, error)
	IsLeader(context.Context, *Empty) (*IsLeaderResponse, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) ExecuteOnLeader(context.Context, *ExecuteOnLeaderRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteOnLeader not implemented")
}
func (UnimplementedServiceServer) IsLeader(context.Context, *Empty) (*IsLeaderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsLeader not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_ExecuteOnLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecuteOnLeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ExecuteOnLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Service/ExecuteOnLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ExecuteOnLeader(ctx, req.(*ExecuteOnLeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_IsLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).IsLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Service/IsLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).IsLeader(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteOnLeader",
			Handler:    _Service_ExecuteOnLeader_Handler,
		},
		{
			MethodName: "IsLeader",
			Handler:    _Service_IsLeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/proto.proto",
}
